package metrics // white-box test

import (
	"testing"
	"time"

	"github.com/Travix-International/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestMetrics_AddHistogram(t *testing.T) {
	log := logger.New()
	sut := NewMetrics("ns", log)

	result := sut.AddHistogram("s1", "n1", "h1")

	assert.Equal(t, "s1/n1", result.Key)
	assert.NotNil(t, result.hist)
	assert.NotNil(t, result.sum)
}

func TestMetrics_AddHistogramWithCustomBuckets(t *testing.T) {
	log := logger.New()
	sut := NewMetrics("ns", log)

	result := sut.AddHistogramWithCustomBuckets("s2", "n2", "h2", []float64{10, 20, 30.5})

	assert.Equal(t, "s2/n2", result.Key)
	assert.NotNil(t, result.hist)
	assert.NotNil(t, result.sum)

	result.Observe(25)
}

func TestMetrics_Count(t *testing.T) {
	sys := "count"
	log := logger.New()
	sut := NewMetrics(sys, log)

	sut.Count(sys, sys, sys)
	assert.Equal(t, 1, len(sut.Counters))

	sut.Count("x", sys, sys)
	assert.Equal(t, 2, len(sut.Counters))

	sut.IncreaseCounter("x", sys, sys, 5)
	assert.Equal(t, 2, len(sut.Counters))

	sut.Count(sys, "x", sys)
	assert.Equal(t, 3, len(sut.Counters))

	sut.Count(sys, sys, sys)
	assert.Equal(t, 3, len(sut.Counters))

	sut.IncreaseCounter("y", sys, sys, 3)
	assert.Equal(t, 4, len(sut.Counters))
}

func TestMetrics_SetGauge(t *testing.T) {
	sys := "gauge"
	log := logger.New()
	sut := NewMetrics(sys, log)

	sut.SetGauge(5, sys, sys, sys)
	assert.Equal(t, 1, len(sut.Gauges))

	sut.SetGauge(6, "x", sys, sys)
	assert.Equal(t, 2, len(sut.Gauges))

	sut.SetGauge(7, sys, "x", sys)
	assert.Equal(t, 3, len(sut.Gauges))

	sut.SetGauge(8, sys, sys, sys)
	assert.Equal(t, 3, len(sut.Gauges))
}

func TestMetrics_CountLabels(t *testing.T) {
	sys := "labels"
	log := logger.New()
	sut := NewMetrics(sys, log)
	labels := []string{"lbl1", "lbl2", "lbl3"}
	values := []string{"val1", "val2", "val3"}

	sut.CountLabels(sys, sys, sys, labels, values)
	assert.Equal(t, 1, len(sut.CounterVecs))

	sut.CountLabels("x", sys, sys, labels, values)
	assert.Equal(t, 2, len(sut.CounterVecs))

	sut.CountLabels(sys, "x", sys, labels, values)
	assert.Equal(t, 3, len(sut.CounterVecs))

	sut.CountLabels(sys, sys, sys, labels, values)
	assert.Equal(t, 3, len(sut.CounterVecs))
}

func TestMetrics_RecordTimeElapsed(t *testing.T) {
	log := logger.New()
	sut := NewMetrics("ns", log)
	sys := "elapsed"
	start := time.Now().Add(+1 * time.Second)

	hist := sut.AddHistogram(sys, sys, sys)

	// Act
	hist.RecordTimeElapsed(start)
	hist.RecordDuration(start, time.Millisecond)

	// Alas, nothing to assert!
}

/* Benchmarks */

func BenchmarkMetrics_AddHistogram(b *testing.B) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	log := logger.New()
	sut := NewMetrics("ns", log)
	sut.AddHistogram("s3", "n3", "h3") // this is where we really create the histogram
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sut.AddHistogram("s3", "n3", "h3")
	}
}
