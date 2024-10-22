package logger

import (
	"Cloud/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// RequestLogger представляет собой структуру для логирования запросов
type RequestLogger struct {
	collection *mongo.Collection
}

func NewRequestLogger(client *mongo.Client, dbName, collectionName string) *RequestLogger {
	collection := client.Database(dbName).Collection(collectionName)
	return &RequestLogger{collection: collection}
}

// Log записывает лог запроса в MongoDB
func (rl *RequestLogger) Log(method, endpoint, userID, ip, userAgent string, statusCode int, duration time.Duration) error {
	logEntry := models.RequestLog{
		Method:     method,
		Endpoint:   endpoint,
		UserID:     userID,
		IP:         ip,
		UserAgent:  userAgent,
		Time:       time.Now(),
		StatusCode: statusCode,
		Duration:   time.Duration(duration.Milliseconds()),
	}

	_, err := rl.collection.InsertOne(context.Background(), logEntry)
	return err
}
