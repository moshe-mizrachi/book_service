package res

import (
	"book_service/pkg/models/common"

	"github.com/google/uuid"
)

type AddBook struct {
	ID uuid.UUID `json:"id"`
}

type UpdateBook struct {
	ID uuid.UUID `json:"id"`
}

type DeleteBook struct {
	ID uuid.UUID `json:"message"`
}

type SearchBooks struct {
	Books []common.Book `json:"books"`
}

type BookSummary struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	AuthorName string  `json:"author_name"`
	Price      float64 `json:"price"`
}

type StoreSummary struct {
	BookCount       int `json:"book_count"`
	DistinctAuthors int `json:"distinct_authors"`
}
