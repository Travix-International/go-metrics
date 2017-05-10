package metrics // whitebox test

import (
	"testing"

	"github.com/Travix-International/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestMetrics_AddHistogram(t *testing.T) {
	log, _ := logger.New(map[string]string{})
	sut := NewMetrics("ns", log)

	result := sut.AddHistogram("s1", "n1", "h1")

	assert.Equal(t, "s1/n1", result.Key)
	assert.NotNil(t, result.hist)
	assert.NotNil(t, result.sum)
}

func TestMetrics_AddHistogramWithCustomBuckets(t *testing.T) {
	log, _ := logger.New(map[string]string{})
	sut := NewMetrics("ns", log)

	result := sut.AddHistogramWithCustomBuckets("s2", "n2", "h2", []float64{10, 20, 30.5})

	assert.Equal(t, "s2/n2", result.Key)
	assert.NotNil(t, result.hist)
	assert.NotNil(t, result.sum)

	result.Observe(25)
}

func BenchmarkMetrics_AddHistogram(b *testing.B) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	log, _ := logger.New(map[string]string{})
	sut := NewMetrics("ns", log)
	sut.AddHistogram("s3", "n3", "h3") // this is where we really create the histogram
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sut.AddHistogram("s3", "n3", "h3")
	}
}
