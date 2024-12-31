package clients

import (
	_const "book_service/pkg/constants"
	"book_service/pkg/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

var BooksIndex, _ = utils.GetEnvVar[string]("BOOKS_INDEX", "books")

type IndexRequest struct {
	Ctx          context.Context
	Index        string
	ID           string
	Document     interface{}
	ResponseChan chan *IndexResult
	_const.Function
}

type IndexResult struct {
	Response *esapi.Response
	Err      error
}

var (
	esClient       *elasticsearch.Client
	taskQueueIndex chan IndexRequest
	workersDone    chan struct{}
	oncePool       sync.Once
)

func InitElasticsearchClient() error {
	tr := &http.Transport{
		MaxIdleConns:        _const.MaxIdleConnections,
		MaxIdleConnsPerHost: _const.MaxIdleConnectionsPerHost,
		IdleConnTimeout:     _const.IdleConnectionTimeout,
		DialContext: (&net.Dialer{
			Timeout:   _const.DialTimeout,
			KeepAlive: _const.KeepAliveTime,
		}).DialContext,
		TLSHandshakeTimeout:   _const.TLSHandshakeTimeout,
		ExpectContinueTimeout: _const.ExpectContinueTimeout,
	}

	elasticUri, _ := utils.GetEnvVar[string]("ELS_URI", "")
	cfg := elasticsearch.Config{
		Addresses:     []string{elasticUri},
		Transport:     tr,
		RetryOnStatus: []int{502, 503, 504},
		MaxRetries:    3,
	}

	client, err := elasticsearch.NewClient(cfg) // TODO: V8 -> elasticsearch.NewTypedClient(cfg)
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
	InitializeIndices()
	log.Info("Setup indexes successfully")
	log.Info("Elasticsearch client initialized successfully")
	return nil
}

func InitElasticWorkerPool(numWorkers int) {
	oncePool.Do(func() {
		taskQueueIndex = make(chan IndexRequest, 1000)
		workersDone = make(chan struct{})

		for i := 0; i < numWorkers; i++ {
			go indexWorker(taskQueueIndex, workersDone)
		}
		log.Infof("Started %d Elasticsearch index worker(s)", numWorkers)
	})
}

func ShutdownWorkerPool() {
	close(taskQueueIndex)
	for range workersDone {
	}
}

func EnqueueIndexTask(ctx context.Context, index, id string, document interface{}, function _const.Function) {
	responseChan := make(chan *IndexResult, 1)
	req := IndexRequest{
		Ctx:          ctx,
		Index:        index,
		ID:           id,
		Document:     document,
		ResponseChan: responseChan,
		Function:     function,
	}
	taskQueueIndex <- req
	log.Infof("Task enqueued for %s in index %s", function, index)
}

func SearchIndex(ctx context.Context, index string, query interface{}, size, from int) ([]map[string]interface{}, map[string]interface{}, error) {
	if esClient == nil {
		return nil, nil, errors.New("elasticsearch client not initialized")
	}

	res, err := esClient.Search(
		esClient.Search.WithContext(ctx),
		esClient.Search.WithIndex(index),
		esClient.Search.WithBody(esutil.NewJSONReader(query)),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithSize(size), 
		esClient.Search.WithFrom(from),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("search request failed: %w", err)
	}

	if res.IsError() {
		return nil, nil, fmt.Errorf("search returned error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing response body: %v", err)
		return nil, nil, err
	}

	hitsArray, ok := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		log.Fatalf("Error: Unable to extract hits from response")
		return nil, nil, fmt.Errorf("unable to extract hits from response")
	}

	hits := make([]map[string]interface{}, 0, len(hitsArray))
	for _, hit := range hitsArray {
		if hitMap, ok := hit.(map[string]interface{}); ok {
			if source, ok := hitMap["_source"].(map[string]interface{}); ok {
				hits = append(hits, source)
			}
		}
	}

	aggregations, ok := r["aggregations"].(map[string]interface{})
	if !ok {
		aggregations = nil
	}

	return hits, aggregations, nil
}

// Internal Functions
func InitializeIndices() {
	for _, indexMapping := range _const.IndexMappings {
		err := createIndex(esClient, indexMapping.IndexName, indexMapping.Mapping)
		if err != nil {
			log.Errorf("Failed to create index %s: %v", indexMapping.IndexName, err)
		} else {
			log.Infof("Index %s initialized successfully", indexMapping.IndexName)
		}
	}
}

func createIndex(client *elasticsearch.Client, index string, mapping string) error {
	exists, err := client.Indices.Exists([]string{index})
	if err != nil {
		log.Errorf("Error checking if index %s exists: %v", index, err)
		return err
	}
	defer exists.Body.Close()

	if exists.StatusCode == 200 {
		log.Infof("Index %s already exists", index)
		return nil
	}

	res, err := client.Indices.Create(index, client.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		log.Errorf("Error creating index %s: %v", index, err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Errorf("Error response from Elasticsearch for index %s: %s", index, res.String())
		return fmt.Errorf("failed to create index %s: %s", index, res.String())
	}

	log.Infof("Index %s created successfully", index)
	return nil
}

func indexWorker(tasks <-chan IndexRequest, done chan<- struct{}) {
	defer func() { done <- struct{}{} }()
	for req := range tasks {
		var (
			res *esapi.Response
			err error
		)

		switch req.Function {
		case _const.DoCreateIndex:
			res, err = addToIndex(req.Ctx, req.Index, req.ID, req.Document)
		case _const.DoUpdateIndex:
			err = updateIndex(req.Index, req.ID, req.Document)
		case _const.DoDeleteIndex:
			err = deleteIndex(req.Index, req.ID)
		default:
			err = fmt.Errorf("invalid function type: %d", req.Function)
		}

		req.ResponseChan <- &IndexResult{Response: res, Err: err}
		close(req.ResponseChan)
	}

}

func addToIndex(ctx context.Context, index, id string, doc interface{}) (*esapi.Response, error) {
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
		log.Errorf("Index request failed: %v", err)
		return nil, fmt.Errorf("index request failed: %w", err)
	}

	if res.IsError() {
		log.Errorf("Index returned error: %s", res.String())
		return res, fmt.Errorf("index returned error: %s", res.String())
	}

	log.Infof("Document %s indexed successfully in %s", id, index)
	return res, nil
}

func deleteIndex(index, docID string) error {
	res, err := esClient.Delete(index, docID)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting document: %s", res.String())
	}

	log.Infof("Document %s deleted successfully from index %s", docID, index)
	return nil
}

func updateIndex(index, docID string, updateData interface{}) error {
	updateBody, err := json.Marshal(map[string]interface{}{
		"doc": updateData,
	})
	if err != nil {
		return fmt.Errorf("error marshalling update data: %w", err)
	}

	res, err := esClient.Update(index, docID, bytes.NewReader(updateBody))
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error updating document: %s", res.String())
	}

	log.Infof("Document %s updated successfully in index %s", docID, index)
	return nil
}
