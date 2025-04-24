package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var ClickCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "url_redirect_total",
		Help: "Total number of URL redirects by short code",
	},
	[]string{"short_code"},
)
