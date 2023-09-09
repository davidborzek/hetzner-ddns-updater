package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// UpdateCounter is a prometheus counter to count the dyndns updates.
	UpdateCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hdu_total_updates",
			Help: "The total number of dyndns updates",
		},
	)

	// UpdateFailedCounter is a prometheus counter to count the failed dyndns updates.
	UpdateFailedCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hdu_failed_updates",
			Help: "The total number of failed dyndns updates",
		},
	)
)

func EnableMetrics() {
	prometheus.MustRegister(
		UpdateCounter,
		UpdateFailedCounter,
	)
}
