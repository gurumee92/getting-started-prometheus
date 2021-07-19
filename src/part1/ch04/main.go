package main

import (
	"fmt"
	"html"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	REQUEST = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"path", "method"},
	)
)

func init() {
	prometheus.MustRegister(REQUEST)
}

func index(w http.ResponseWriter, r *http.Request) {
	REQUEST.WithLabelValues(r.URL.Path, r.Method).Inc()
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
