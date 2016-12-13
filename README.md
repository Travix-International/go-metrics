# go-metrics
Package to register metrics into Prometheus from golang

# Dependencies

* [logger](https://github.com/Travix-International/logger) 
* [Prometheus](https://github.com/prometheus/client_golang/prometheus)

Use [gvt](https://github.com/FiloSottile/gvt) to manage the dependencies.

# Usage

Note: please refer to the Prometheus documentation for the exact details on each counter and 
on how to choose the counters that best fit your use cases.

```golang
// Initialization
logger := // instantiate logger
metrics := metrics.NewMetrics("namespace", logger)

// Example: simple count
metrics.Count("searchcache", "hit", "SearchCache cache hits")

// Example: count with labels
metrics.CountLabels("service", "get_liveness_request_duration_milliseconds", "Response times for GET requests to liveness in milliseconds.",
    []string{"status", "method"}, []string{strconv.Itoa(ww.status), r.Method})

// Example: increase counter by 5
metrics.IncreaseCounter("kpimetrics", "fullflowpublishing_errors_count", "total number of errors while publishing fares to Datastore.", 5)

// Example: histogram of processing times
func UseHistogram(req *someRequest) someResponse {
		histogram := metrics.AddHistogram("datastore", "getroutefare_request_duration_milliseconds", "Response times for getRouteFare requests to Datastore in milliseconds.")
		start := time.Now()
		defer func() {
			histogram.RecordTimeElapsed(start)
		}()

		// Do work that runs for some time

		return myResponse
	}
```

# Known limitations

## Histograms

When adding a histogram, the library will create both a `summary` and `histogram`. The summary is currently hard-coded
to use the following percentiles: `0.5 / 0.75 / 0.9 / 0.95 / 0.99 / 0.999`. This should cover most use cases. As for the
histogram, this is currently using the default buckets from the Prometheus library, which will not fit all use cases. It will
work correctly for simple network traffic.