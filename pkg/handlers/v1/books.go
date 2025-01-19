package v1

import (
	"book_service/pkg/clients"
	"book_service/pkg/consts"
	"book_service/pkg/models/common"
	"book_service/pkg/models/common/req"
	"book_service/pkg/models/common/res"
	"book_service/pkg/query"
	"book_service/pkg/utils"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

func GetBookById(c *gin.Context) {
	bookReq, err := utils.GetValidatedPayload[req.GetBook](c)
	if err != nil {
		log.Errorf("Error getting book by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	esQuery := query.NewQueryBuilder().ID(bookReq.ID).Build()
	hits, _, err := clients.SearchIndex(c, esQuery, 1, 0)
	if err != nil {
		log.Errorf("Error searching for book with ID %s: %v", bookReq.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if len(hits) == 0 {
		log.Infof("Book with ID %s not found", bookReq.ID)
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	log.Infof("Book with ID %s retrieved successfully", bookReq.ID)
	c.JSON(http.StatusOK, hits[0])
}

func CreateBook(c *gin.Context) {
	bodyBookReq, err := utils.GetValidatedPayload[req.AddBook](c)
	if err != nil {
		log.Errorf("Error getting book by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	book := common.Book{
		ID:          uuid.New(),
		PublishDate: bodyBookReq.PublishDate.Format(time.DateOnly),
	}

	if err := copier.Copy(&book, &bodyBookReq); err != nil {
		log.Errorf("Error copying AddBook payload to Book: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to map request data"})
		return
	}

	clients.EnqueueIndexTask(c, book.ID.String(), book, consts.DoCreateIndex)
	log.Infof("Book with ID %s queued for creation successfully", book.ID)
	c.JSON(http.StatusAccepted, res.AddBook{ID: book.ID})
}

func UpdateBook(c *gin.Context) {
	bodyBookReq, err := utils.GetValidatedPayload[req.UpdateBook](c)
	if err != nil {
		log.Errorf("Error getting book by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	titleUpdate := common.TitleUpdate{}
	if err := copier.Copy(&titleUpdate, &bodyBookReq); err != nil {
		log.Errorf("Error copying UpdateBook payload to TitleUpdate: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to map request data"})
		return
	}

	clients.EnqueueIndexTask(c, bodyBookReq.ID, titleUpdate, consts.DoUpdateIndex)
	log.Infof("Book with ID %s queued for update successfully", bodyBookReq.ID)
	c.JSON(http.StatusAccepted, res.UpdateBook{ID: uuid.MustParse(bodyBookReq.ID)})
}

func DeleteBook(c *gin.Context) {
	deleteReq, err := utils.GetValidatedPayload[req.DeleteBook](c)
	if err != nil {
		log.Errorf("Error getting book by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	clients.EnqueueIndexTask(c, deleteReq.ID, "", consts.DoDeleteIndex)
	log.Infof("Book with ID %s queued for deletion successfully", deleteReq.ID)
	c.JSON(http.StatusAccepted, res.DeleteBook{ID: uuid.MustParse(deleteReq.ID)})
}

func SearchBooks(c *gin.Context) {
	searchReq, err := utils.GetValidatedPayload[req.SearchBooks](c)
	if err != nil {
		log.Errorf("Error getting book by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	esQuery := query.NewQueryBuilder().
		Title(searchReq.Title).
		PriceRange(searchReq.PriceRange.Min, searchReq.PriceRange.Max).
		Build()

	hits, _, err := clients.SearchIndex(c, esQuery, searchReq.Size, searchReq.From)
	if err != nil {
		log.Errorf("Error searching books: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	log.Infof("SearchBooks esQuery executed successfully, retrieved %d results", len(hits))
	c.JSON(http.StatusOK, hits)
}

func GetBooksStats(c *gin.Context) {
	esQuery := query.NewQueryBuilder().
		DistinctAuthors().
		Build()

	_, aggregations, err := clients.SearchIndex(c, esQuery, 0, 0, clients.EsClient.Search.WithTrackTotalHits(true))
	if err != nil {
		log.Errorf("Error fetching books statistics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	stats, err := utils.ParseAggregations(aggregations, consts.AggregationConfigs["BookStats"])
	if err != nil {
		log.Errorf("Error parsing aggregations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	log.Infof("GetBooksStats executed successfully: %+v", stats)
	c.JSON(http.StatusOK, stats)
}
