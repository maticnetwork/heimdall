package grpc

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBlockNumArg(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected string
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: "latest",
		},
		{
			name:     "Positive number",
			input:    big.NewInt(12345),
			expected: "12345",
		},
		{
			name:     "Zero",
			input:    big.NewInt(0),
			expected: "0",
		},
		{
			name:     "Negative number",
			input:    big.NewInt(-1),
			expected: "pending",
		},
		{
			name:     "Large negative number",
			input:    big.NewInt(-1234567890),
			expected: "<invalid -1234567890>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToBlockNumArg(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
