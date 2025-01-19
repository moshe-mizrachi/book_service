package req

import (
	face "book_service/pkg/interfaces"
	"book_service/pkg/utils"
	"errors"
)

var _ face.Validatable = (*GetBook)(nil)

func (g *GetBook) Validate() error {
	validUUID := utils.IsValidUUID(g.ID)
	if !validUUID {
		return errors.New("invalid uuid")
	}
	return nil
}

type GetBook struct {
	ID string `uri:"id" binding:"required"`
}
