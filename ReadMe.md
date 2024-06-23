Enpoints Curl

curl -X POST http://localhost:8080/exoplanets -H "Content-Type: application/json" -d '{
  "name": "PlanetX",
  "description": "A mysterious gas giant",
  "distance": 200,
  "radius": 2.5,
  "type": "GasGiant"
}'




curl -X POST http://localhost:8080/exoplanets -H "Content-Type: application/json" -d '{
  "name": "PlanetY",
  "description": "A rocky terrestrial planet",
  "distance": 150,
  "radius": 1.2,
  "mass": 5.5,
  "type": "Terrestrial"
}'


List All Exoplanets (GET /exoplanets)

curl -X GET http://localhost:8080/exoplanets


curl -X GET http://localhost:8080/exoplanets/<id>


curl -X PUT http://localhost:8080/exoplanets/1234 -H "Content-Type: application/json" -d '{
  "name": "UpdatedPlanet",
  "description": "An updated description",
  "distance": 250,
  "radius": 3.0,
  "mass": 6.0,
  "type": "Terrestrial"
}'


curl -X DELETE http://localhost:8080/exoplanets/1234

Fuel Estimation (GET /exoplanets/{id}/fuel-estimation)
curl -X GET "http://localhost:8080/exoplanets/1234/fuel-estimation?crew_capacity=10"


# Exoplanet Microservice

This is a simple Go microservice for managing exoplanets. It supports adding, listing, updating, deleting exoplanets, and estimating fuel requirements for space voyages.

## Running the Service

### Prerequisites

- Docker
- Go (if running locally)

### Running with Docker

```sh
docker build -t exoplanet-microservice .
docker run -p 8080:8080 exoplanet-microservice



