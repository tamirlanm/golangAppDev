package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Movie struct {
	ID     string   `json:"id"`
	Genre  string   `json:"genre"`
	Budget int      `json:"budget"`
	Title  string   `json:"title"`
	Actors []string `json:"actors"`
}

var db *sql.DB

func connectDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var conn *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		conn, err = sql.Open("postgres", dsn)
		if err == nil {
			err = conn.Ping()
		}
		if err == nil {
			log.Println("Connected to database!")
			return conn
		}
		log.Printf("DB not ready, retrying in 2s... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Could not connect to DB: %v", err)
	return nil
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, genre, budget, title, actors FROM movies")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		var actors []byte
		if err := rows.Scan(&m.ID, &m.Genre, &m.Budget, &m.Title, &actors); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.Unmarshal(actors, &m.Actors)
		movies = append(movies, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var m Movie
	var actors []byte
	err := db.QueryRow(
		"SELECT id, genre, budget, title, actors FROM movies WHERE id=$1", id,
	).Scan(&m.ID, &m.Genre, &m.Budget, &m.Title, &actors)

	if err == sql.ErrNoRows {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Unmarshal(actors, &m.Actors)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var m Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actors, _ := json.Marshal(m.Actors)
	_, err := db.Exec(
		"INSERT INTO movies (id, genre, budget, title, actors) VALUES ($1,$2,$3,$4,$5)",
		m.ID, m.Genre, m.Budget, m.Title, actors,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var m Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actors, _ := json.Marshal(m.Actors)
	res, err := db.Exec(
		"UPDATE movies SET genre=$1, budget=$2, title=$3, actors=$4 WHERE id=$5",
		m.Genre, m.Budget, m.Title, actors, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	m.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	res, err := db.Exec("DELETE FROM movies WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	db = connectDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down gracefully...")
		db.Close()
		os.Exit(0)
	}()

	log.Println("Starting the Server on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
