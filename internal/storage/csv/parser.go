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
	// We know targetType is a struct type now, so we can remove the pointer check
	result := reflect.New(targetType).Elem()

	// Pre-calculate field mappings once
	fields := make([]struct {
		fieldTyp reflect.Type
		name     string
		index    int
		csvPos   int
	}, 0, targetType.NumField())

	// Build field mapping
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		if csvTag := field.Tag.Get("csv"); csvTag != "" {
			if pos, exists := positions[csvTag]; exists {
				fields = append(fields, struct {
					fieldTyp reflect.Type
					name     string
					index    int
					csvPos   int
				}{field.Type, csvTag, i, pos})
			}
		}
	}

	// Validate record length once
	maxPos := 0
	for _, f := range fields {
		if f.csvPos > maxPos {
			maxPos = f.csvPos
		}
	}
	if len(record) <= maxPos {
		return nil, fmt.Errorf("record too short: got %d fields, need at least %d", len(record), maxPos+1)
	}

	// Parse all fields
	var parseErr *ParseError
	for _, f := range fields {
		value := record[f.csvPos]
		parsedValue, err := p.parseValue(value, f.fieldTyp, f.name)
		if err != nil {
			parseErr = &ParseError{
				Field: f.name,
				Value: value,
				Type:  f.fieldTyp.String(),
				Err:   err,
			}
			break
		}
		result.Field(f.index).Set(reflect.ValueOf(parsedValue))
	}

	if parseErr != nil {
		return nil, parseErr
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
