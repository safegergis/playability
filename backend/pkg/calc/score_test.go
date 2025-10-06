package calc

import (
	"testing"
)

func TestCalculateAccessibilityScore(t *testing.T) {
	tests := []struct {
		name     string
		scores   []int
		expected float64
	}{
		{
			name:     "empty array returns 0",
			scores:   []int{},
			expected: 0,
		},
		{
			name:     "single score",
			scores:   []int{7},
			expected: 7.0,
		},
		{
			name:     "multiple scores - average",
			scores:   []int{5, 7, 9},
			expected: 7.0,
		},
		{
			name:     "scores requiring rounding",
			scores:   []int{5, 6, 7},
			expected: 6.0,
		},
		{
			name:     "scores with decimal average - rounds to 1 decimal",
			scores:   []int{5, 5, 6},
			expected: 5.3,
		},
		{
			name:     "all zeros",
			scores:   []int{0, 0, 0},
			expected: 0.0,
		},
		{
			name:     "all max scores",
			scores:   []int{10, 10, 10},
			expected: 10.0,
		},
		{
			name:     "large dataset",
			scores:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: 5.5,
		},
		{
			name:     "precision test - 1 decimal place",
			scores:   []int{7, 8, 9},
			expected: 8.0,
		},
		{
			name:     "precision test - requires rounding up",
			scores:   []int{8, 9},
			expected: 8.5,
		},
		{
			name:     "precision test - requires rounding down",
			scores:   []int{1, 2, 3, 4},
			expected: 2.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateAccessibilityScore(tt.scores)
			if result != tt.expected {
				t.Errorf("CalculateAccessibilityScore(%v) = %v; expected %v", tt.scores, result, tt.expected)
			}
		})
	}
}

func TestToFixed(t *testing.T) {
	tests := []struct {
		name      string
		num       float64
		precision int
		expected  float64
	}{
		{
			name:      "1 decimal place",
			num:       5.333333,
			precision: 1,
			expected:  5.3,
		},
		{
			name:      "2 decimal places",
			num:       5.666666,
			precision: 2,
			expected:  5.67,
		},
		{
			name:      "round up at 0.5",
			num:       5.55,
			precision: 1,
			expected:  5.6,
		},
		{
			name:      "round down below 0.5",
			num:       5.44,
			precision: 1,
			expected:  5.4,
		},
		{
			name:      "zero precision",
			num:       5.6,
			precision: 0,
			expected:  6.0,
		},
		{
			name:      "negative number",
			num:       -3.666,
			precision: 1,
			expected:  -3.7,
		},
		{
			name:      "zero",
			num:       0.0,
			precision: 1,
			expected:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toFixed(tt.num, tt.precision)
			if result != tt.expected {
				t.Errorf("toFixed(%v, %d) = %v; expected %v", tt.num, tt.precision, result, tt.expected)
			}
		})
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		name     string
		num      float64
		expected int
	}{
		{
			name:     "round up at 0.5",
			num:      5.5,
			expected: 6,
		},
		{
			name:     "round up above 0.5",
			num:      5.6,
			expected: 6,
		},
		{
			name:     "round down below 0.5",
			num:      5.4,
			expected: 5,
		},
		{
			name:     "exact integer",
			num:      7.0,
			expected: 7,
		},
		{
			name:     "negative round up",
			num:      -5.5,
			expected: -6,
		},
		{
			name:     "negative round down",
			num:      -5.4,
			expected: -5,
		},
		{
			name:     "zero",
			num:      0.0,
			expected: 0,
		},
		{
			name:     "very small positive",
			num:      0.4,
			expected: 0,
		},
		{
			name:     "very small negative",
			num:      -0.4,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := round(tt.num)
			if result != tt.expected {
				t.Errorf("round(%v) = %v; expected %v", tt.num, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkCalculateAccessibilityScore(b *testing.B) {
	scores := []int{5, 6, 7, 8, 9, 10, 5, 6, 7, 8}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateAccessibilityScore(scores)
	}
}

func BenchmarkToFixed(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		toFixed(5.666666, 1)
	}
}
