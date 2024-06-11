package model

type Book struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Tags   string `json:"tags"`
}

// A bleveClassifier is an interface describing any object which knows how
// to identify its own type.
func (b *Book) BleveType() string {
	return "book"
}
