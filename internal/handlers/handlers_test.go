package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"exoplanet-microservice/internal/models"
	"exoplanet-microservice/internal/repository"
	"exoplanet-microservice/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setup() *mux.Router {
	db := initTestDB()
	svc := service.NewExoplanetService(repository.NewSQLiteExoplanetRepository(db))
	handler := NewExoplanetHandler(svc)

	r := mux.NewRouter()
	r.HandleFunc("/exoplanets", handler.AddExoplanet).Methods("POST")
	r.HandleFunc("/exoplanets", handler.ListExoplanets).Methods("GET")
	r.HandleFunc("/exoplanets/{id}", handler.GetExoplanetByID).Methods("GET")
	r.HandleFunc("/exoplanets/{id}", handler.UpdateExoplanet).Methods("PUT")
	r.HandleFunc("/exoplanets/{id}", handler.DeleteExoplanet).Methods("DELETE")
	r.HandleFunc("/exoplanets/{id}/fuel-estimation", handler.FuelEstimation).Methods("GET")
	return r
}

func initTestDB() *sql.DB {
	// Initialize SQLite database for testing
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}

	// Create table if not exists
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS exoplanets (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			distance INTEGER NOT NULL,
			radius REAL NOT NULL,
			mass REAL,
			type TEXT NOT NULL
		);
	`
	_, err = db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}

	return db
}

func TestAddExoplanet(t *testing.T) {
	router := setup()

	payload := []byte(`{
		"name": "TestPlanet",
		"description": "A test gas giant",
		"distance": 300,
		"radius": 2.0,
		"type": "GasGiant"
	}`)

	req, _ := http.NewRequest("POST", "/exoplanets", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusCreated, response.Code)
}

func TestListExoplanets(t *testing.T) {
	router := setup()

	req, _ := http.NewRequest("GET", "/exoplanets", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)
	var planets []models.Exoplanet
	json.NewDecoder(response.Body).Decode(&planets)
	assert.GreaterOrEqual(t, len(planets), 0)
}

func TestGetExoplanetByID(t *testing.T) {
	router := setup()

	// First, add a planet to retrieve
	payload := []byte(`{
		"name": "TestPlanet",
		"description": "A test gas giant",
		"distance": 300,
		"radius": 2.0,
		"type": "GasGiant"
	}`)

	req, _ := http.NewRequest("POST", "/exoplanets", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	// Now, get the planet by ID
	var result models.Exoplanet
	json.NewDecoder(response.Body).Decode(&result)

	req, _ = http.NewRequest("GET", "/exoplanets/"+result.ID, nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestUpdateExoplanet(t *testing.T) {
	router := setup()

	// First, add a planet to update
	payload := []byte(`{
		"name": "TestPlanet",
		"description": "A test gas giant",
		"distance": 300,
		"radius": 2.0,
		"type": "GasGiant"
	}`)

	req, _ := http.NewRequest("POST", "/exoplanets", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	// Now, update the planet by ID
	var result models.Exoplanet
	json.NewDecoder(response.Body).Decode(&result)

	updatePayload := []byte(`{
		"name": "UpdatedPlanet",
		"description": "An updated description",
		"distance": 350,
		"radius": 2.5,
		"type": "GasGiant"
	}`)

	req, _ = http.NewRequest("PUT", "/exoplanets/"+result.ID, bytes.NewBuffer(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestDeleteExoplanet(t *testing.T) {
	router := setup()

	// First, add a planet to delete
	payload := []byte(`{
		"name": "TestPlanet",
		"description": "A test gas giant",
		"distance": 300,
		"radius": 2.0,
		"type": "GasGiant"
	}`)

	req, _ := http.NewRequest("POST", "/exoplanets", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	// Now, delete the planet by ID
	var result models.Exoplanet
	json.NewDecoder(response.Body).Decode(&result)

	req, _ = http.NewRequest("DELETE", "/exoplanets/"+result.ID, nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestFuelEstimation(t *testing.T) {
	router := setup()

	// First, add a planet to estimate fuel for
	payload := []byte(`{
		"name": "TestPlanet",
		"description": "A test gas giant",
		"distance": 300,
		"radius": 2.0,
		"type": "GasGiant"
	}`)

	req, _ := http.NewRequest("POST", "/exoplanets", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	// Now, estimate fuel for the planet by ID
	var result models.Exoplanet
	json.NewDecoder(response.Body).Decode(&result)

	req, _ = http.NewRequest("GET", "/exoplanets/"+result.ID+"/fuel-estimation?crew_capacity=10", nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)
	var fuelEstimation map[string]interface{}
	json.NewDecoder(response.Body).Decode(&fuelEstimation)
	assert.NotNil(t, fuelEstimation["fuel"])
}
