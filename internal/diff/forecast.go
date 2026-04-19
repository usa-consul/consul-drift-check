package diff

import (
	"math"
	"sort"
	"time"
)

// ForecastPoint represents a predicted drift count at a future time.
type ForecastPoint struct {
	At    time.Time
	Count float64
}

// ForecastOptions controls how the forecast is generated.
type ForecastOptions struct {
	// Steps is the number of future intervals to predict.
	Steps int
	// Interval is the duration between each forecast point.
	Interval time.Duration
}

// Forecast uses simple linear regression over historical trend points
// to predict future drift counts.
func Forecast(trend []TrendPoint, opts ForecastOptions) []ForecastPoint {
	if opts.Steps <= 0 {
		opts.Steps = 3
	}
	if opts.Interval <= 0 {
		opts.Interval = 24 * time.Hour
	}
	if len(trend) < 2 {
		return nil
	}

	sorted := make([]TrendPoint, len(trend))
	copy(sorted, trend)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].At.Before(sorted[j].At)
	})

	origin := sorted[0].At
	n := float64(len(sorted))
	var sumX, sumY, sumXY, sumX2 float64
	for _, p := range sorted {
		x := p.At.Sub(origin).Hours()
		y := float64(p.Total)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	denom := n*sumX2 - sumX*sumX
	var slope, intercept float64
	if denom != 0 {
		slope = (n*sumXY - sumX*sumY) / denom
		intercept = (sumY - slope*sumX) / n
	}

	last := sorted[len(sorted)-1].At
	points := make([]ForecastPoint, opts.Steps)
	for i := 0; i < opts.Steps; i++ {
		at := last.Add(time.Duration(i+1) * opts.Interval)
		x := at.Sub(origin).Hours()
		predicted := intercept + slope*x
		points[i] = ForecastPoint{
			At:    at,
			Count: math.Max(0, predicted),
		}
	}
	return points
}
