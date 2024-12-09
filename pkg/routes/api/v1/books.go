package v1

import (
	handlers "book_service/pkg/handlers/v1"
	mw "book_service/pkg/middlewares"
	"book_service/pkg/models/common/req"
	"github.com/gin-gonic/gin"
)

type empty struct {
}

func RegisterBooksRoutes(rgp *gin.RouterGroup) {
	v1 := rgp.Group("/v1/books")
	{
		v1.GET("/:id", mw.Validation[req.GetBook](), handlers.GetBookById)

		v1.POST("/", mw.Validation[req.AddBook](), handlers.CreateBook)

		//v1.GET("/:id", func(c *gin.Context) {
		//	id := c.Param("id")
		//	c.JSON(200, gin.H{"book_id": id, "message": "Book details fetched successfully!"})
		//})
	}
}
