package model

type Book struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Tags string `json:"tags"`
}
