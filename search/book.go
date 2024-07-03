package search

import (
	"context"
	"fmt"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type BookSearchClient struct {
	interval     time.Duration
	index        bleve.Index
	indexMapping *mapping.IndexMappingImpl
}

// 搜索字段
type SearchField struct {
	Name, Value string
}

func NewBookSearchClient(interval time.Duration) *BookSearchClient {
	docMapping := bleve.NewDocumentMapping()
	newTextMapping := bleve.NewTextFieldMapping()
	newTextMapping.Analyzer = standard.Name
	docMapping.AddFieldMappingsAt("id", newTextMapping)
	docMapping.AddFieldMappingsAt("tags", newTextMapping)
	simpleTextMapping := bleve.NewTextFieldMapping()
	simpleTextMapping.Analyzer = simple.Name
	// simple 兼容中文和英文
	docMapping.AddFieldMappingsAt("name", simpleTextMapping)
	docMapping.AddFieldMappingsAt("author", newTextMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultAnalyzer = standard.Name
	indexMapping.AddDocumentMapping("book", docMapping)

	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		panic(err)
	}
	cli := &BookSearchClient{
		interval:     interval,
		index:        index,
		indexMapping: indexMapping,
	}
	cli.RebuildIndex()
	return cli
}

// RebuildIndex ...
func (client *BookSearchClient) RebuildIndex() {
	// TODO: Rebuild this from database
	books := GetLocalBooks()
	for _, v := range books {
		client.AddDoc(fmt.Sprint(v.Id), v)
	}
}

// AddDoc 添加文档
func (client *BookSearchClient) AddDoc(id string, data any) {
	err := client.index.Index(id, data)
	if err != nil {
		panic(err)
	}
}

func (client *BookSearchClient) BuildShouldQuery(fields []*SearchField) *query.BooleanQuery {
	query := bleve.NewBooleanQuery()
	for _, field := range fields {
		termQuery := bleve.NewTermQuery(field.Value) // 匹配字段值
		termQuery.SetField(field.Name)
		query.AddShould(termQuery)
	}
	return query
}

func (client *BookSearchClient) Search(req *bleve.SearchRequest) (*bleve.SearchResult, error) {
	return client.index.Search(req)
}

func (client *BookSearchClient) Run(ctx context.Context) {
	ticker := time.NewTicker(client.interval)
	defer ticker.Stop()
	hlog.Info("BookSearchClient is running.")

	for {
		select {
		case <-ctx.Done():
			hlog.Info("BookSearchClient Done.")
			return
		case <-ticker.C:
			now := time.Now()
			hlog.Infof("rebuild book index at: %v", now.Unix())
			client.RebuildIndex()
		}
	}
}
