package v1

import (
	"book_service/pkg/clients"
	"book_service/pkg/models/common"
	"book_service/pkg/models/common/req"
	"book_service/pkg/query"
	"book_service/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"net/http"
)

var booksIndex, _ = utils.GetEnvVar[string]("BOOKS_INDEX", "books")

func GetBookById(c *gin.Context) {
	payload, _ := c.Get("validated")
	bookId, _ := payload.(req.GetBook)

	// TODO: it should be in the validation level with extra validators and transformations
	_, parseErr := uuid.Parse(bookId.ID)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": parseErr})
		return
	}

	esQuery := query.NewQueryBuilder().ID(bookId.ID).Build()
	fmt.Println(esQuery)
	hits, err := clients.DoSearch(c, booksIndex, esQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	if len(hits) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	c.JSON(200, hits[0])
}

func CreateBook(c *gin.Context) {
	payload, _ := c.Get("validated")
	bodyBookReq, _ := payload.(req.AddBook)
	book := common.Book{ID: uuid.New(), PublishDate: bodyBookReq.PublishDate.Format("2006-01-02")}

	_ = copier.Copy(&book, &bodyBookReq)

	clients.EnqueueIndexTask(c, booksIndex, book.ID.String(), book)
	logrus.Info("INSERTING THE BOOK")
	c.JSON(http.StatusAccepted, gin.H{
		"id": book.ID,
	})
}
