package service

import (
	"errors"
	"exoplanet-microservice/internal/models"
	"exoplanet-microservice/internal/repository"

	"github.com/google/uuid"
)

type ExoplanetService struct {
	Repo repository.ExoplanetRepository
}

func NewExoplanetService(repo repository.ExoplanetRepository) *ExoplanetService {
	return &ExoplanetService{Repo: repo}
}

func (s *ExoplanetService) AddExoplanet(exoplanet models.Exoplanet) (models.Exoplanet, error) {
	if err := validateExoplanet(exoplanet); err != nil {
		return exoplanet, err
	}
	exoplanet.ID = uuid.New().String()
	return exoplanet, s.Repo.AddExoplanet(exoplanet)
}

func (s *ExoplanetService) ListExoplanets() ([]models.Exoplanet, error) {
	return s.Repo.ListExoplanets()
}

func (s *ExoplanetService) GetExoplanetByID(id string) (models.Exoplanet, error) {
	return s.Repo.GetExoplanetByID(id)
}

func (s *ExoplanetService) UpdateExoplanet(exoplanet models.Exoplanet) error {
	if err := validateExoplanet(exoplanet); err != nil {
		return err
	}
	return s.Repo.UpdateExoplanet(exoplanet)
}

func (s *ExoplanetService) DeleteExoplanet(id string) error {
	return s.Repo.DeleteExoplanet(id)
}

func (s *ExoplanetService) FuelEstimation(id string, crewCapacity int) (float64, error) {
	exoplanet, err := s.Repo.GetExoplanetByID(id)
	if err != nil {
		return 0, err
	}

	gravity := calculateGravity(exoplanet)
	fuel := calculateFuel(exoplanet.Distance, gravity, crewCapacity)

	return fuel, nil
}

func validateExoplanet(exoplanet models.Exoplanet) error {
	if exoplanet.Distance < 10 || exoplanet.Distance > 1000 {
		return errors.New("distance must be between 10 and 1000 light years")
	}
	if exoplanet.Radius < 0.1 || exoplanet.Radius > 10 {
		return errors.New("radius must be between 0.1 and 10 Earth-radius units")
	}
	if exoplanet.Type == models.Terrestrial {
		if exoplanet.Mass < 0.1 || exoplanet.Mass > 10 {
			return errors.New("mass must be between 0.1 and 10 Earth-mass units")
		}
	} else if exoplanet.Type != models.GasGiant {
		return errors.New("type must be either GasGiant or Terrestrial")
	}
	return nil
}

func calculateGravity(exoplanet models.Exoplanet) float64 {
	if exoplanet.Type == models.GasGiant {
		return 0.5 / (exoplanet.Radius * exoplanet.Radius)
	} else if exoplanet.Type == models.Terrestrial {
		return exoplanet.Mass / (exoplanet.Radius * exoplanet.Radius)
	}
	return 0
}

func calculateFuel(distance int, gravity float64, crewCapacity int) float64 {
	return float64(distance) / (gravity * gravity) * float64(crewCapacity)
}
