package retry

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"errors"
	"fmt"
	"io"
	mathrand "math/rand"
	"net/http"
	"time"
)

const (
	defaultMaxRetries = 5
	baseDelay         = 500 * time.Millisecond
	maxDelay          = 5 * time.Second
)

type Client struct {
	httpClient *http.Client
	maxRetries int
	rng        *mathrand.Rand
}

func NewClient(maxRetries int) *Client {
	if maxRetries <= 0 {
		maxRetries = defaultMaxRetries
	}

	return &Client{
		httpClient: http.DefaultClient,
		maxRetries: maxRetries,
		rng:        mathrand.New(mathrand.NewSource(time.Now().UnixNano())),
	}
}

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return false
		}

		var timeout interface{ Timeout() bool }
		if errors.As(err, &timeout) && timeout.Timeout() {
			return true
		}

		return true
	}

	if resp == nil {
		return false
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	case http.StatusUnauthorized, http.StatusNotFound:
		return false
	default:
		return false
	}
}

func CalculateBackoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	backoff := baseDelay
	for i := 0; i < attempt; i++ {
		if backoff >= maxDelay/2 {
			backoff = maxDelay
			break
		}
		backoff *= 2
	}

	if backoff > maxDelay {
		backoff = maxDelay
	}
	if backoff <= 0 {
		return 0
	}

	seed, err := randomInt64()
	if err != nil {
		return backoff / 2
	}
	if seed < 0 {
		seed = -seed
	}
	return time.Duration(seed % int64(backoff))
}

func (c *Client) ExecutePayment(ctx context.Context, endpoint string, payload []byte) ([]byte, int, error) {
	if c == nil {
		return nil, 0, fmt.Errorf("retry client is nil")
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	var lastErr error
	var lastStatus int

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, 0, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
		if err != nil {
			return nil, 0, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if resp != nil {
			body, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr != nil {
				return nil, 0, readErr
			}

			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return body, resp.StatusCode, nil
			}

			lastStatus = resp.StatusCode
			lastErr = fmt.Errorf("http %d: %s", resp.StatusCode, string(body))

			if !IsRetryable(resp, err) {
				return body, resp.StatusCode, lastErr
			}
		} else if err != nil {
			lastErr = err
			if !IsRetryable(nil, err) {
				return nil, 0, err
			}
		}

		if attempt == c.maxRetries-1 {
			break
		}

		delay := CalculateBackoff(attempt)
		fmt.Printf("Attempt %d failed: waiting %v...\n", attempt+1, delay)

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, 0, ctx.Err()
		case <-timer.C:
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("payment failed after %d attempts", c.maxRetries)
	}

	return nil, lastStatus, lastErr
}

func randomInt64() (int64, error) {
	var b [8]byte
	if _, err := crand.Read(b[:]); err != nil {
		return 0, err
	}

	var v int64
	for _, x := range b {
		v = (v << 8) | int64(x)
	}
	return v, nil
}
