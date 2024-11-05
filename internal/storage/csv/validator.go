package csv

import (
	"fmt"
	"reflect"
	"strings"
)

// ValidationError wraps errors that occur during header validation
type ValidationError struct {
	Expected []string
	Got      []string
	Missing  []string
	Extra    []string
}

func (e *ValidationError) Error() string {
	var msg strings.Builder
	msg.WriteString("CSV header validation failed:\n")

	if len(e.Missing) > 0 {
		msg.WriteString(fmt.Sprintf("Missing required headers: %v\n", e.Missing))
	}
	if len(e.Extra) > 0 {
		msg.WriteString(fmt.Sprintf("Unexpected headers found: %v\n", e.Extra))
	}
	msg.WriteString(fmt.Sprintf("Expected: %v\n", e.Expected))
	msg.WriteString(fmt.Sprintf("Got: %v", e.Got))

	return msg.String()
}

// HeaderValidator handles CSV header validation using struct reflection
type HeaderValidator struct {
	// RequireExact determines if extra columns should cause validation failure
	RequireExact bool
}

// ValidateHeaders checks if the provided headers match the struct's csv tags
func (v *HeaderValidator) ValidateHeaders(headers []string, structType reflect.Type) error {
	expected := v.getExpectedHeaders(structType)

	// Create maps for easy lookup
	expectedMap := make(map[string]bool)
	for _, h := range expected {
		expectedMap[h] = true
	}

	gotMap := make(map[string]bool)
	for _, h := range headers {
		gotMap[h] = true
	}

	// Find missing and extra headers
	var missing, extra []string

	for _, exp := range expected {
		if !gotMap[exp] {
			missing = append(missing, exp)
		}
	}

	if v.RequireExact {
		for _, got := range headers {
			if !expectedMap[got] {
				extra = append(extra, got)
			}
		}
	}

	// If we found any problems, return validation error
	if len(missing) > 0 || (v.RequireExact && len(extra) > 0) {
		return &ValidationError{
			Expected: expected,
			Got:      headers,
			Missing:  missing,
			Extra:    extra,
		}
	}

	return nil
}

// getExpectedHeaders extracts CSV header names from struct tags
func (v *HeaderValidator) getExpectedHeaders(structType reflect.Type) []string {
	var headers []string
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Get the csv tag, skip if empty
		if csvTag := field.Tag.Get("csv"); csvTag != "" {
			// Split the tag on comma in case there are options
			tagParts := strings.Split(csvTag, ",")
			headers = append(headers, tagParts[0])
		}
	}

	return headers
}

// GetHeaderPositions returns a map of header names to their positions
func (v *HeaderValidator) GetHeaderPositions(headers []string, structType reflect.Type) (map[string]int, error) {
	// Validate headers first
	if err := v.ValidateHeaders(headers, structType); err != nil {
		return nil, err
	}

	positions := make(map[string]int)
	for i, header := range headers {
		positions[header] = i
	}

	return positions, nil
}
