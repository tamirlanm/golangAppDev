package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetRate_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"base":"USD","target":"EUR","rate":0.92}`))
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)

	rate, err := svc.GetRate("USD", "EUR")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rate != 0.92 {
		t.Errorf("rate = %v; want %v", rate, 0.92)
	}
}

func TestGetRate_APIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid currency pair"}`))
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)

	_, err := svc.GetRate("AAA", "BBB")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetRate_MalformedJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"base":"USD","target":"EUR"`))
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)

	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}

func TestGetRate_EmptyBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)

	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}

func TestGetRate_UnexpectedStatusWithoutErrorMsg(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"base":"USD","target":"EUR","rate":0}`))
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)

	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetRate_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(6 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"base":"USD","target":"EUR","rate":0.92}`))
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)
	svc.Client.Timeout = 1 * time.Second

	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestGetRate_BusinessError404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"invalid currency pair"}`))
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)

	_, err := svc.GetRate("X", "Y")
	if err == nil {
		t.Fatal("expected api error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid currency pair") {
		t.Errorf("unexpected error: %v", err)
	}
}
