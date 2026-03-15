package main

import "net/http"

func RegisterRoutes(h *Handler) {
	http.HandleFunc("/users", h.GetUsers)
	http.HandleFunc("/common-friends", h.GetCommonFriends)
}
