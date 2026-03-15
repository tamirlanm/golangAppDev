package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	query := r.URL.Query()

	page := parseIntOrDefault(query.Get("page"), 1)
	pageSize := parseIntOrDefault(query.Get("page_size"), 5)

	filter := UserFilter{
		Page:     page,
		PageSize: pageSize,
		OrderBy:  query.Get("order_by"),
	}

	if value := strings.TrimSpace(query.Get("id")); value != "" {
		id, err := strconv.Atoi(value)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
			return
		}
		filter.ID = &id
	}

	if value := strings.TrimSpace(query.Get("name")); value != "" {
		filter.Name = &value
	}

	if value := strings.TrimSpace(query.Get("email")); value != "" {
		filter.Email = &value
	}

	if value := strings.TrimSpace(query.Get("gender")); value != "" {
		filter.Gender = &value
	}

	if value := strings.TrimSpace(query.Get("birth_date")); value != "" {
		birthDate, err := time.Parse("2006-01-02", value)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "birth_date must be in YYYY-MM-DD format"})
			return
		}
		filter.BirthDate = &birthDate
	}

	result, err := h.repo.GetPaginatedUsers(filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	query := r.URL.Query()

	user1ID, err := strconv.Atoi(query.Get("user1_id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid user1_id"})
		return
	}

	user2ID, err := strconv.Atoi(query.Get("user2_id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid user2_id"})
		return
	}

	if user1ID == user2ID {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "user1_id and user2_id must be different"})
		return
	}

	friends, err := h.repo.GetCommonFriends(user1ID, user2ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, friends)
}

func parseIntOrDefault(value string, defaultValue int) int {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}

	number, err := strconv.Atoi(value)
	if err != nil || number <= 0 {
		return defaultValue
	}

	return number
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
