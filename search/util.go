package search

import (
	"encoding/json"
	"gosearcher/model"
	"os"
)

var (
	DataPath = "data"
)

func GetLocalBooks() map[int64]*model.Book {
	content, err := os.ReadFile(DataPath + "/books.json")
	if err != nil {
		panic(err)
	}

	var books []*model.Book
	err = json.Unmarshal(content, &books)
	if err != nil {
		panic(err)
	}

	result := make(map[int64]*model.Book, len(books))
	for _, v := range books {
		result[v.Id] = v
	}
	return result
}

func JsonCompact(val any) string {
	b, _ := json.Marshal(val)
	return string(b)
}
