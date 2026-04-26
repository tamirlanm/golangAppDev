package idempotency

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrInProgress      = errors.New("idempotency request still processing")
	ErrPayloadMismatch = errors.New("idempotency key reused with different payload")
)

type CachedResponse struct {
	StatusCode int
	Body       []byte
}

type Store interface {
	TryStart(ctx context.Context, key, payloadHash string) (*CachedResponse, bool, error)
	Finish(ctx context.Context, key string, statusCode int, body []byte) error
	Close() error
}

type RedisStore struct {
	client *redis.Client
}

type storedValue struct {
	Status      string `json:"status"` // processing | completed
	StatusCode  int    `json:"status_code"`
	Body        string `json:"body"`
	PayloadHash string `json:"payload_hash"`
}

func NewRedisStore(addr string) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisStore{client: client}, nil
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}

func (s *RedisStore) TryStart(ctx context.Context, key, payloadHash string) (*CachedResponse, bool, error) {
	processingKey := "idempotency:" + key

	current := storedValue{
		Status:      "processing",
		PayloadHash: payloadHash,
	}

	data, err := json.Marshal(current)
	if err != nil {
		return nil, false, err
	}

	// Atomic: create key only if it does not exist, with TTL.
	// Redis SET key value NX EX is the standard atomic gate for this pattern.
	ok, err := s.client.SetNX(ctx, processingKey, data, 5*time.Minute).Result()
	if err != nil {
		return nil, false, err
	}

	if ok {
		// We are the first request.
		return nil, true, nil
	}

	// Key already exists. Read current value.
	raw, err := s.client.Get(ctx, processingKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, ErrInProgress
		}
		return nil, false, err
	}

	var existing storedValue
	if err := json.Unmarshal(raw, &existing); err != nil {
		return nil, false, err
	}

	if existing.PayloadHash != payloadHash {
		return nil, false, ErrPayloadMismatch
	}

	if existing.Status == "completed" {
		return &CachedResponse{
			StatusCode: existing.StatusCode,
			Body:       []byte(existing.Body),
		}, false, nil
	}

	return nil, false, ErrInProgress
}

func (s *RedisStore) Finish(ctx context.Context, key string, statusCode int, body []byte) error {
	processingKey := "idempotency:" + key

	raw, err := s.client.Get(ctx, processingKey).Bytes()
	if err != nil {
		return err
	}

	var current storedValue
	if err := json.Unmarshal(raw, &current); err != nil {
		return err
	}

	current.Status = "completed"
	current.StatusCode = statusCode
	current.Body = string(body)

	updated, err := json.Marshal(current)
	if err != nil {
		return err
	}

	// Save final response with a longer TTL.
	return s.client.Set(ctx, processingKey, updated, 24*time.Hour).Err()
}
