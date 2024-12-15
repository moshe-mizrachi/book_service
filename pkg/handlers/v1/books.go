package v1

import (
	"book_service/pkg/clients"
	_const "book_service/pkg/constants"
	"book_service/pkg/models/common"
	"book_service/pkg/models/common/req"
	"book_service/pkg/models/common/res"
	"book_service/pkg/query"
	"book_service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

var booksIndex, _ = utils.GetEnvVar[string]("BOOKS_INDEX", "books")

func GetBookById(c *gin.Context) {
	bookReq := utils.GetValidatedPayload[req.GetBook](c)

	esQuery := query.NewQueryBuilder().ID(bookReq.ID).Build()
	hits, _, err := clients.DoSearch(c, booksIndex, esQuery, 1, 0)
	if err != nil {
		logrus.Errorf("Error searching for book with ID %s: %v", bookReq.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if len(hits) == 0 {
		logrus.Infof("Book with ID %s not found", bookReq.ID)
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	logrus.Infof("Book with ID %s retrieved successfully", bookReq.ID)
	c.JSON(http.StatusOK, hits[0])
}

func CreateBook(c *gin.Context) {
	bodyBookReq := utils.GetValidatedPayload[req.AddBook](c)
	book := common.Book{
		ID:          uuid.New(),
		PublishDate: bodyBookReq.PublishDate.Format("2006-01-02"),
	}

	if err := copier.Copy(&book, &bodyBookReq); err != nil {
		logrus.Errorf("Error copying AddBook payload to Book: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to map request data"})
		return
	}

	clients.EnqueueIndexTask(c, booksIndex, book.ID.String(), book, _const.CreateIndex)
	logrus.Infof("Book with ID %s queued for creation successfully", book.ID)
	c.JSON(http.StatusAccepted, res.AddBook{ID: book.ID})
}

func UpdateBook(c *gin.Context) {
	bodyBookReq := utils.GetValidatedPayload[req.UpdateBook](c)
	titleUpdate := common.TitleUpdate{}
	if err := copier.Copy(&titleUpdate, &bodyBookReq); err != nil {
		logrus.Errorf("Error copying UpdateBook payload to TitleUpdate: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to map request data"})
		return
	}

	clients.EnqueueIndexTask(c, booksIndex, bodyBookReq.ID, titleUpdate, _const.UpdateIndex)
	logrus.Infof("Book with ID %s queued for update successfully", bodyBookReq.ID)
	c.JSON(http.StatusAccepted, res.UpdateBook{ID: uuid.MustParse(bodyBookReq.ID)})
}

func DeleteBook(c *gin.Context) {
	deleteReq := utils.GetValidatedPayload[req.DeleteBook](c)

	clients.EnqueueIndexTask(c, booksIndex, deleteReq.ID, "", _const.DeleteIndex)
	logrus.Infof("Book with ID %s queued for deletion successfully", deleteReq.ID)
	c.JSON(http.StatusAccepted, res.DeleteBook{ID: uuid.MustParse(deleteReq.ID)})
}

func SearchBooks(c *gin.Context) {
	searchReq := utils.GetValidatedPayload[req.SearchBooks](c)

	esQuery := query.NewQueryBuilder().
		Title(searchReq.Title).
		PriceRange(searchReq.PriceRange.Min, searchReq.PriceRange.Max).
		Build()

	hits, _, err := clients.DoSearch(c, booksIndex, esQuery, searchReq.Size, searchReq.From)
	if err != nil {
		logrus.Errorf("Error searching books: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	logrus.Infof("SearchBooks query executed successfully, retrieved %d results", len(hits))
	c.JSON(http.StatusOK, hits)
}

func GetBooksStats(c *gin.Context) {
	esQuery := query.NewQueryBuilder().
		DistinctAuthors().
		TotalBooks().
		Build()

	_, aggregations, err := clients.DoSearch(c, booksIndex, esQuery, 0, 0)
	if err != nil {
		logrus.Errorf("Error fetching books statistics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	stats, err := utils.ParseAggregations(aggregations, _const.AggregationConfigs["BookStats"])
	if err != nil {
		logrus.Errorf("Error parsing aggregations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	logrus.Infof("GetBooksStats executed successfully: %+v", stats)
	c.JSON(http.StatusOK, stats)
}
