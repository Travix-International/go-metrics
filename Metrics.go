package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/Travix-International/logger"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	// Metrics provides a set of convenience functions that wrap Prometheus
	Metrics struct {
		Namespace       string
		Counters        map[string]prometheus.Counter
		CounterVecs     map[string]*prometheus.CounterVec
		Summaries       map[string]prometheus.Summary
		SummaryVecs     map[string]*prometheus.SummaryVec
		Histograms      map[string]prometheus.Histogram
		HistogramVecs   map[string]*prometheus.HistogramVec
		Gauges          map[string]prometheus.Gauge
		Logger          *logger.Logger
		countMutex      *sync.RWMutex
		countVecMutex   *sync.RWMutex
		histMutex       *sync.RWMutex
		histVecMutex    *sync.RWMutex
		summaryVecMutex *sync.RWMutex
		gaugeMutex      *sync.RWMutex
	}

	// MetricsHistogram combines a histogram and summary
	MetricsHistogram struct {
		Key  string
		hist prometheus.Histogram
		sum  prometheus.Summary
	}

	// HistogramVec wraps prometheus.HistogramVec
	HistogramVec struct {
		Key         string
		Labels      []string
		LabelValues []string
		histVec     *prometheus.HistogramVec
	}

	// SummaryVec wraps prometheus.SummaryVec
	SummaryVec struct {
		Key         string
		Labels      []string
		LabelValues []string
		summaryVec  *prometheus.SummaryVec
	}
)

// NewMetrics will instantiate a new Metrics wrapper object
func NewMetrics(namespace string, logger *logger.Logger) *Metrics {
	m := Metrics{
		Namespace:       namespace,
		Logger:          logger,
		Counters:        make(map[string]prometheus.Counter),
		CounterVecs:     make(map[string]*prometheus.CounterVec),
		Histograms:      make(map[string]prometheus.Histogram),
		HistogramVecs:   make(map[string]*prometheus.HistogramVec),
		Summaries:       make(map[string]prometheus.Summary),
		SummaryVecs:     make(map[string]*prometheus.SummaryVec),
		Gauges:          make(map[string]prometheus.Gauge),
		countMutex:      &sync.RWMutex{},
		countVecMutex:   &sync.RWMutex{},
		histMutex:       &sync.RWMutex{},
		histVecMutex:    &sync.RWMutex{},
		summaryVecMutex: &sync.RWMutex{},
		gaugeMutex:      &sync.RWMutex{},
	}
	return &m
}

// DefaultObjectives returns a default map of quantiles to be used in summaries.
func DefaultObjectives() map[float64]float64 {
	return map[float64]float64{0.5: 0.05, 0.75: 0.025, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001, 0.999: 0.0001}
}

// Count increases the counter for the specified subsystem and name.
func (m *Metrics) Count(subsystem, name, help string) {
	m.countMutex.RLock()
	key := fmt.Sprintf("%s/%s", subsystem, name)
	counter, exists := m.Counters[key]
	m.countMutex.RUnlock()

	if !exists {
		m.countMutex.Lock()
		if counter, exists = m.Counters[key]; !exists {
			counter = prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: m.Namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
			})
			m.Counters[key] = counter
			err := prometheus.Register(counter)
			if err != nil {
				m.Logger.Warn("MetricsCounterRegistrationFailed",
					fmt.Sprintf("CounterHandler: Counter registration %v failed: %v", counter, err))
			}
		}
		m.countMutex.Unlock()
	}

	counter.Inc()
}

// SetGauge sets the gauge value for the specified subsystem and name.
func (m *Metrics) SetGauge(value float64, subsystem, name, help string) {
	m.gaugeMutex.RLock()
	key := fmt.Sprintf("%s/%s", subsystem, name)
	gauge, exists := m.Gauges[key]
	m.gaugeMutex.RUnlock()

	if !exists {
		m.gaugeMutex.Lock()
		if gauge, exists = m.Gauges[key]; !exists {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: m.Namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
			})
			m.Gauges[key] = gauge
			err := prometheus.Register(gauge)
			if err != nil {
				m.Logger.Warn("MetricsSetGaugeFailed",
					fmt.Sprintf("SetGauge: Gauge registration %v failed: %v", gauge, err))
			}
		}
		m.gaugeMutex.Unlock()
	}

	gauge.Set(value)
}

// CountLabels increases the counter for the specified subsystem and name and adds the specified labels with values.
func (m *Metrics) CountLabels(subsystem, name, help string, labels, values []string) {
	m.countVecMutex.RLock()
	key := fmt.Sprintf("%s/%s", subsystem, name)
	counter, exists := m.CounterVecs[key]
	m.countVecMutex.RUnlock()

	if !exists {
		m.countVecMutex.Lock()
		if counter, exists = m.CounterVecs[key]; !exists {
			counter = prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: m.Namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
			}, labels)
			m.CounterVecs[key] = counter
			err := prometheus.Register(counter)
			if err != nil {
				m.Logger.Warn("MetricsCounterLabelRegistrationFailed",
					fmt.Sprintf("CounterLabelHandler: Counter registration %v failed: %v", counter, err))
			}
		}
		m.countVecMutex.Unlock()
	}

	counter.WithLabelValues(values...).Inc()
}

// IncreaseCounter increases the counter for the specified subsystem and name with the specified increment.
func (m *Metrics) IncreaseCounter(subsystem, name, help string, increment int) {
	m.countMutex.RLock()
	key := fmt.Sprintf("%s/%s", subsystem, name)
	counter, exists := m.Counters[key]
	m.countMutex.RUnlock()

	if !exists {
		m.countMutex.Lock()
		if counter, exists = m.Counters[key]; !exists {
			counter = prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: m.Namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
			})
			m.Counters[key] = counter
			err := prometheus.Register(counter)
			if err != nil {
				m.Logger.Warn("MetricsIncreaseCounterRegistrationFailed",
					fmt.Sprintf("CounterHandler: Counter registration failed: %v: %v", counter, err))
			}
		}
		m.countMutex.Unlock()
	}

	counter.Add(float64(increment))
}

// AddHistogram returns the MetricsHistogram for the specified subsystem and name.
func (m *Metrics) AddHistogram(subsystem, name, help string) *MetricsHistogram {
	return m.addHistogramWithBuckets(subsystem, name, help, prometheus.DefBuckets)
}

// AddHistogramWithCustomBuckets returns the MetricsHistogram for the specified subsystem and name with the specified buckets.
func (m *Metrics) AddHistogramWithCustomBuckets(subsystem, name, help string, buckets []float64) *MetricsHistogram {
	return m.addHistogramWithBuckets(subsystem, name, help, buckets)
}

func (m *Metrics) addHistogramWithBuckets(subsystem, name, help string, buckets []float64) *MetricsHistogram {
	m.histMutex.RLock()
	key := fmt.Sprintf("%s/%s", subsystem, name)
	sum, exists := m.Summaries[key]
	hist := m.Histograms[key]
	m.histMutex.RUnlock()

	if !exists {
		m.histMutex.Lock()
		if sum, exists = m.Summaries[key]; !exists {
			// todo: remove Summary creation/observation
			sum = prometheus.NewSummary(prometheus.SummaryOpts{
				Namespace:  m.Namespace,
				Subsystem:  subsystem,
				Name:       name + "_summary",
				Help:       help,
				Objectives: DefaultObjectives(),
			})
			prometheus.MustRegister(sum)
			m.Summaries[key] = sum

			hist = prometheus.NewHistogram(prometheus.HistogramOpts{
				Namespace: m.Namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
				Buckets:   buckets,
			})
			prometheus.MustRegister(hist)
			m.Histograms[key] = hist
		}
		m.histMutex.Unlock()
	}

	mh := MetricsHistogram{
		Key:  key,
		hist: hist,
		sum:  sum,
	}
	return &mh
}

// AddHistogramVec returns the HistogramVec for the specified subsystem and name.
func (m *Metrics) AddHistogramVec(subsystem, name, help string, labels, labelValues []string) *HistogramVec {
	return m.addHistogramVecWithBuckets(subsystem, name, help, labels, labelValues, prometheus.DefBuckets)
}

// AddHistogramVecWithCustomBuckets returns the HistogramVec for the specified subsystem and name with the specified buckets.
func (m *Metrics) AddHistogramVecWithCustomBuckets(subsystem, name, help string, labels, labelValues []string,
	buckets []float64) *HistogramVec {

	return m.addHistogramVecWithBuckets(subsystem, name, help, labels, labelValues, buckets)
}

func (m *Metrics) addHistogramVecWithBuckets(subsystem, name, help string, labels, labelValues []string,
	buckets []float64) *HistogramVec {

	m.histVecMutex.RLock()
	key := fmt.Sprintf("%s/%s", subsystem, name)
	vec, exists := m.HistogramVecs[key]
	m.histVecMutex.RUnlock()

	if !exists {
		m.histVecMutex.Lock()
		vec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: m.Namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		}, labels)
		prometheus.MustRegister(vec)
		m.HistogramVecs[key] = vec
		m.histVecMutex.Unlock()
	}

	mh := HistogramVec{
		Key:         key,
		Labels:      labels,
		LabelValues: labelValues,
		histVec:     vec,
	}
	return &mh
}

// AddSummaryVec returns the SummaryVec for the specified subsystem and name.
func (m *Metrics) AddSummaryVec(subsystem, name, help string, labels, labelValues []string) *SummaryVec {
	return m.addSummaryVecWithObjectives(subsystem, name, help, labels, labelValues, DefaultObjectives())
}

// AddSummaryVecWithCustomObjectives returns the SummaryVec for the specified subsystem and name with the specified objectives.
func (m *Metrics) AddSummaryVecWithCustomObjectives(subsystem, name, help string, labels, labelValues []string,
	objectives map[float64]float64) *SummaryVec {

	return m.addSummaryVecWithObjectives(subsystem, name, help, labels, labelValues, objectives)
}

func (m *Metrics) addSummaryVecWithObjectives(subsystem, name, help string, labels, labelValues []string,
	objectives map[float64]float64) *SummaryVec {

	m.summaryVecMutex.RLock()
	key := fmt.Sprintf("%s/%s", subsystem, name)
	vec, exists := m.SummaryVecs[key]
	m.summaryVecMutex.RUnlock()

	if !exists {
		m.summaryVecMutex.Lock()
		vec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace:  m.Namespace,
			Subsystem:  subsystem,
			Name:       name + "_summary",
			Help:       help,
			Objectives: objectives,
		}, labels)
		prometheus.MustRegister(vec)
		m.SummaryVecs[key] = vec
		m.summaryVecMutex.Unlock()
	}

	mh := SummaryVec{
		Key:         key,
		Labels:      labels,
		LabelValues: labelValues,
		summaryVec:  vec,
	}
	return &mh
}

// RecordTimeElapsed adds the elapsed time since the specified start to the histogram in seconds and to the linked
// summary in milliseconds.
func (histogram *MetricsHistogram) RecordTimeElapsed(start time.Time) {
	elapsed := float64(time.Since(start).Seconds())
	histogram.hist.Observe(elapsed)         // The default histogram buckets are recorded in seconds
	histogram.sum.Observe(elapsed * 1000.0) // While we have summaries in milliseconds
}

// RecordDuration adds the elapsed time since the specified start to the histogram in the specified unit of time
// and to the linked summary in milliseconds.
func (histogram *MetricsHistogram) RecordDuration(start time.Time, unit time.Duration) {
	since := time.Since(start)
	elapsedSeconds := float64(since.Seconds())
	elapsedUnits := float64(since.Truncate(unit))

	histogram.hist.Observe(elapsedUnits)
	histogram.sum.Observe(elapsedSeconds * 1000.0)
}

// Observe adds the specified value to the histogram.
func (histogram *MetricsHistogram) Observe(value float64) {
	histogram.hist.Observe(value)
}

// RecordTimeElapsed adds the elapsed time since the specified start to the histogram in seconds.
func (vec *HistogramVec) RecordTimeElapsed(start time.Time) {
	elapsed := float64(time.Since(start).Seconds())
	vec.Observe(elapsed) // The default histogram buckets are recorded in seconds
}

// RecordDuration adds the elapsed time since the specified start to the histogram in the specified unit of time.
func (vec *HistogramVec) RecordDuration(start time.Time, unit time.Duration) {
	since := time.Since(start)
	elapsedUnits := float64(since.Truncate(unit))

	vec.Observe(elapsedUnits)
}

// Observe adds the specified value to the histogram.
func (vec *HistogramVec) Observe(value float64) {
	vec.histVec.WithLabelValues(vec.LabelValues...).Observe(value)
}

// RecordTimeElapsed adds the elapsed time since the specified start to the summary in milliseconds.
func (vec *SummaryVec) RecordTimeElapsed(start time.Time) {
	elapsed := float64(time.Since(start).Seconds())
	vec.summaryVec.WithLabelValues(vec.LabelValues...).Observe(elapsed * 1000.0) // Summaries are in milliseconds
}
