package handlers

import (
	"context"
	"gosearcher/model"
	"gosearcher/search"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type SearchBooksRequest struct {
	Size int    `query:"size" binding:"required"`
	From int    `query:"from"`
	Tags string `query:"tags"`
}

func SearchBooks(c context.Context, ctx *app.RequestContext) {
	var req SearchBooksRequest
	err := ctx.BindQuery(&req)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, utils.H{"message": err.Error()})
		return
	}
	// Build search fields
	tagList := strings.Split(req.Tags, ",")
	searchFields := make([]*search.SearchField, 0, len(tagList))
	for _, tag := range tagList {
		searchFields = append(searchFields, &search.SearchField{
			Name:  "tags",
			Value: tag,
		})
	}

	query := search.BookCli.BuildShouldQuery(searchFields)
	searchReq := bleve.NewSearchRequest(query)
	// sort by doc _id desc
	searchReq.SortBy([]string{"-_id"})
	searchReq.Size = req.Size
	searchReq.From = req.From
	searchReq.Fields = []string{"*"}
	searchResults, err := search.BookCli.Search(searchReq)
	hlog.Infof("query took: %v", searchResults.Took)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, utils.H{"message": err})
		return
	}
	// Get book ids
	ids := make([]int64, 0, searchResults.Total)
	if searchResults.Total > 0 {
		for _, hit := range searchResults.Hits {
			id, _ := strconv.ParseInt(hit.ID, 10, 64)
			ids = append(ids, id)
		}
	}

	// Get book info from database by this ids
	result := make([]*model.Book, 0, len(ids))
	bookMap := search.GetLocalBooks()
	for _, id := range ids {
		if v, ok := bookMap[id]; ok {
			result = append(result, v)
		}
	}
	ctx.JSON(consts.StatusOK, result)
}

type FuzzySearchBooksRequest struct {
	Size    int    `query:"size" binding:"required"`
	From    int    `query:"from"`
	Keyword string `query:"keyword"`
}

func FuzzySearchBooks(c context.Context, ctx *app.RequestContext) {
	var req FuzzySearchBooksRequest
	err := ctx.BindQuery(&req)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, utils.H{"message": err.Error()})
		return
	}

	var query query.Query
	if len(req.Keyword) == 0 {
		query = bleve.NewMatchAllQuery()
	} else {
		authorQuery := bleve.NewWildcardQuery("*" + req.Keyword + "*")
		authorQuery.SetField("author")

		// nameQuery := bleve.NewBooleanQuery()
		// for _, v := range req.Keyword {
		// 	q := bleve.NewTermQuery(string(v))
		// 	q.SetField("name")
		// 	nameQuery.AddMust(q)
		// }

		nameWildcardQuery := bleve.NewWildcardQuery("*" + req.Keyword + "*")
		nameWildcardQuery.SetField("name")

		boolQuery := bleve.NewBooleanQuery()
		boolQuery.AddShould(authorQuery)
		// boolQuery.AddShould(nameQuery)
		boolQuery.AddShould(nameWildcardQuery)
		hlog.Infof("query: %v", search.JsonCompact(boolQuery))
		query = boolQuery
	}

	searchReq := bleve.NewSearchRequest(query)
	// sort by doc _id desc
	searchReq.SortBy([]string{"-_id"})
	// searchReq.Size = req.Size
	// searchReq.From = req.From
	searchReq.Fields = []string{"*"}
	searchResults, err := search.BookCli.Search(searchReq)
	hlog.Infof("query took: %v", searchResults.Took)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, utils.H{"message": err})
		return
	}
	// Get book ids
	ids := make([]int64, 0, searchResults.Total)
	if searchResults.Total > 0 {
		for _, hit := range searchResults.Hits {
			id, _ := strconv.ParseInt(hit.ID, 10, 64)
			ids = append(ids, id)
		}
	}

	// Get book info from database by this ids
	result := make([]*model.Book, 0, len(ids))
	bookMap := search.GetLocalBooks()
	for _, id := range ids {
		if v, ok := bookMap[id]; ok {
			result = append(result, v)
		}
	}
	ctx.JSON(consts.StatusOK, result)
}
