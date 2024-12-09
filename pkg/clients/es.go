package clients

import (
	"book_service/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

var (
	esClient *elasticsearch.Client
)

func InitElasticsearchClient() error {
	tr := &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	elasticUri, _ := utils.GetEnvVar[string]("ELS_URI", "")
	println(elasticUri)
	cfg := elasticsearch.Config{
		Addresses:     []string{elasticUri},
		Transport:     tr,
		RetryOnStatus: []int{502, 503, 504},
		MaxRetries:    3,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	res, err := client.Info()
	if err != nil {
		return fmt.Errorf("failed to get Elasticsearch info: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("elasticsearch returned error: %s", res.String())
	}

	esClient = client
	log.Println("Elasticsearch client initialized successfully")
	return nil
}

// =====================
// WORKER POOL SETUP
// =====================

type SearchRequest struct {
	Ctx          context.Context
	Index        string
	Query        interface{}
	ResponseChan chan *SearchResult
}

type SearchResult struct {
	Response *esapi.Response
	Err      error
}

type IndexRequest struct {
	Ctx          context.Context
	Index        string
	ID           string
	Document     interface{}
	ResponseChan chan *IndexResult
}

type IndexResult struct {
	Response *esapi.Response
	Err      error
}

var (
	taskQueueIndex chan IndexRequest
	workersDone    chan struct{}
	oncePool       sync.Once
)

func InitWorkerPool(numWorkers int) {
	oncePool.Do(func() {
		taskQueueIndex = make(chan IndexRequest, 1000)
		workersDone = make(chan struct{})

		for i := 0; i < numWorkers; i++ {
			go indexWorker(taskQueueIndex, workersDone)
		}
		log.Printf("Started %d Elasticsearch search worker(s) and %d index worker(s)", numWorkers, numWorkers)
	})
}

func ShutdownWorkerPool() {
	close(taskQueueIndex)
	for range workersDone {
		// consume until all workers are done
	}
}

func indexWorker(tasks <-chan IndexRequest, done chan<- struct{}) {
	for req := range tasks {
		res, err := doIndex(req.Ctx, req.Index, req.ID, req.Document)
		req.ResponseChan <- &IndexResult{Response: res, Err: err}
		close(req.ResponseChan)
	}
	done <- struct{}{}
}

func DoSearch(ctx context.Context, index string, query interface{}) ([]interface{}, error) {
	if esClient == nil {
		return nil, errors.New("elasticsearch client not initialized")
	}

	//queryBuild := map[string]interface{}{
	//	"query": map[string]interface{}{
	//		"term": map[string]interface{}{
	//			"_id": query,
	//		},
	//	},
	//}

	res, err := esClient.Search(
		esClient.Search.WithContext(ctx),
		esClient.Search.WithIndex(index),
		esClient.Search.WithBody(esutil.NewJSONReader(query)),
		esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	if res.IsError() {
		return nil, fmt.Errorf("search returned error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing response body: %s", err)
	}

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	fmt.Printf("Number of hits: %d\n", len(hits))
	for _, hit := range hits {
		doc := hit.(map[string]interface{})["_source"]
		fmt.Printf("Document: %v\n", doc)
	}
	return hits, nil
}

func doIndex(ctx context.Context, index, id string, doc interface{}) (*esapi.Response, error) {
	fmt.Printf("Document to index: %+v\n", doc)

	if esClient == nil {
		return nil, errors.New("elasticsearch client not initialized")
	}

	if doc == nil {
		return nil, fmt.Errorf("document cannot be nil")
	}

	res, err := esClient.Index(
		index,
		esutil.NewJSONReader(doc),
		esClient.Index.WithContext(ctx),
		esClient.Index.WithDocumentID(id),
	)
	if err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("index request failed: %w", err)
	}

	if res.IsError() {
		logrus.Error(res.String())
		return res, fmt.Errorf("index returned error: %s", res.String())
	}

	return res, nil
}

// =====================
// PUBLIC FUNCTIONS
// =====================

// EnqueueSearchTask
// shouldn't be used due to latency, use doSearch
//func EnqueueSearchTask(ctx context.Context, index string, query interface{}) <-chan *SearchResult {
//	responseChan := make(chan *SearchResult, 1)
//	req := SearchRequest{
//		Ctx:          ctx,
//		Index:        index,
//		Query:        query,
//		ResponseChan: responseChan,
//	}
//	taskQueueSearch <- req
//	return responseChan
//}

func EnqueueIndexTask(ctx context.Context, index, id string, document interface{}) <-chan *IndexResult {
	responseChan := make(chan *IndexResult, 1)
	req := IndexRequest{
		Ctx:          ctx,
		Index:        index,
		ID:           id,
		Document:     document,
		ResponseChan: responseChan,
	}
	taskQueueIndex <- req
	return responseChan
}
