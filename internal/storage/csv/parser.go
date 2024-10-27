package csv

import (
	"fmt"
	"reflect"
	"strconv"
)

// CSVParser handles parsing of CSV records into structs
type CSVParser struct {
	// Add configuration options here, like:
	// - strict mode for type conversion
	// - custom type handlers
	// - null/empty value handling
}

// ParseError provides detailed information about parsing failures
type ParseError struct {
	Err   error
	Field string
	Value string
	Type  string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("failed to parse field %q (value: %q) as %s: %v",
		e.Field, e.Value, e.Type, e.Err)
}

// ParseRecord parses a CSV record into a new instance of the given type
func (p *CSVParser) ParseRecord(record []string, positions map[string]int, targetType reflect.Type) (interface{}, error) {
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	// Create a new instance of the target type
	result := reflect.New(targetType).Elem()

	// Iterate through the struct fields
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)

		// Get the csv tag
		csvTag := field.Tag.Get("csv")
		if csvTag == "" {
			continue // Skip fields without csv tag
		}

		// Find position in record
		pos, exists := positions[csvTag]
		if !exists {
			continue // Skip if column wasn't in the CSV
		}

		// Get the value from the record
		if pos >= len(record) {
			return nil, fmt.Errorf("record too short, missing field %q", csvTag)
		}
		value := record[pos]

		// Parse the value according to the field type
		parsedValue, err := p.parseValue(value, field.Type, csvTag)
		if err != nil {
			return nil, &ParseError{
				Field: csvTag,
				Value: value,
				Type:  field.Type.String(),
				Err:   err,
			}
		}

		// Set the field value
		result.Field(i).Set(reflect.ValueOf(parsedValue))
	}

	return result.Interface(), nil
}

// parseValue converts a string value to the appropriate type
func (p *CSVParser) parseValue(value string, targetType reflect.Type, fieldName string) (interface{}, error) {
	// Handle empty values by returning the zero value for the type
	if value == "" {
		return reflect.Zero(targetType).Interface(), nil
	}

	switch targetType.Kind() {
	case reflect.String:
		return value, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing %s as unsigned integer: %w", fieldName, err)
		}
		// Convert to the specific uint type
		return reflect.ValueOf(val).Convert(targetType).Interface(), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing %s as integer: %w", fieldName, err)
		}
		return reflect.ValueOf(val).Convert(targetType).Interface(), nil

	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing %s as float: %w", fieldName, err)
		}
		return reflect.ValueOf(val).Convert(targetType).Interface(), nil

	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("parsing %s as float: %w", fieldName, err)
		}
		return val, nil

	default:
		return nil, fmt.Errorf("unsupported type: %v", targetType)
	}
}
