package main

import (
	"database/sql"
	"exoplanet-microservice/internal/handlers"
	"exoplanet-microservice/internal/repository"
	"exoplanet-microservice/internal/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect to SQLite database
	db, err := sql.Open("sqlite3", "./exoplanets.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create exoplanets table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS exoplanets (
		id TEXT PRIMARY KEY,
		name TEXT,
		description TEXT,
		distance INTEGER,
		radius REAL,
		mass REAL,
		type TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Set up repository, service, and handlers
	repo := repository.NewSQLiteExoplanetRepository(db)
	service := service.NewExoplanetService(repo)
	handler := handlers.NewExoplanetHandler(service)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/exoplanets", handler.AddExoplanet).Methods("POST")
	r.HandleFunc("/exoplanets", handler.ListExoplanets).Methods("GET")
	r.HandleFunc("/exoplanets/{id}", handler.GetExoplanetByID).Methods("GET")
	r.HandleFunc("/exoplanets/{id}", handler.UpdateExoplanet).Methods("PUT")
	r.HandleFunc("/exoplanets/{id}", handler.DeleteExoplanet).Methods("DELETE")
	r.HandleFunc("/exoplanets/{id}/fuel-estimation", handler.FuelEstimation).Methods("GET")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
