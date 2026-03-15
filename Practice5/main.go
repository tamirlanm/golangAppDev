package main

import (
	"log"
	"net/http"
)

func main() {
	db, err := OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := NewRepository(db)
	handler := NewHandler(repo)

	RegisterRoutes(handler)

	log.Println("server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
