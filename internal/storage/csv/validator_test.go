// csv/validator_test.go
package csv

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name      string `csv:"name"`
	Optional  string `csv:"optional,omitempty"`
	NoCSV     string
	Different string `csv:"csv_name"`
	ID        uint   `csv:"id"`
}

func Test_HeaderValidator_ValidateHeaders(testContext *testing.T) {
	validator := &HeaderValidator{RequireExact: true}
	structType := reflect.TypeOf(TestStruct{})

	testCases := []struct {
		name      string
		errorType string
		headers   []string
		wantError bool
	}{
		{
			name:      "valid headers exact match",
			headers:   []string{"id", "name", "optional", "csv_name"},
			wantError: false,
		},
		{
			name:      "missing required header",
			headers:   []string{"id", "name", "csv_name"},
			wantError: true,
			errorType: "missing",
		},
		{
			name:      "extra header",
			headers:   []string{"id", "name", "optional", "csv_name", "extra"},
			wantError: true,
			errorType: "extra",
		},
		{
			name:      "wrong order still valid",
			headers:   []string{"name", "id", "csv_name", "optional"},
			wantError: false,
		},
	}

	for _, testCase := range testCases {
		testContext.Run(testCase.name, func(testContext *testing.T) {
			err := validator.ValidateHeaders(testCase.headers, structType)

			if testCase.wantError {
				assert.Error(testContext, err)
				validationErr, ok := err.(*ValidationError)
				assert.True(testContext, ok, "error should be ValidationError")

				switch testCase.errorType {
				case "missing":
					assert.NotEmpty(testContext, validationErr.Missing)
				case "extra":
					assert.NotEmpty(testContext, validationErr.Extra)
				}
			} else {
				assert.NoError(testContext, err)
			}
		})
	}
}

func Test_HeaderValidator_NonExactMode(testContext *testing.T) {
	validator := &HeaderValidator{RequireExact: false}
	structType := reflect.TypeOf(TestStruct{})

	// Test that extra headers don't cause errors when RequireExact is false
	headers := []string{"id", "name", "optional", "csv_name", "extra_ok"}
	err := validator.ValidateHeaders(headers, structType)
	assert.NoError(testContext, err)
}

func Test_HeaderValidator_GetHeaderPositions(testContext *testing.T) {
	validator := &HeaderValidator{RequireExact: true}
	structType := reflect.TypeOf(TestStruct{})

	headers := []string{"name", "id", "csv_name", "optional"}
	positions, err := validator.GetHeaderPositions(headers, structType)

	assert.NoError(testContext, err)
	assert.Equal(testContext, 1, positions["id"])
	assert.Equal(testContext, 0, positions["name"])
	assert.Equal(testContext, 2, positions["csv_name"])
	assert.Equal(testContext, 3, positions["optional"])
}
