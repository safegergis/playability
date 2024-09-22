package calc

import (
	"playability/types"
	"slices"
)

// SupportLevel represents the count and percentage of a support level
type SupportLevel struct {
	Count      float64
	Percentage float64
}

// FeatureStats holds the support levels for a feature
type FeatureStats struct {
	FeatureName        string
	SupportLevels      map[string]*SupportLevel
	Consensus          string
	SecondaryConsensus string
}

// CalculateFeatureScore calculates the feature score for a game
func CalculateFeatureScore(reports []types.FeatureReport) []types.Feature {
	features := []string{"ColorBlind", "ClosedCaptions", "FullControllerSupport", "ControllerRemapping"}
	supportStatuses := []string{"true", "limited", "false"}

	featureStatsMap := initializeFeatureStats(features, supportStatuses)

	// Count support levels
	for _, report := range reports {
		updateSupportCount(featureStatsMap["ColorBlind"], report.ColorBlind)
		updateSupportCount(featureStatsMap["ClosedCaptions"], report.ClosedCaptions)
		updateSupportCount(featureStatsMap["FullControllerSupport"], report.FullControllerSupport)
		updateSupportCount(featureStatsMap["ControllerRemapping"], report.ControllerRemapping)
	}

	// Calculate percentages and determine consensus
	var result []types.Feature
	for _, featureName := range features {
		stats := featureStatsMap[featureName]
		total := calculateTotal(stats.SupportLevels)
		calculatePercentages(stats.SupportLevels, total)
		determineConsensus(stats)
		result = append(result, types.Feature{
			FeatureName:        featureNameToDisplayName(featureName),
			TruePercentage:     toFixed(stats.SupportLevels["true"].Percentage, 2),
			LimitedPercentage:  toFixed(stats.SupportLevels["limited"].Percentage, 2),
			FalsePercentage:    toFixed(stats.SupportLevels["false"].Percentage, 2),
			Consensus:          stats.Consensus,
			SecondaryConsensus: stats.SecondaryConsensus,
		})
	}

	return result
}

// initializeFeatureStats initializes the feature statistics map
func initializeFeatureStats(features, supportStatuses []string) map[string]*FeatureStats {
	featureStatsMap := make(map[string]*FeatureStats)
	for _, feature := range features {
		supportMap := make(map[string]*SupportLevel)
		for _, status := range supportStatuses {
			supportMap[status] = &SupportLevel{Count: 0.0, Percentage: 0.0}
		}
		featureStatsMap[feature] = &FeatureStats{
			FeatureName:   feature,
			SupportLevels: supportMap,
		}
	}
	return featureStatsMap
}

// updateSupportCount increments the count for the given support level
func updateSupportCount(feature *FeatureStats, status string) {
	if level, exists := feature.SupportLevels[status]; exists {
		level.Count++
	}
}

// calculateTotal sums up all counts for a feature
func calculateTotal(supportLevels map[string]*SupportLevel) float64 {
	total := 0.0
	for _, level := range supportLevels {
		total += level.Count
	}
	return total
}

// calculatePercentages computes the percentage for each support level
func calculatePercentages(supportLevels map[string]*SupportLevel, total float64) {
	for _, level := range supportLevels {
		if total > 0 {
			level.Percentage = level.Count / total
		}
	}
}

// determineConsensus finds the highest and second highest support levels
func determineConsensus(feature *FeatureStats) {
	type levelPercentage struct {
		Status     string
		Percentage float64
	}

	var levels []levelPercentage
	for status, level := range feature.SupportLevels {
		levels = append(levels, levelPercentage{Status: status, Percentage: level.Percentage})
	}

	slices.SortFunc(levels, func(a, b levelPercentage) int {
		if a.Percentage > b.Percentage {
			return -1
		}
		if a.Percentage < b.Percentage {
			return 1
		}
		return 0
	})

	feature.Consensus = levels[0].Status
	if len(levels) > 1 && levels[1].Percentage > 0 {
		feature.SecondaryConsensus = levels[1].Status
	} else {
		feature.SecondaryConsensus = ""
	}
}

// featureNameToDisplayName converts feature keys to display names
func featureNameToDisplayName(name string) string {
	switch name {
	case "ColorBlind":
		return "color_blind"
	case "ClosedCaptions":
		return "closed_captions"
	case "FullControllerSupport":
		return "full_controller_support"
	case "ControllerRemapping":
		return "controller_remapping"
	default:
		return name
	}
}
