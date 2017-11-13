package metrics // white-box test

import (
	"testing"
	"time"

	"github.com/Travix-International/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestMetrics_AddHistogramWithCustomBuckets(t *testing.T) {
	log, _ := logger.New(make(map[string]string))
	sut := NewMetrics("ns", log)

	result := sut.AddHistogramWithCustomBuckets("s2", "n2", "h2", []float64{10, 20, 30.5})

	assert.Equal(t, "s2/n2", result.Key)
	assert.NotNil(t, result.hist)
	assert.NotNil(t, result.sum)

	result.Observe(25)
}

func TestMetrics_AddHistogramVecWithCustomBuckets(t *testing.T) {
	log, _ := logger.New(make(map[string]string))
	sut := NewMetrics("ns", log)
	labels := []string{"lbl1", "lbl2", "lbl3"}
	values := []string{"val1", "val2", "val3"}

	result := sut.AddHistogramVecWithCustomBuckets("s3", "n3", "h3", labels, values,
		[]float64{10, 20, 30.5})

	assert.Equal(t, "s3/n3", result.Key)
	assert.NotNil(t, result.histVec)

	result.Observe(25)
}

func TestMetrics_AddSummaryVecWithCustomObjectives(t *testing.T) {
	log, _ := logger.New(make(map[string]string))
	sut := NewMetrics("ns", log)
	labels := []string{"lbl1", "lbl2", "lbl3"}
	values := []string{"val1", "val2", "val3"}

	result := sut.AddSummaryVecWithCustomObjectives("s4", "n4", "h4", labels, values,
		map[float64]float64{0.5: 0.05, 0.75: 0.025, 0.9: 0.01, 0.99: 0.001})

	assert.Equal(t, "s4/n4", result.Key)
	assert.NotNil(t, result.summaryVec)
}

func TestMetrics_Count(t *testing.T) {
	sys := "count"
	log, _ := logger.New(make(map[string]string))
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
	log, _ := logger.New(make(map[string]string))
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
	log, _ := logger.New(make(map[string]string))
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

func TestMetrics_AddHistogram(t *testing.T) {
	log, _ := logger.New(make(map[string]string))
	sut := NewMetrics("ns", log)
	sys := "addhist"
	start := time.Now().Add(+1 * time.Second)

	hist := sut.AddHistogram(sys, sys, sys)

	// Act
	hist.RecordTimeElapsed(start)
	hist.RecordDuration(start, time.Millisecond)

	assert.Equal(t, "addhist/addhist", hist.Key)
	assert.NotNil(t, hist.hist)
	assert.NotNil(t, hist.sum)
}

func TestMetrics_AddHistogramVec(t *testing.T) {
	log, _ := logger.New(make(map[string]string))
	sut := NewMetrics("ns", log)
	sys := "addhistvec"
	start := time.Now().Add(+1 * time.Second)
	labels := []string{"label1", "label2", "label3"}
	values := []string{"val1", "val2", "val3"}

	vec := sut.AddHistogramVec(sys, sys, sys, labels, values)

	// Act
	vec.RecordTimeElapsed(start)
	vec.RecordDuration(start, time.Millisecond)

	assert.Equal(t, "addhistvec/addhistvec", vec.Key)
	assert.Equal(t, labels, vec.Labels)
	assert.NotNil(t, values, vec.LabelValues)
}

func TestMetrics_AddSummaryVec(t *testing.T) {
	log, _ := logger.New(make(map[string]string))
	sut := NewMetrics("ns", log)
	sys := "addsumvec"
	start := time.Now().Add(+1 * time.Second)
	labels := []string{"label1", "label2", "label3"}
	values := []string{"val1", "val2", "val3"}

	vec := sut.AddSummaryVec(sys, sys, sys, labels, values)

	// Act
	vec.RecordTimeElapsed(start)
	vec.RecordDuration(start, time.Microsecond)

	assert.Equal(t, "addsumvec/addsumvec", vec.Key)
	assert.Equal(t, labels, vec.Labels)
	assert.NotNil(t, values, vec.LabelValues)
}

/* Benchmarks */

func BenchmarkMetrics_AddHistogram(b *testing.B) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	log, _ := logger.New(make(map[string]string))
	sut := NewMetrics("ns", log)
	sut.AddHistogram("s3", "n3", "h3") // this is where we really create the histogram
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sut.AddHistogram("s3", "n3", "h3")
	}
}
