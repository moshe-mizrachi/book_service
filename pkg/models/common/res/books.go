package res

import (
	"time"
)

type AddBookResponse struct {
	ID int `json:"id"`
}

type UpdateBookResponse struct {
	Message string `json:"message"`
}

type GetBookResponse struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	AuthorName     string    `json:"author_name"`
	Price          float64   `json:"price"`
	EbookAvailable bool      `json:"ebook_available"`
	PublishDate    time.Time `json:"publish_date"`
}

type DeleteBookResponse struct {
	Message string `json:"message"`
}

type SearchBooksResponse struct {
	Books []BookSummary `json:"books"`
}

type BookSummary struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	AuthorName string  `json:"author_name"`
	Price      float64 `json:"price"`
}

type StoreSummaryResponse struct {
	BookCount       int `json:"book_count"`
	DistinctAuthors int `json:"distinct_authors"`
}
