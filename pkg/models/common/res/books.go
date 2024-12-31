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
