package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"Practice9/internal/idempotency"
	"Practice9/internal/retry"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	var err error
	switch os.Args[1] {
	case "retry":
		err = runRetryDemo()
	case "idempotency":
		err = runIdempotencyDemo()
	case "serve":
		err = runServer()
	case "all":
		err = runRetryDemo()
		if err == nil {
			err = runIdempotencyDemo()
		}
	default:
		printUsage()
		return
	}

	if err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run . retry        # task 1 demo")
	fmt.Println("  go run . idempotency  # task 2 demo")
	fmt.Println("  go run . serve        # real HTTP server on :8080")
	fmt.Println("  go run . all          # run both demos one after another")
}

func runRetryDemo() error {
	log.Println("--- Retry demo started ---")

	var hits int32
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt32(&hits, 1)
		log.Printf("Gateway request #%d", current)

		type response struct {
			Status string `json:"status"`
		}

		if current <= 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "service unavailable"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response{Status: "success"})
	}))
	defer testServer.Close()

	client := retry.NewClient(5)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body, status, err := client.ExecutePayment(ctx, testServer.URL, []byte(`{"amount":1000}`))
	if err != nil {
		return err
	}

	log.Printf("Final status: %d", status)
	log.Printf("Final body: %s", string(body))
	log.Println("--- Retry demo finished ---")
	return nil
}

func runIdempotencyDemo() error {
	log.Println("--- Idempotency demo started ---")

	store, err := idempotency.NewRedisStore("localhost:6379")
	if err != nil {
		return err
	}
	defer func() {
		_ = store.Close()
	}()

	handler := idempotency.Middleware(store, http.HandlerFunc(idempotency.PaymentHandler))
	server := httptest.NewServer(handler)
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	key := "demo-key-123"
	payload := []byte(`{"amount":1000}`)

	var wg sync.WaitGroup
	start := make(chan struct{})
	requestCount := 8

	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-start

			req, err := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(payload))
			if err != nil {
				log.Printf("request %d build error: %v", id, err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Idempotency-Key", key)

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("request %d error: %v", id, err)
				return
			}
			defer resp.Body.Close()

			data, _ := io.ReadAll(resp.Body)
			log.Printf("request %d -> %d %s", id, resp.StatusCode, string(data))
		}(i + 1)
	}

	close(start)
	wg.Wait()

	req, err := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", key)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	log.Printf("repeat request -> %d %s", resp.StatusCode, string(data))
	log.Println("--- Idempotency demo finished ---")
	return nil
}

func runServer() error {
	store, err := idempotency.NewRedisStore("localhost:6379")
	if err != nil {
		return err
	}
	defer func() {
		_ = store.Close()
	}()

	mux := http.NewServeMux()
	mux.Handle("/pay", idempotency.Middleware(store, http.HandlerFunc(idempotency.PaymentHandler)))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Println("Server running on http://localhost:8080")
	return http.ListenAndServe(":8080", mux)
}
