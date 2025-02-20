package grpc

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBlockNumArg(t *testing.T) {
	t.Parallel()

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
			expected: "0x3039",
		},
		{
			name:     "Zero",
			input:    big.NewInt(0),
			expected: "0x0",
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ToBlockNumArg(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
