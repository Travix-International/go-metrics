# Metrics [![Build Status](https://travis-ci.org/Travix-International/go-metrics.svg?branch=master)](https://travis-ci.org/Travix-International/go-metrics?branch=master)

[![Go Report Card](https://goreportcard.com/badge/github.com/Travix-International/go-metrics)](https://goreportcard.com/report/github.com/Travix-International/go-metrics) [![Coverage Status](https://coveralls.io/repos/github/Travix-International/go-metrics/badge.svg?branch=master)](https://coveralls.io/github/Travix-International/go-metrics?branch=master) 
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/Travix-International/go-metrics/blob/master/LICENSE)


> Package to register metrics into Prometheus from Go

## Package usage

Include this package into your project with:

```
go get github.com/Travix-International/go-metrics
```

# Usage

Note: please refer to the Prometheus documentation for the exact details on each counter and 
on how to choose the counters that best fit your use cases.

```golang
// Initialization
metrics := metrics.NewMetrics("namespace", logger.New())

// Example: simple count
metrics.Count("searchcache", "hit", "SearchCache cache hits")

// Example: count with labels
metrics.CountLabels("service", "get_liveness_request_duration_milliseconds", 
	"Response times for GET requests to liveness in milliseconds.",
	[]string{"status", "method"}, []string{strconv.Itoa(ww.status), r.Method})

// Example: increase counter by 5
metrics.IncreaseCounter("kpimetrics", "fullflowpublishing_errors_count", 
	"total number of errors while publishing fares to Datastore.", 5)

// Example: histogram of processing times
func UseHistogram(req *someRequest) someResponse {
		histogram := metrics.AddHistogram("datastore", "getroutefare_request_duration_milliseconds", 
			"Response times for getRouteFare requests to Datastore in milliseconds.")
		
		start := time.Now()
		defer func() {
			histogram.RecordTimeElapsed(start)
		}()

		// Do work that runs for some time

		return myResponse
	}
```

# Dependencies

* [logger](https://github.com/Travix-International/logger) 
* [Prometheus](https://github.com/prometheus/client_golang/prometheus)

# Known limitations

## Histograms

When adding a histogram, the library will create both a `summary` and `histogram`. The summary is currently hard-coded
to use the following percentiles: `0.5 / 0.75 / 0.9 / 0.95 / 0.99 / 0.999`. This should cover most use cases. As for the
histogram, this is currently using the default buckets from the Prometheus library, which will not fit all use cases. It will
work correctly for simple network traffic, but not for cases which typically have a higher value (more than one second).
The buckets are currently not configurable.


[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/Travix-International/go-metrics/blob/master/LICENSE)
