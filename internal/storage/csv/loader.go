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

// Store is a generic store for any model type
type Store[T Model] struct {
	items     map[uint]T
	validator *HeaderValidator
}

func NewStore[T Model]() *Store[T] {
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

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("reading CSV headers: %w", err)
	}

	// Create a zero value of T to get its type
	var zero T
	modelType := reflect.TypeOf(zero)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// Validate headers using our validator
	if err := s.validator.ValidateHeaders(headers, modelType); err != nil {
		return fmt.Errorf("validating CSV headers: %w", err)
	}

	// Get header positions for efficient record parsing
	positions, err := s.validator.GetHeaderPositions(headers, modelType)
	if err != nil {
		return fmt.Errorf("getting header positions: %w", err)
	}

	// Create parser
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

		// Parse the record into our model type
		result, err := parser.ParseRecord(record, positions, modelType)
		if err != nil {
			return fmt.Errorf("parsing record: %w", err)
		}

		// Convert the result to our generic type T
		var item T
		switch v := result.(type) {
		case T:
			item = v
		default:
			// If result is a struct but we need a pointer
			if reflect.TypeOf(zero).Kind() == reflect.Ptr {
				// Create a new pointer to the result
				ptr := reflect.New(reflect.TypeOf(result))
				ptr.Elem().Set(reflect.ValueOf(result))
				item = ptr.Interface().(T)
			} else {
				return fmt.Errorf("unexpected type: got %T, want %T", result, zero)
			}
		}

		// Store the item using its ID
		s.items[item.GetID()] = item
	}

	return nil
}

func (s *Store[T]) GetByID(id uint) (T, error) {
	if item, ok := s.items[id]; ok {
		return item, nil
	}
	var zero T
	return zero, fmt.Errorf("item with ID %d not found", id)
}

func (s *Store[T]) List() []T {
	result := make([]T, 0, len(s.items))
	for _, item := range s.items {
		result = append(result, item)
	}
	return result
}

// Search returns Pokemon matching the given criteria
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
