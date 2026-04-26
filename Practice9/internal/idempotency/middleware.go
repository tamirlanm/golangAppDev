package idempotency

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func Middleware(store Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Idempotency-Key header required", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read request body", http.StatusBadRequest)
			return
		}
		_ = r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(body))

		hash := sha256.Sum256(body)
		payloadHash := hex.EncodeToString(hash[:])

		cached, started, err := store.TryStart(r.Context(), key, payloadHash)
		if err != nil {
			switch {
			case errors.Is(err, ErrInProgress):
				http.Error(w, "Duplicate request in progress", http.StatusConflict)
			case errors.Is(err, ErrPayloadMismatch):
				http.Error(w, "Idempotency-Key reused with a different payload", http.StatusConflict)
			default:
				http.Error(w, fmt.Sprintf("storage error: %v", err), http.StatusInternalServerError)
			}
			return
		}

		if cached != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(cached.StatusCode)
			_, _ = w.Write(cached.Body)
			return
		}

		if !started {
			http.Error(w, "unexpected idempotency state", http.StatusInternalServerError)
			return
		}

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		if err := store.Finish(r.Context(), key, recorder.Code, recorder.Body.Bytes()); err != nil {
			http.Error(w, fmt.Sprintf("cannot save result: %v", err), http.StatusInternalServerError)
			return
		}

		copyHeaders(w.Header(), recorder.Header())
		w.WriteHeader(recorder.Code)
		_, _ = w.Write(recorder.Body.Bytes())
	})
}

func copyHeaders(dst, src http.Header) {
	for k, vals := range src {
		for _, v := range vals {
			dst.Add(k, v)
		}
	}
}
