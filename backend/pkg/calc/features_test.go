package calc

import (
	"playability/types"
	"testing"
)

func TestCalculateFeatureScore(t *testing.T) {
	tests := []struct {
		name     string
		reports  []types.FeatureReport
		expected []types.FeatureStat
	}{
		{
			name:    "empty reports array",
			reports: []types.FeatureReport{},
			// When all percentages are 0, consensus is non-deterministic, so we just check structure
			expected: nil, // We'll handle this specially
		},
		{
			name: "single report - all true",
			reports: []types.FeatureReport{
				{ColorBlind: "true", ClosedCaptions: "true", FullControllerSupport: "true", ControllerRemapping: "true"},
			},
			expected: []types.FeatureStat{
				{FeatureName: "color_blind", Consensus: "true", SecondaryConsensus: "", TruePercentage: 1.0, LimitedPercentage: 0, FalsePercentage: 0},
				{FeatureName: "closed_captions", Consensus: "true", SecondaryConsensus: "", TruePercentage: 1.0, LimitedPercentage: 0, FalsePercentage: 0},
				{FeatureName: "full_controller_support", Consensus: "true", SecondaryConsensus: "", TruePercentage: 1.0, LimitedPercentage: 0, FalsePercentage: 0},
				{FeatureName: "controller_remapping", Consensus: "true", SecondaryConsensus: "", TruePercentage: 1.0, LimitedPercentage: 0, FalsePercentage: 0},
			},
		},
		{
			name: "single report - all false",
			reports: []types.FeatureReport{
				{ColorBlind: "false", ClosedCaptions: "false", FullControllerSupport: "false", ControllerRemapping: "false"},
			},
			expected: []types.FeatureStat{
				{FeatureName: "color_blind", Consensus: "false", SecondaryConsensus: "", TruePercentage: 0, LimitedPercentage: 0, FalsePercentage: 1.0},
				{FeatureName: "closed_captions", Consensus: "false", SecondaryConsensus: "", TruePercentage: 0, LimitedPercentage: 0, FalsePercentage: 1.0},
				{FeatureName: "full_controller_support", Consensus: "false", SecondaryConsensus: "", TruePercentage: 0, LimitedPercentage: 0, FalsePercentage: 1.0},
				{FeatureName: "controller_remapping", Consensus: "false", SecondaryConsensus: "", TruePercentage: 0, LimitedPercentage: 0, FalsePercentage: 1.0},
			},
		},
		{
			name: "multiple reports - clear consensus (true)",
			reports: []types.FeatureReport{
				{ColorBlind: "true", ClosedCaptions: "true", FullControllerSupport: "true", ControllerRemapping: "true"},
				{ColorBlind: "true", ClosedCaptions: "true", FullControllerSupport: "true", ControllerRemapping: "true"},
				{ColorBlind: "false", ClosedCaptions: "false", FullControllerSupport: "false", ControllerRemapping: "false"},
			},
			expected: []types.FeatureStat{
				{FeatureName: "color_blind", Consensus: "true", SecondaryConsensus: "false", TruePercentage: 0.67, LimitedPercentage: 0, FalsePercentage: 0.33},
				{FeatureName: "closed_captions", Consensus: "true", SecondaryConsensus: "false", TruePercentage: 0.67, LimitedPercentage: 0, FalsePercentage: 0.33},
				{FeatureName: "full_controller_support", Consensus: "true", SecondaryConsensus: "false", TruePercentage: 0.67, LimitedPercentage: 0, FalsePercentage: 0.33},
				{FeatureName: "controller_remapping", Consensus: "true", SecondaryConsensus: "false", TruePercentage: 0.67, LimitedPercentage: 0, FalsePercentage: 0.33},
			},
		},
		{
			name: "multiple reports - limited support",
			reports: []types.FeatureReport{
				{ColorBlind: "limited", ClosedCaptions: "limited", FullControllerSupport: "limited", ControllerRemapping: "limited"},
				{ColorBlind: "limited", ClosedCaptions: "limited", FullControllerSupport: "limited", ControllerRemapping: "limited"},
				{ColorBlind: "limited", ClosedCaptions: "limited", FullControllerSupport: "limited", ControllerRemapping: "limited"},
			},
			expected: []types.FeatureStat{
				{FeatureName: "color_blind", Consensus: "limited", SecondaryConsensus: "", TruePercentage: 0, LimitedPercentage: 1.0, FalsePercentage: 0},
				{FeatureName: "closed_captions", Consensus: "limited", SecondaryConsensus: "", TruePercentage: 0, LimitedPercentage: 1.0, FalsePercentage: 0},
				{FeatureName: "full_controller_support", Consensus: "limited", SecondaryConsensus: "", TruePercentage: 0, LimitedPercentage: 1.0, FalsePercentage: 0},
				{FeatureName: "controller_remapping", Consensus: "limited", SecondaryConsensus: "", TruePercentage: 0, LimitedPercentage: 1.0, FalsePercentage: 0},
			},
		},
		{
			name: "mixed support levels with clear secondary consensus",
			reports: []types.FeatureReport{
				{ColorBlind: "true", ClosedCaptions: "true", FullControllerSupport: "true", ControllerRemapping: "true"},
				{ColorBlind: "true", ClosedCaptions: "true", FullControllerSupport: "true", ControllerRemapping: "true"},
				{ColorBlind: "true", ClosedCaptions: "true", FullControllerSupport: "true", ControllerRemapping: "true"},
				{ColorBlind: "limited", ClosedCaptions: "limited", FullControllerSupport: "limited", ControllerRemapping: "limited"},
				{ColorBlind: "limited", ClosedCaptions: "limited", FullControllerSupport: "limited", ControllerRemapping: "limited"},
				{ColorBlind: "false", ClosedCaptions: "false", FullControllerSupport: "false", ControllerRemapping: "false"},
			},
			expected: []types.FeatureStat{
				{FeatureName: "color_blind", Consensus: "true", SecondaryConsensus: "limited", TruePercentage: 0.5, LimitedPercentage: 0.33, FalsePercentage: 0.17},
				{FeatureName: "closed_captions", Consensus: "true", SecondaryConsensus: "limited", TruePercentage: 0.5, LimitedPercentage: 0.33, FalsePercentage: 0.17},
				{FeatureName: "full_controller_support", Consensus: "true", SecondaryConsensus: "limited", TruePercentage: 0.5, LimitedPercentage: 0.33, FalsePercentage: 0.17},
				{FeatureName: "controller_remapping", Consensus: "true", SecondaryConsensus: "limited", TruePercentage: 0.5, LimitedPercentage: 0.33, FalsePercentage: 0.17},
			},
		},
		{
			name: "different features have different support levels",
			reports: []types.FeatureReport{
				{ColorBlind: "true", ClosedCaptions: "false", FullControllerSupport: "limited", ControllerRemapping: "true"},
				{ColorBlind: "true", ClosedCaptions: "false", FullControllerSupport: "limited", ControllerRemapping: "true"},
			},
			expected: []types.FeatureStat{
				{FeatureName: "color_blind", Consensus: "true", SecondaryConsensus: "", TruePercentage: 1.0, LimitedPercentage: 0, FalsePercentage: 0},
				{FeatureName: "closed_captions", Consensus: "false", SecondaryConsensus: "", TruePercentage: 0, LimitedPercentage: 0, FalsePercentage: 1.0},
				{FeatureName: "full_controller_support", Consensus: "limited", SecondaryConsensus: "", TruePercentage: 0, LimitedPercentage: 1.0, FalsePercentage: 0},
				{FeatureName: "controller_remapping", Consensus: "true", SecondaryConsensus: "", TruePercentage: 1.0, LimitedPercentage: 0, FalsePercentage: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateFeatureScore(tt.reports)

			// Special case for empty reports - just check structure
			if tt.expected == nil {
				if len(result) != 4 {
					t.Errorf("Expected 4 features, got %d", len(result))
				}
				for _, stat := range result {
					if stat.TruePercentage != 0 || stat.LimitedPercentage != 0 || stat.FalsePercentage != 0 {
						t.Errorf("Expected all percentages to be 0 for %s", stat.FeatureName)
					}
				}
				return
			}

			// Check we have the right number of features
			if len(result) != len(tt.expected) {
				t.Fatalf("Expected %d features, got %d", len(tt.expected), len(result))
			}

			// Check each feature
			for i, expectedStat := range tt.expected {
				actualStat := result[i]

				if actualStat.FeatureName != expectedStat.FeatureName {
					t.Errorf("Feature %d: expected name %s, got %s", i, expectedStat.FeatureName, actualStat.FeatureName)
				}
				if actualStat.Consensus != expectedStat.Consensus {
					t.Errorf("Feature %s: expected consensus %s, got %s", actualStat.FeatureName, expectedStat.Consensus, actualStat.Consensus)
				}
				if actualStat.SecondaryConsensus != expectedStat.SecondaryConsensus {
					t.Errorf("Feature %s: expected secondary consensus %s, got %s", actualStat.FeatureName, expectedStat.SecondaryConsensus, actualStat.SecondaryConsensus)
				}
				if actualStat.TruePercentage != expectedStat.TruePercentage {
					t.Errorf("Feature %s: expected true percentage %v, got %v", actualStat.FeatureName, expectedStat.TruePercentage, actualStat.TruePercentage)
				}
				if actualStat.LimitedPercentage != expectedStat.LimitedPercentage {
					t.Errorf("Feature %s: expected limited percentage %v, got %v", actualStat.FeatureName, expectedStat.LimitedPercentage, actualStat.LimitedPercentage)
				}
				if actualStat.FalsePercentage != expectedStat.FalsePercentage {
					t.Errorf("Feature %s: expected false percentage %v, got %v", actualStat.FeatureName, expectedStat.FalsePercentage, actualStat.FalsePercentage)
				}
			}
		})
	}
}

func TestFeatureNameToDisplayName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ColorBlind", "color_blind"},
		{"ClosedCaptions", "closed_captions"},
		{"FullControllerSupport", "full_controller_support"},
		{"ControllerRemapping", "controller_remapping"},
		{"UnknownFeature", "UnknownFeature"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := featureNameToDisplayName(tt.input)
			if result != tt.expected {
				t.Errorf("featureNameToDisplayName(%s) = %s; expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDetermineConsensus(t *testing.T) {
	tests := []struct {
		name                   string
		supportLevels          map[string]*SupportLevel
		expectedConsensus      string
		expectedSecondary      string
	}{
		{
			name: "clear winner - true",
			supportLevels: map[string]*SupportLevel{
				"true":    {Count: 10, Percentage: 0.8},
				"limited": {Count: 2, Percentage: 0.16},
				"false":   {Count: 0.5, Percentage: 0.04},
			},
			expectedConsensus: "true",
			expectedSecondary: "limited",
		},
		{
			name: "clear winner - false with no secondary",
			supportLevels: map[string]*SupportLevel{
				"true":    {Count: 0, Percentage: 0},
				"limited": {Count: 0, Percentage: 0},
				"false":   {Count: 10, Percentage: 1.0},
			},
			expectedConsensus: "false",
			expectedSecondary: "",
		},
		{
			name: "different clear winner and runner-up",
			supportLevels: map[string]*SupportLevel{
				"true":    {Count: 1, Percentage: 0.1},
				"limited": {Count: 3, Percentage: 0.3},
				"false":   {Count: 6, Percentage: 0.6},
			},
			expectedConsensus: "false",
			expectedSecondary: "limited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feature := &FeatureStats{
				FeatureName:   "test",
				SupportLevels: tt.supportLevels,
			}

			determineConsensus(feature)

			if feature.Consensus != tt.expectedConsensus {
				t.Errorf("Expected consensus %s, got %s", tt.expectedConsensus, feature.Consensus)
			}
			if feature.SecondaryConsensus != tt.expectedSecondary {
				t.Errorf("Expected secondary consensus %s, got %s", tt.expectedSecondary, feature.SecondaryConsensus)
			}
		})
	}
}

func TestCalculatePercentages(t *testing.T) {
	tests := []struct {
		name          string
		supportLevels map[string]*SupportLevel
		total         float64
		expected      map[string]float64
	}{
		{
			name: "normal distribution",
			supportLevels: map[string]*SupportLevel{
				"true":    {Count: 6},
				"limited": {Count: 3},
				"false":   {Count: 1},
			},
			total: 10,
			expected: map[string]float64{
				"true":    0.6,
				"limited": 0.3,
				"false":   0.1,
			},
		},
		{
			name: "zero total",
			supportLevels: map[string]*SupportLevel{
				"true":    {Count: 0},
				"limited": {Count: 0},
				"false":   {Count: 0},
			},
			total: 0,
			expected: map[string]float64{
				"true":    0,
				"limited": 0,
				"false":   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculatePercentages(tt.supportLevels, tt.total)

			for status, expectedPercentage := range tt.expected {
				if tt.supportLevels[status].Percentage != expectedPercentage {
					t.Errorf("Status %s: expected percentage %v, got %v",
						status, expectedPercentage, tt.supportLevels[status].Percentage)
				}
			}
		})
	}
}

func TestCalculateTotal(t *testing.T) {
	supportLevels := map[string]*SupportLevel{
		"true":    {Count: 5},
		"limited": {Count: 3},
		"false":   {Count: 2},
	}

	result := calculateTotal(supportLevels)
	expected := 10.0

	if result != expected {
		t.Errorf("calculateTotal() = %v; expected %v", result, expected)
	}
}

// Benchmark tests
func BenchmarkCalculateFeatureScore(b *testing.B) {
	reports := []types.FeatureReport{
		{ColorBlind: "true", ClosedCaptions: "true", FullControllerSupport: "limited", ControllerRemapping: "false"},
		{ColorBlind: "true", ClosedCaptions: "false", FullControllerSupport: "limited", ControllerRemapping: "true"},
		{ColorBlind: "limited", ClosedCaptions: "true", FullControllerSupport: "true", ControllerRemapping: "true"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateFeatureScore(reports)
	}
}
