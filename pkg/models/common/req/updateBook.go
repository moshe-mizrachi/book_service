package req

type UpdateBook struct {
	ID    string `uri:"id" binding:"required" validate:"required"`
	Title string `json:"title" validate:"required,min=1"`
}
