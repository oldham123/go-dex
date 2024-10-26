// internal/storage/csv/loader.go
package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Pokemon represents a single row from our CSV data
type Pokemon struct {
	Identifier     string `json:"identifier"`
	Id             uint   `json:"id"`
	SpeciesId      uint   `json:"species_id"`
	Height         uint   `json:"height"`
	Weight         uint   `json:"weight"`
	BaseExperience uint   `json:"base_experience"`
	Order          uint   `json:"order"`
	IsDefault      bool   `json:"is_default"`
}

// Store handles the loaded Pokemon data
type Store struct {
	pokemon map[uint]*Pokemon
}

func NewStore() *Store {
	return &Store{
		pokemon: make(map[uint]*Pokemon),
	}
}

// LoadFromCSV loads Pokemon data from a CSV file
func (s *Store) LoadFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("opening CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read and validate headers
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("reading CSV headers: %w", err)
	}
	if err := validateHeaders(headers); err != nil {
		return fmt.Errorf("validating CSV headers: %w", err)
	}

	// Read data rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading CSV record: %w", err)
		}

		pokemon, err := parsePokemonRecord(record)
		if err != nil {
			return fmt.Errorf("parsing Pokemon record: %w", err)
		}

		s.pokemon[pokemon.Id] = pokemon
	}

	return nil
}

// validateHeaders checks if the CSV has the expected column structure
func validateHeaders(headers []string) error {
	expected := []string{
		"id",
		"identifier",
		"species_id",
		"height",
		"weight",
		"base_experience",
		"order",
		"is_default",
	}

	if len(headers) != len(expected) {
		return fmt.Errorf("expected %d columns, got %d", len(expected), len(headers))
	}

	for i, header := range headers {
		if header != expected[i] {
			return fmt.Errorf("expected column %d to be %s, got %s", i, expected[i], header)
		}
	}
	return nil
}

// parsePokemonRecord converts a CSV record into a Pokemon struct
func parsePokemonRecord(record []string) (*Pokemon, error) {
	if len(record) != 8 {
		return nil, fmt.Errorf("invalid record length: expected 8, got %d", len(record))
	}

	// Parse unsigned integer values
	id, err := parseUint(record[0], "ID")
	if err != nil {
		return nil, err
	}

	speciesID, err := parseUint(record[2], "species_id")
	if err != nil {
		return nil, err
	}

	height, err := parseUint(record[3], "height")
	if err != nil {
		return nil, err
	}

	weight, err := parseUint(record[4], "weight")
	if err != nil {
		return nil, err
	}

	baseExp, err := parseUint(record[5], "base_experience")
	if err != nil {
		return nil, err
	}

	order, err := parseUint(record[6], "order")
	if err != nil {
		return nil, err
	}

	isDefault, err := strconv.ParseBool(record[7])
	if err != nil {
		return nil, fmt.Errorf("parsing is_default: %w", err)
	}

	return &Pokemon{
		Id:             id,
		Identifier:     record[1],
		SpeciesId:      speciesID,
		Height:         height,
		Weight:         weight,
		BaseExperience: baseExp,
		Order:          order,
		IsDefault:      isDefault,
	}, nil
}

// Helper function to parse unsigned integers with meaningful error messages
func parseUint(s string, field string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %s: %w", field, err)
	}
	return uint(val), nil
}

// Access methods for the API

// GetByID returns a Pokemon by its ID
func (s *Store) GetByID(id uint) (*Pokemon, error) {
	pokemon, exists := s.pokemon[id]
	if !exists {
		return nil, fmt.Errorf("pokemon with ID %d not found", id)
	}
	return pokemon, nil
}

// List returns all Pokemon
func (s *Store) List() []*Pokemon {
	pokemon := make([]*Pokemon, 0, len(s.pokemon))
	for _, p := range s.pokemon {
		pokemon = append(pokemon, p)
	}
	return pokemon
}

// Search returns Pokemon matching the given criteria
func (s *Store) Search(criteria func(*Pokemon) bool) []*Pokemon {
	var results []*Pokemon
	for _, p := range s.pokemon {
		if criteria(p) {
			results = append(results, p)
		}
	}
	return results
}
