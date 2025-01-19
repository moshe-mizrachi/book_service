package v1

import (
	handlers "book_service/pkg/handlers/v1"
	mw "book_service/pkg/middlewares"
	"book_service/pkg/models/common/req"

	"github.com/gin-gonic/gin"
)

func RegisterBooksRoutes(rgp *gin.RouterGroup) {
	v1 := rgp.Group("/v1/books")
	{
		v1.GET("/:id", mw.Validation[req.GetBook](), handlers.GetBookById)
		v1.PUT("/:id", mw.Validation[req.UpdateBook](), handlers.UpdateBook) // why just title
		v1.DELETE("/:id", mw.Validation[req.DeleteBook](), handlers.DeleteBook)
		v1.GET("/search", mw.Validation[req.SearchBooks](), handlers.SearchBooks) // the good pattern for search is to put it into body due to size
		v1.POST("/", mw.Validation[req.AddBook](), handlers.CreateBook)
	}
}
