package req

type DeleteBook struct {
	ID string `uri:"id" binding:"required" validate:"required"`
}
