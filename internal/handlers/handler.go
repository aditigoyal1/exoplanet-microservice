package handlers

import (
	"encoding/json"
	"exoplanet-microservice/internal/models"
	"exoplanet-microservice/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ExoplanetHandler struct {
	Service *service.ExoplanetService
}

func NewExoplanetHandler(service *service.ExoplanetService) *ExoplanetHandler {
	return &ExoplanetHandler{Service: service}
}

func (h *ExoplanetHandler) AddExoplanet(w http.ResponseWriter, r *http.Request) {
	var exoplanet models.Exoplanet
	err := json.NewDecoder(r.Body).Decode(&exoplanet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exoplanet, err = h.Service.AddExoplanet(exoplanet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exoplanet)
}

func (h *ExoplanetHandler) ListExoplanets(w http.ResponseWriter, r *http.Request) {
	exoplanets, err := h.Service.ListExoplanets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(exoplanets)
}

func (h *ExoplanetHandler) GetExoplanetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	exoplanet, err := h.Service.GetExoplanetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(exoplanet)
}

func (h *ExoplanetHandler) UpdateExoplanet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var exoplanet models.Exoplanet
	err := json.NewDecoder(r.Body).Decode(&exoplanet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exoplanet.ID = id
	err = h.Service.UpdateExoplanet(exoplanet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(exoplanet)
}

func (h *ExoplanetHandler) DeleteExoplanet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	err := h.Service.DeleteExoplanet(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ExoplanetHandler) FuelEstimation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	crewCapacity, err := strconv.Atoi(r.URL.Query().Get("crew_capacity"))
	if err != nil || crewCapacity <= 0 {
		http.Error(w, "Invalid crew capacity", http.StatusBadRequest)
		return
	}

	fuel, err := h.Service.FuelEstimation(id, crewCapacity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"fuel": fuel})
}
