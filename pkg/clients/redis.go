package clients

import (
	_const "book_service/pkg/constants"
	log "github.com/sirupsen/logrus"
	"book_service/pkg/utils"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type UserAction struct {
	Method string    `json:"method"`
	Path   string    `json:"path"`
	Time   time.Time `json:"time"`
	User   string    `json:"user"`
}

var (
	redisClient   *redis.Client
	actionBuffers = make(map[string][]UserAction)
	bufferMutex   sync.Mutex
	actionsChan   = make(chan UserAction, _const.ActionsChanelSize)
	flushInterval = _const.FlushInterval
	flushSize     = _const.FlushSize
)

func InitRedisClient() {
	redisUri, _ := utils.GetEnvVar[string]("REDIS_URI", "localhost:6379")
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisUri,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	go actionWorker()
}

func AppendAction(action UserAction) {
	select {
	case actionsChan <- action:
	default:
		log.Println("Warning: actions channel is full, dropping action")
	}
}

func GetLastActions(user string) ([]UserAction, error) {
	vals, err := redisClient.LRange(context.Background(), "user:"+user+":actions", 0, -1).Result()
	if err != nil {
		log.Errorf("Was unable to get user action for %s with error: %v", user, err)
		return nil, err
	}

	var actions []UserAction
	var userAction UserAction
	for _, val := range vals {
		err := json.Unmarshal([]byte(val), &userAction)
		if err != nil {
			log.Errorf("Was unable to unmarshal actions for %s with error: %v", user, err)
			return nil, err
		}
		actions = append(actions, userAction)
	}

	return actions, nil
}

func actionWorker() {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case action, ok := <-actionsChan:
			if !ok {
				flushAllToRedis()
				return
			}
			addToBuffer(action)
		case <-ticker.C:
			flushAllToRedis()
		}
	}
}

func addToBuffer(action UserAction) {
	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	buffer, exists := actionBuffers[action.User]
	if !exists {
		buffer = []UserAction{}
	}
	buffer = append(buffer, action)
	actionBuffers[action.User] = buffer

	if len(buffer) >= flushSize {
		flushAllToRedis()
	}
}

func flushAllToRedis() {
	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	if len(actionBuffers) == 0 {
		return
	}

	pipe := redisClient.Pipeline()
	for user, actions := range actionBuffers {
		flushUserActionsToPipeline(pipe, user, actions)
		delete(actionBuffers, user)
	}

	_, err := pipe.Exec(context.Background())
	if err != nil {
		log.Printf("Failed to execute pipeline: %v", err)
	}
}

func flushUserActionsToPipeline(pipe redis.Pipeliner, user string, actions []UserAction) {
	for _, action := range actions {
		data, err := json.Marshal(action)
		if err != nil {
			log.Printf("Failed to marshal action for user %s: %v", user, err)
			continue
		}
		pipe.LPush(context.Background(), "user:"+user+":actions", data)
		pipe.LTrim(context.Background(), "user:"+user+":actions", 0, 2) // Keep only the 3 most recent elements
	}
}
