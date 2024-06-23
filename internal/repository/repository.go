package repository

import (
	"exoplanet-microservice/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type ExoplanetRepository interface {
	AddExoplanet(exoplanet models.Exoplanet) error
	ListExoplanets() ([]models.Exoplanet, error)
	GetExoplanetByID(id string) (models.Exoplanet, error)
	UpdateExoplanet(exoplanet models.Exoplanet) error
	DeleteExoplanet(id string) error
}
