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
		Namespace     string
		Counters      map[string]prometheus.Counter
		CounterVecs   map[string]*prometheus.CounterVec
		Summaries     map[string]prometheus.Summary
		Histograms    map[string]prometheus.Histogram
		Gauges        map[string]prometheus.Gauge
		Logger        *logger.Logger
		countMutex    *sync.RWMutex
		countVecMutex *sync.RWMutex
		histMutex     *sync.RWMutex
		gaugeMutex    *sync.RWMutex
	}

	// MetricsHistogram combines a histogram and summary
	MetricsHistogram struct {
		Key  string
		hist prometheus.Histogram
		sum  prometheus.Summary
	}
)

// NewMetrics will instantiate a new Metrics wrapper object
func NewMetrics(namespace string, logger *logger.Logger) *Metrics {
	m := Metrics{
		Namespace:     namespace,
		Logger:        logger,
		Counters:      make(map[string]prometheus.Counter),
		CounterVecs:   make(map[string]*prometheus.CounterVec),
		Histograms:    make(map[string]prometheus.Histogram),
		Summaries:     make(map[string]prometheus.Summary),
		Gauges:        make(map[string]prometheus.Gauge),
		countMutex:    &sync.RWMutex{},
		countVecMutex: &sync.RWMutex{},
		histMutex:     &sync.RWMutex{},
		gaugeMutex:    &sync.RWMutex{},
	}
	return &m
}

// Count increases the counter for the specified subsystem and name.
func (ctx *Metrics) Count(subsystem, name, help string) {
	ctx.countMutex.RLock()
	key := fmt.Sprintf("%s/%s", subsystem, name)
	counter, exists := ctx.Counters[key]
	ctx.countMutex.RUnlock()

	if !exists {
		ctx.countMutex.Lock()
		if counter, exists = ctx.Counters[key]; !exists {
			counter = prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: ctx.Namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
			})
			ctx.Counters[key] = counter
			err := prometheus.Register(counter)
			if err != nil {
				ctx.Logger.Warn("MetricsCounterRegistrationFailed",
					fmt.Sprintf("CounterHandler: Counter registration %v failed: %v", counter, err))
			}
		}
		ctx.countMutex.Unlock()
	}

	counter.Inc()
}

// SetGauge sets the gauge value for the specified subsystem and name.
func (ctx *Metrics) SetGauge(value float64, subsystem, name, help string) {
	ctx.gaugeMutex.RLock()
	key := fmt.Sprintf("%s/%s", subsystem, name)
	gauge, exists := ctx.Gauges[key]
	ctx.gaugeMutex.RUnlock()

	if !exists {
		ctx.gaugeMutex.Lock()
		if gauge, exists = ctx.Gauges[key]; !exists {
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: ctx.Namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
			})
			ctx.Gauges[key] = gauge
			err := prometheus.Register(gauge)
			if err != nil {
				ctx.Logger.Warn("MetricsSetGaugeFailed",
					fmt.Sprintf("SetGauge: Gauge registration %v failed: %v", gauge, err))
			}
		}
		ctx.gaugeMutex.Unlock()
	}

	gauge.Set(value)
}

// CountLabels increases the counter for the specified subsystem and name and adds the specified labels with values.
func (ctx *Metrics) CountLabels(subsystem, name, help string, labels, values []string) {
	ctx.countVecMutex.RLock()
	key := fmt.Sprintf("%s/%s", subsystem, name)
	counter, exists := ctx.CounterVecs[key]
	ctx.countVecMutex.RUnlock()

	if !exists {
		ctx.countVecMutex.Lock()
		if counter, exists = ctx.CounterVecs[key]; !exists {
			counter = prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: ctx.Namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
			}, labels)
			ctx.CounterVecs[key] = counter
			err := prometheus.Register(counter)
			if err != nil {
				ctx.Logger.Warn("MetricsCounterLabelRegistrationFailed",
					fmt.Sprintf("CounterLabelHandler: Counter registration %v failed: %v", counter, err))
			}
		}
		ctx.countVecMutex.Unlock()
	}

	counter.WithLabelValues(values...).Inc()
}

// IncreaseCounter increases the counter for the specified subsystem and name with the specified increment.
func (ctx *Metrics) IncreaseCounter(subsystem, name, help string, increment int) {
	ctx.countMutex.RLock()
	key := fmt.Sprintf("%s/%s", subsystem, name)
	counter, exists := ctx.Counters[key]
	ctx.countMutex.RUnlock()

	if !exists {
		ctx.countMutex.Lock()
		if counter, exists = ctx.Counters[key]; !exists {
			counter = prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: ctx.Namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
			})
			ctx.Counters[key] = counter
			err := prometheus.Register(counter)
			if err != nil {
				ctx.Logger.Warn("MetricsIncreaseCounterRegistrationFailed",
					fmt.Sprintf("CounterHandler: Counter registration failed: %v: %v", counter, err))
			}
		}
		ctx.countMutex.Unlock()
	}

	counter.Add(float64(increment))
}

// AddHistogram returns the MetricsHistogram for the specified subsystem and name.
func (ctx *Metrics) AddHistogram(subsystem, name, help string) *MetricsHistogram {
	return ctx.addHistogramWithBuckets(subsystem, name, help, prometheus.DefBuckets)
}

// AddHistogramWithCustomBuckets returns the MetricsHistogram for the specified subsystem and name with the specified buckets.
func (ctx *Metrics) AddHistogramWithCustomBuckets(subsystem, name, help string, buckets []float64) *MetricsHistogram {
	return ctx.addHistogramWithBuckets(subsystem, name, help, buckets)
}

func (ctx *Metrics) addHistogramWithBuckets(subsystem, name, help string, buckets []float64) *MetricsHistogram {
	ctx.histMutex.RLock()
	key := fmt.Sprintf("%s/%s", subsystem, name)
	sum, exists := ctx.Summaries[key]
	hist := ctx.Histograms[key]
	ctx.histMutex.RUnlock()

	if !exists {
		ctx.histMutex.Lock()
		if sum, exists = ctx.Summaries[key]; !exists {
			// todo: remove Summary creation/observation
			sum = prometheus.NewSummary(prometheus.SummaryOpts{
				Namespace:  ctx.Namespace,
				Subsystem:  subsystem,
				Name:       name + "_summary",
				Help:       help,
				Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.025, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001, 0.999: 0.0001},
			})
			prometheus.MustRegister(sum)
			ctx.Summaries[key] = sum

			hist = prometheus.NewHistogram(prometheus.HistogramOpts{
				Namespace: ctx.Namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
				Buckets:   buckets,
			})
			prometheus.MustRegister(hist)
			ctx.Histograms[key] = hist
		}
		ctx.histMutex.Unlock()
	}

	mh := MetricsHistogram{
		Key:  key,
		hist: hist,
		sum:  sum,
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
