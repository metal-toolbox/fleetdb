package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	apiLatencySec *prometheus.HistogramVec
	dbErrorCount  *prometheus.CounterVec
)

const (
	appName = "fleetdb"
	apiSub  = "api"
	dbSub   = "database"
)

func init() {
	apiLatencySec = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: appName,
			Subsystem: apiSub,
			Name:      "latency_seconds",
			Help:      "api latency measurements in seconds",
			// XXX: These buckets will likely need some tuning
			Buckets: []float64{0.025, 0.05, 0.1, 0.25, 0.5, 0.75, 1.0, 2.5, 5.0, 7.5, 10.0, 15.0, 20.0},
		},
		[]string{
			"code",
			"endpoint",
			"method",
		},
	)

	dbErrorCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: appName,
			Subsystem: dbSub,
			Name:      "error_count",
			Help:      "total count of database errors",
		},
		[]string{
			"operation",
		},
	)
}

// APICallEpilog observes the response and elapsed time of a call to a given endpoint
func APICallEpilog(start time.Time, endpoint, method string, responseCode int) {
	code := strconv.Itoa(responseCode)
	elapsed := time.Since(start).Seconds()
	apiLatencySec.WithLabelValues(code, endpoint, method).Observe(elapsed)
}

// DBError observes errors arising from an attempt to read or write data to the remote database
func DBError(op string) {
	dbErrorCount.WithLabelValues(op).Inc()
}
