package csv

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CSVParser_ParseRecord(testContext *testing.T) {
	type TestStruct struct {
		StringField string `csv:"string"`
		SkipField   string
		IntField    int     `csv:"int"`
		UintField   uint    `csv:"uint"`
		FloatField  float64 `csv:"float"`
		BoolField   bool    `csv:"bool"`
	}

	parser := &CSVParser{}

	positions := map[string]int{
		"string": 0,
		"int":    1,
		"uint":   2,
		"float":  3,
		"bool":   4,
	}

	testCases := []struct {
		name      string
		record    []string
		wantError bool
	}{
		{
			name:      "valid record",
			record:    []string{"hello", "-42", "123", "3.14", "true"},
			wantError: false,
		},
		{
			name:      "invalid int",
			record:    []string{"hello", "not-a-number", "123", "3.14", "true"},
			wantError: true,
		},
		{
			name:      "invalid bool",
			record:    []string{"hello", "-42", "123", "3.14", "not-a-bool"},
			wantError: true,
		},
	}

	for _, testCase := range testCases {
		testContext.Run(testCase.name, func(testContext *testing.T) {
			result, err := parser.ParseRecord(
				testCase.record,
				positions,
				reflect.TypeOf(TestStruct{}),
			)

			if testCase.wantError {
				assert.Error(testContext, err)
				return
			}

			assert.NoError(testContext, err)
			parsed := result.(TestStruct)

			// Verify some values
			if !testCase.wantError {
				assert.Equal(testContext, "hello", parsed.StringField)
				assert.Equal(testContext, -42, parsed.IntField)
				assert.Equal(testContext, uint(123), parsed.UintField)
				assert.Equal(testContext, 3.14, parsed.FloatField)
				assert.True(testContext, parsed.BoolField)
			}
		})
	}
}
