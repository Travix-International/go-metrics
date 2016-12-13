package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/Travix-International/logger"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	Metrics struct {
		Namespace     string
		Counters      map[string]prometheus.Counter
		CounterVecs   map[string]*prometheus.CounterVec
		Summaries     map[string]prometheus.Summary
		Histograms    map[string]prometheus.Histogram
		Loggy         *logger.Logger
		countMutex    *sync.Mutex
		countVecMutex *sync.Mutex
		histMutex     *sync.Mutex
	}

	MetricsHistogram struct {
		Key  string
		hist prometheus.Histogram
		sum  prometheus.Summary
	}
)

func NewMetrics(namespace string, loggy *logger.Logger) *Metrics {
	m := Metrics{
		Namespace:     namespace,
		Loggy:         loggy,
		Counters:      make(map[string]prometheus.Counter),
		CounterVecs:   make(map[string]*prometheus.CounterVec),
		Histograms:    make(map[string]prometheus.Histogram),
		Summaries:     make(map[string]prometheus.Summary),
		countMutex:    &sync.Mutex{},
		countVecMutex: &sync.Mutex{},
		histMutex:     &sync.Mutex{},
	}
	return &m
}

func (ctx *Metrics) Count(subsystem, name, help string) {
	ctx.countMutex.Lock()
	defer ctx.countMutex.Unlock()

	key := fmt.Sprintf("%s/%s", subsystem, name)
	counter, exists := ctx.Counters[key]

	if !exists {
		counter = prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: ctx.Namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		})
		ctx.Counters[key] = counter
		err := prometheus.Register(counter)
		if err != nil {
			ctx.Loggy.Warn("MetricsCounterRegistrationFailed", fmt.Sprintf("CounterHandler: Counter registration failed", counter, err))
		}
	}

	counter.Inc()
}

func (ctx *Metrics) CountLabels(subsystem, name, help string, labels, values []string) {
	ctx.countVecMutex.Lock()
	defer ctx.countVecMutex.Unlock()

	key := fmt.Sprintf("%s/%s", subsystem, name)
	counter, exists := ctx.CounterVecs[key]

	if !exists {
		counter = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: ctx.Namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		}, labels)
		ctx.CounterVecs[key] = counter
		err := prometheus.Register(counter)
		if err != nil {
			ctx.Loggy.Warn("MetricsCounterLabelRegistrationFailed", fmt.Sprintf("CounterLabelHandler: Counter registration failed", counter, err))
		}
	}

	counter.WithLabelValues(values...).Inc()
}

func (ctx *Metrics) IncreaseCounter(subsystem, name, help string, increment int) {
	ctx.countMutex.Lock()
	defer ctx.countMutex.Unlock()

	key := fmt.Sprintf("%s/%s", subsystem, name)
	counter, exists := ctx.Counters[key]

	if !exists {
		counter = prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: ctx.Namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		})
		ctx.Counters[key] = counter
		err := prometheus.Register(counter)
		if err != nil {
			ctx.Loggy.Warn("MetricsIncreaseCounterRegistrationFailed", fmt.Sprintf("CounterHandler: Counter registration failed: %v: %v", counter, err))
		}
	}

	counter.Add(float64(increment))
}

func (ctx *Metrics) AddHistogram(subsystem, name, help string) *MetricsHistogram {
	ctx.histMutex.Lock()
	defer ctx.histMutex.Unlock()

	key := fmt.Sprintf("%s/%s", subsystem, name)

	sum, exists := ctx.Summaries[key]
	if !exists {
		sum = prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  ctx.Namespace,
			Subsystem:  subsystem,
			Name:       name + "_summary",
			Help:       help,
			Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.025, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001, 0.999: 0.0001},
		})
		prometheus.MustRegister(sum)
		ctx.Summaries[key] = sum
	}

	hist, exists := ctx.Histograms[key]
	if !exists {
		hist = prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: ctx.Namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		})
		prometheus.MustRegister(hist)
		ctx.Histograms[key] = hist
	}

	mh := MetricsHistogram{
		Key:  key,
		hist: hist,
		sum:  sum,
	}
	return &mh
}

func (histogram *MetricsHistogram) RecordTimeElapsed(start time.Time) {
	elapsed := float64(time.Since(start).Seconds())
	histogram.hist.Observe(elapsed)         // The default histogram buckets are recorded in seconds
	histogram.sum.Observe(elapsed * 1000.0) // While we have summaries in milliseconds
}
