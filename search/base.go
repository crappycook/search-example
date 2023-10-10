package search

import (
	"context"
	"time"
)

var (
	BookCli *BookSearchClient
)

func Init() {
	ctx := context.Background()
	BookCli = NewBookSearchClient(time.Second * 10)
	BookCli.Run(ctx)
}
