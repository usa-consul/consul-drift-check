package diff

import (
	"testing"
	"time"
)

func makeForecastTrend(counts []int) []TrendPoint {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	points := make([]TrendPoint, len(counts))
	for i, c := range counts {
		points[i] = TrendPoint{
			At:    base.Add(time.Duration(i) * 24 * time.Hour),
			Total: c,
		}
	}
	return points
}

func TestForecast_TooFewPoints_ReturnsNil(t *testing.T) {
	result := Forecast(makeForecastTrend([]int{5}), ForecastOptions{Steps: 3})
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestForecast_ReturnsExpectedStepCount(t *testing.T) {
	trend := makeForecastTrend([]int{2, 4, 6, 8, 10})
	result := Forecast(trend, ForecastOptions{Steps: 3, Interval: 24 * time.Hour})
	if len(result) != 3 {
		t.Fatalf("expected 3 points, got %d", len(result))
	}
}

func TestForecast_PointsAreChronological(t *testing.T) {
	trend := makeForecastTrend([]int{1, 2, 3, 4, 5})
	result := Forecast(trend, ForecastOptions{Steps: 4, Interval: 24 * time.Hour})
	for i := 1; i < len(result); i++ {
		if !result[i].At.After(result[i-1].At) {
			t.Errorf("points not chronological at index %d", i)
		}
	}
}

func TestForecast_FlatTrend_CountStaysConstant(t *testing.T) {
	trend := makeForecastTrend([]int{10, 10, 10, 10, 10})
	result := Forecast(trend, ForecastOptions{Steps: 2, Interval: 24 * time.Hour})
	for _, p := range result {
		if p.Count < 9.0 || p.Count > 11.0 {
			t.Errorf("expected ~10, got %f", p.Count)
		}
	}
}

func TestForecast_DefaultSteps_ThreePoints(t *testing.T) {
	trend := makeForecastTrend([]int{1, 3, 5})
	result := Forecast(trend, ForecastOptions{})
	if len(result) != 3 {
		t.Fatalf("expected default 3 steps, got %d", len(result))
	}
}

func TestForecast_NeverNegative(t *testing.T) {
	trend := makeForecastTrend([]int{100, 50, 10, 5, 1})
	result := Forecast(trend, ForecastOptions{Steps: 5, Interval: 24 * time.Hour})
	for _, p := range result {
		if p.Count < 0 {
			t.Errorf("forecast count should not be negative, got %f", p.Count)
		}
	}
}
