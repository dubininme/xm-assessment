//go:build unit

package company

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCompanyName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid_single_char",
			input:   "A",
			wantErr: nil,
		},
		{
			name:    "valid_max_length",
			input:   "123456789012345", // exactly 15 characters
			wantErr: nil,
		},
		{
			name:    "invalid_empty",
			input:   "",
			wantErr: ErrInvalidCompanyNameLength,
		},
		{
			name:    "invalid_too_long",
			input:   "1234567890123456", // 16 characters
			wantErr: ErrInvalidCompanyNameLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := NewCompanyName(tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, tt.input, res.String())
			}

		})
	}
}

func TestNewCompanyDescription(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid_empty",
			input:   "",
			wantErr: nil,
		},
		{
			name:    "valid_short",
			input:   "A great company",
			wantErr: nil,
		},
		{
			name:    "valid_max_length",
			input:   strings.Repeat("a", 3000), // exactly 3000 characters
			wantErr: nil,
		},
		{
			name:    "invalid_too_long",
			input:   strings.Repeat("a", 3001), // 3001 characters
			wantErr: ErrInvalidCompanyDescriptionLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewCompanyDescription(tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.input, result.String())
			}
		})
	}
}

func TestNewEmployeesCount(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr error
	}{
		{
			name:    "valid_min",
			input:   1,
			wantErr: nil,
		},
		{
			name:    "valid_small",
			input:   10,
			wantErr: nil,
		},
		{
			name:    "valid_large",
			input:   1000000,
			wantErr: nil,
		},
		{
			name:    "invalid_zero",
			input:   0,
			wantErr: ErrInvalidEmployeesCount,
		},
		{
			name:    "invalid_negative",
			input:   -1,
			wantErr: ErrInvalidEmployeesCount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewEmployeesCount(tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.input, result.Int())
			}
		})
	}
}

func TestNewCompanyType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid_corporations",
			input:   "Corporations",
			wantErr: nil,
		},
		{
			name:    "invalid_random_string",
			input:   "RandomType",
			wantErr: ErrInvalidCompanyType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewCompanyType(tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.input, result.String())
			}
		})
	}
}

func TestCompanyTypeFromInt(t *testing.T) {
	tests := []struct {
		name     string
		input    int16
		expected CompanyType
		wantErr  error
	}{
		{
			name:     "int_to_corporations",
			input:    1,
			expected: CorporationsType,
			wantErr:  nil,
		},
		{
			name:     "invalid_unknown",
			input:    999,
			expected: "",
			wantErr:  ErrInvalidCompanyType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompanyTypeFromInt(tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCompanyType_Roundtrip(t *testing.T) {
	// Critical: ensure String -> Int -> String works
	cType := CorporationsType
	intVal := cType.Int()
	fromInt, err := CompanyTypeFromInt(intVal)

	require.NoError(t, err)
	assert.Equal(t, cType, fromInt)
}
