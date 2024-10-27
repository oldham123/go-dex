// internal/storage/csv/loader.go
package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
)

// Pokemon represents a single row from our CSV data
type Pokemon struct {
	Identifier     string `csv:"identifier"`
	Id             uint   `csv:"id"`
	SpeciesId      uint   `csv:"species_id"`
	Height         uint   `csv:"height"`
	Weight         uint   `csv:"weight"`
	BaseExperience uint   `csv:"base_experience"`
	Order          uint   `csv:"order"`
	IsDefault      bool   `csv:"is_default"`
}

// Store handles the loaded Pokemon data
type Store struct {
	pokemon   map[uint]*Pokemon
	validator *HeaderValidator
}

func NewStore() *Store {
	return &Store{
		pokemon:   make(map[uint]*Pokemon),
		validator: &HeaderValidator{RequireExact: true},
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

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("reading CSV headers: %w", err)
	}

	// Validate headers using our new validator
	if err := s.validator.ValidateHeaders(headers, reflect.TypeOf(Pokemon{})); err != nil {
		return fmt.Errorf("validating CSV headers: %w", err)
	}

	// Get header positions for efficient record parsing
	positions, err := s.validator.GetHeaderPositions(headers, reflect.TypeOf(Pokemon{}))
	if err != nil {
		return fmt.Errorf("getting header positions: %w", err)
	}

	parser := &CSVParser{}

	// Read data rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading CSV record: %w", err)
		}

		// pokemon, err := parsePokemonRecord(record, positions)
		// if err != nil {
		// 	return fmt.Errorf("parsing Pokemon record: %w", err)
		// }

		result, err := parser.ParseRecord(record, positions, reflect.TypeOf(Pokemon{}))
		if err != nil {
			panic(err)
		}

		pokemon := result.(Pokemon)

		s.pokemon[pokemon.Id] = &pokemon
	}

	return nil
}

// parsePokemonRecord converts a CSV record into a Pokemon struct using header positions
func parsePokemonRecord(record []string, positions map[string]int) (*Pokemon, error) {
	if len(record) < len(positions) {
		return nil, fmt.Errorf("invalid record length: expected at least %d, got %d", len(positions), len(record))
	}

	// Parse unsigned integer values using positions map
	id, err := parseUint(record[positions["id"]], "id")
	if err != nil {
		return nil, err
	}

	speciesID, err := parseUint(record[positions["species_id"]], "species_id")
	if err != nil {
		return nil, err
	}

	height, err := parseUint(record[positions["height"]], "height")
	if err != nil {
		return nil, err
	}

	weight, err := parseUint(record[positions["weight"]], "weight")
	if err != nil {
		return nil, err
	}

	baseExp, err := parseUint(record[positions["base_experience"]], "base_experience")
	if err != nil {
		return nil, err
	}

	order, err := parseUint(record[positions["order"]], "order")
	if err != nil {
		return nil, err
	}

	isDefault, err := strconv.ParseBool(record[positions["is_default"]])
	if err != nil {
		return nil, fmt.Errorf("parsing is_default: %w", err)
	}

	return &Pokemon{
		Id:             id,
		Identifier:     record[positions["identifier"]],
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
