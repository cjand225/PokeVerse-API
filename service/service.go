// Package service implements the business logic for the web service.
package service

import (
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"pokeverse/web-service/database"
)

// Define the BaseStats struct separately
type BaseStats struct {
	HP             int `json:"HP"`
	Speed          int `json:"Speed"`
	Attack         int `json:"Attack"`
	Defense        int `json:"Defense"`
	SpecialAttack  int `json:"Special Attack"`
	SpecialDefense int `json:"Special Defense"`
}

// Define the Pokemon struct with the embedded BaseStats field
type Pokemon struct {
	ID         int       `json:"ID"`
	Name       string    `json:"Name"`
	Type       []string  `json:"Type"`
	BaseStats  BaseStats `json:"Base Stats"`
	Generation int       `json:"Generation"`
}

// Service handles the business logic for the web service.
type Service struct {
	pool *pgxpool.Pool
}

// NewService creates a new instance of Service with the provided connection pool.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// GetPokemonByID retrieves a Pokemon by its ID and language from the database.
// It constructs the SQL query using the provided ID and language parameters,
// executes the query, and unmarshals the JSON data into a Pokemon struct.
//
// Parameters:
//   id: The ID of the Pokemon to retrieve.
//   lang: The language code for the Pokemon's data.
//
// Returns:
//   *Pokemon: A pointer to the retrieved Pokemon data.
//   error: An error if any issue occurs during the database query or JSON unmarshaling.
func (s *Service) GetPokemonByID(id int, lang string) (*Pokemon, error) {
	byteData, err := database.Query(s.pool, "SELECT pokedex.getpokemon($1, $2);", id, lang)
	if err != nil {

		log.Print(err.Error())
		return nil, err
	}

	pokemon := &Pokemon{}
	err = json.Unmarshal(byteData, pokemon)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return pokemon, nil
}
