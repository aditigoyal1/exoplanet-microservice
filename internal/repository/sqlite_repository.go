package repository

import (
	"database/sql"
	"errors"
	"exoplanet-microservice/internal/models"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var once sync.Once

type SQLiteExoplanetRepository struct {
	DB *sql.DB
}

// NewSQLiteExoplanetRepository initializes a new instance of SQLiteExoplanetRepository.
func NewSQLiteExoplanetRepository(db *sql.DB) *SQLiteExoplanetRepository {
	return &SQLiteExoplanetRepository{DB: db}
}

// InitDB initializes the SQLite database connection.
func InitDB() (*sql.DB, error) {
	var db *sql.DB
	var err error
	once.Do(func() {
		db, err = sql.Open("sqlite3", "file:exoplanets.db?cache=shared&mode=rwc")
		if err != nil {
			err = errors.New("failed to open database connection")
			return
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
			err = errors.New("failed to create table exoplanets")
			return
		}
	})

	return db, err
}

func (r *SQLiteExoplanetRepository) AddExoplanet(exoplanet models.Exoplanet) error {
	_, err := r.DB.Exec("INSERT INTO exoplanets (id, name, description, distance, radius, mass, type) VALUES (?, ?, ?, ?, ?, ?, ?)",
		exoplanet.ID, exoplanet.Name, exoplanet.Description, exoplanet.Distance, exoplanet.Radius, exoplanet.Mass, exoplanet.Type)
	if err != nil {
		return errors.New("failed to insert exoplanet")
	}
	return nil
}

func (r *SQLiteExoplanetRepository) ListExoplanets() ([]models.Exoplanet, error) {
	rows, err := r.DB.Query("SELECT id, name, description, distance, radius, mass, type FROM exoplanets")
	if err != nil {
		return nil, errors.New("failed to fetch exoplanets")
	}
	defer rows.Close()

	var exoplanets []models.Exoplanet
	for rows.Next() {
		var exoplanet models.Exoplanet
		err := rows.Scan(&exoplanet.ID, &exoplanet.Name, &exoplanet.Description, &exoplanet.Distance, &exoplanet.Radius, &exoplanet.Mass, &exoplanet.Type)
		if err != nil {
			return nil, err
		}
		exoplanets = append(exoplanets, exoplanet)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("error while iterating over exoplanets")
	}

	return exoplanets, nil
}

func (r *SQLiteExoplanetRepository) GetExoplanetByID(id string) (models.Exoplanet, error) {
	var exoplanet models.Exoplanet
	err := r.DB.QueryRow("SELECT id, name, description, distance, radius, mass, type FROM exoplanets WHERE id = ?", id).Scan(
		&exoplanet.ID, &exoplanet.Name, &exoplanet.Description, &exoplanet.Distance, &exoplanet.Radius, &exoplanet.Mass, &exoplanet.Type)
	if err == sql.ErrNoRows {
		return exoplanet, errors.New("exoplanet not found")
	}
	if err != nil {
		return exoplanet, errors.New("failed to get exoplanet by ID")
	}
	return exoplanet, nil
}

func (r *SQLiteExoplanetRepository) UpdateExoplanet(exoplanet models.Exoplanet) error {
	_, err := r.DB.Exec("UPDATE exoplanets SET name = ?, description = ?, distance = ?, radius = ?, mass = ?, type = ? WHERE id = ?",
		exoplanet.Name, exoplanet.Description, exoplanet.Distance, exoplanet.Radius, exoplanet.Mass, exoplanet.Type, exoplanet.ID)
	if err != nil {
		return errors.New("failed to update exoplanet")
	}
	return nil
}

func (r *SQLiteExoplanetRepository) DeleteExoplanet(id string) error {
	_, err := r.DB.Exec("DELETE FROM exoplanets WHERE id = ?", id)
	if err != nil {
		return errors.New("failed to delete exoplanet")
	}
	return nil
}
