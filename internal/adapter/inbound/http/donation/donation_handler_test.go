package donation

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestValidateAndParseDonationAmount(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    decimal.Decimal
		expectError bool
	}{
		{
			name:        "valid decimal",
			input:       "123.45",
			expected:    decimal.NewFromFloat(123.45),
			expectError: false,
		},
		{
			name:        "valid integer",
			input:       "1000",
			expected:    decimal.NewFromFloat(1000),
			expectError: false,
		},
		{
			name:        "valid with currency symbol",
			input:       "$1,000.50",
			expected:    decimal.NewFromFloat(1000.50),
			expectError: false,
		},
		{
			name:        "valid with euro symbol",
			input:       "€1.000,50",
			expected:    decimal.NewFromFloat(1000.50),
			expectError: false,
		},
		{
			name:        "valid with spaces",
			input:       "1 000.50",
			expected:    decimal.NewFromFloat(1000.50),
			expectError: false,
		},
		{
			name:        "valid with multiple commas",
			input:       "1,000,000.50",
			expected:    decimal.NewFromFloat(1000000.50),
			expectError: false,
		},
		{
			name:        "negative amount (parsing only)",
			input:       "-100.50",
			expected:    decimal.NewFromFloat(-100.50),
			expectError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    decimal.Zero,
			expectError: true,
		},
		{
			name:        "invalid format",
			input:       "abc",
			expected:    decimal.NewFromFloat(0),
			expectError: false,
		},
		{
			name:        "zero amount (parsing only)",
			input:       "0",
			expected:    decimal.NewFromFloat(0),
			expectError: false,
		},
		{
			name:        "negative amount (parsing only)",
			input:       "-100",
			expected:    decimal.NewFromFloat(-100),
			expectError: false,
		},
		{
			name:        "amount too large (parsing only)",
			input:       "100000001",
			expected:    decimal.NewFromFloat(100000001),
			expectError: false,
		},
		{
			name:        "amount with USD suffix",
			input:       "1,000 USD",
			expected:    decimal.NewFromFloat(1000),
			expectError: false,
		},
		{
			name:        "amount with trailing decimal",
			input:       "1000.",
			expected:    decimal.NewFromFloat(1000),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDonationAmount(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, result.Equal(tt.expected), "Expected %s, got %s", tt.expected.String(), result.String())
			}
		})
	}
}

func TestCleanAmountString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"currency symbol", "$1,000.50", "1000.50"},
		{"euro format", "€1.000,50", "1000.50"},
		{"with spaces", "1 000.50", "1000.50"},
		{"multiple commas", "1,000,000.50", "1000000.50"},
		{"negative", "-100.50", "-100.50"},
		{"trailing decimal", "1000.", "1000"},
		{"with USD suffix", "1,000 USD", "1000"},
		{"mixed format", "$1,000.50 USD", "1000.50"},
		{"empty", "", "0"},
		{"only minus", "-", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanAmountString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
