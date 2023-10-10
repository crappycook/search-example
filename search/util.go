package search

import "gosearcher/model"

func GetLocalBooks() map[int64]*model.Book {
	data := map[int64]*model.Book{
		1: {
			Id:   1,
			Name: "Network",
			Tags: "cs new",
		},
		2: {
			Id:   2,
			Name: "Operating System",
			Tags: "cs",
		},
		3: {
			Id:   3,
			Name: "442",
			Tags: "new football",
		},
		4: {
			Id:   4,
			Name: "Math",
			Tags: "science",
		},
		5: {
			Id:   5,
			Name: "Deep Learning",
			Tags: "science new",
		},
	}
	return data
}
