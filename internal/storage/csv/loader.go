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

// Model represents any data model that can be stored
type Model interface {
	GetID() uint
}

// Ensure T is a struct type that implements Model
type structModel interface {
	Model
	comparable
}

// Store is a generic store for struct types implementing Model
type Store[T structModel] struct {
	items     map[uint]T
	validator *HeaderValidator
}

func NewStore[T structModel]() *Store[T] {
	// Compile-time validation that T is a struct type
	var zero T
	if reflect.TypeOf(zero).Kind() != reflect.Struct {
		panic("NewStore: type parameter T must be a struct type")
	}

	return &Store[T]{
		items:     make(map[uint]T),
		validator: &HeaderValidator{RequireExact: true},
	}
}

func (s *Store[T]) LoadFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("opening CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read first row of data to get headers
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("reading CSV headers: %w", err)
	}

	// Get type information once at the start
	modelType := reflect.TypeOf(*new(T))

	if err := s.validator.ValidateHeaders(headers, modelType); err != nil {
		return fmt.Errorf("validating CSV headers: %w", err)
	}

	positions, err := s.validator.GetHeaderPositions(headers, modelType)
	if err != nil {
		return fmt.Errorf("getting header positions: %w", err)
	}

	parser := &CSVParser{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading CSV record: %w", err)
		}

		item, err := parser.ParseRecord(record, positions, modelType)
		if err != nil {
			return fmt.Errorf("parsing record: %w", err)
		}

		// We know this must succeed because of our type constraints
		typed := item.(T)
		s.items[typed.GetID()] = typed
	}

	return nil
}

func (s *Store[T]) GetByID(id uint) (T, error) {
	if item, ok := s.items[id]; ok {
		return item, nil
	}
	return *new(T), fmt.Errorf("item with ID %d not found", id)
}

func (s *Store[T]) List() []T {
	result := make([]T, len(s.items))
	i := 0
	for _, item := range s.items {
		result[i] = item
		i++
	}
	return result
}

// Search returns instances matching the given criteria
func (s *Store[T]) Search(criteria func(T) bool) []T {
	var results []T
	for _, item := range s.items {
		if criteria(item) {
			results = append(results, item)
		}
	}
	return results
}

// Helper function to parse unsigned integers with meaningful error messages
func parseUint(s string, field string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %s: %w", field, err)
	}
	return uint(val), nil
}
