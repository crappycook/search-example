package search

import (
	"context"
	"fmt"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
)

type BookSearchClient struct {
	interval     time.Duration
	index        bleve.Index
	indexMapping *mapping.IndexMappingImpl
	// documentMapping *mapping.DocumentMapping
}

// 搜索字段
type SearchField struct {
	Name, Value string
}

func NewBookSearchClient(interval time.Duration) *BookSearchClient {
	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultAnalyzer = simple.Name
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

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				client.RebuildIndex()
			}
		}
	}()
}
