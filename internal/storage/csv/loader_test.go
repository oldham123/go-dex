// internal/storage/csv/loader_test.go
package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseUint(testContext *testing.T) {
	testCases := []struct {
		name      string
		input     string
		field     string
		want      uint
		wantError bool
	}{
		{
			name:      "valid number",
			input:     "123",
			field:     "test_field",
			want:      123,
			wantError: false,
		},
		{
			name:      "zero",
			input:     "0",
			field:     "test_field",
			want:      0,
			wantError: false,
		},
		{
			name:      "invalid number",
			input:     "abc",
			field:     "test_field",
			want:      0,
			wantError: true,
		},
		{
			name:      "negative number",
			input:     "-123",
			field:     "test_field",
			want:      0,
			wantError: true,
		},
		{
			name:      "empty string",
			input:     "",
			field:     "test_field",
			want:      0,
			wantError: true,
		},
	}

	for _, testCase := range testCases {
		testContext.Run(testCase.name, func(testContext *testing.T) {
			got, err := parseUint(testCase.input, testCase.field)

			if testCase.wantError {
				assert.Error(testContext, err)
				assert.Contains(testContext, err.Error(), testCase.field)
			} else {
				assert.NoError(testContext, err)
				assert.Equal(testContext, testCase.want, got)
			}
		})
	}
}
