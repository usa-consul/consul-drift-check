package diff

import (
	"testing"
	"time"
)

func BenchmarkForecast_LargeHistory(b *testing.B) {
	base := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	trend := make([]TrendPoint, 365)
	for i := range trend {
		trend[i] = TrendPoint{
			At:    base.Add(time.Duration(i) * 24 * time.Hour),
			Total: i*2 + 1,
		}
	}
	opts := ForecastOptions{Steps: 30, Interval: 24 * time.Hour}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Forecast(trend, opts)
	}
}
