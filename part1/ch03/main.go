package main

import (
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	COUNTER = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hello_world_total",
		Help: "Hello World requested",
	})

	GAUGE = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hello_world_connection",
		Help: "Number of /gauge in progress",
	})

	SUMMARY = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "hello_world_latency_seconds",
		Help: "Latency Time for a request /summary",
	})

	HISTOGRAM = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "hello_world_latency_histogram",
		Help:    "A histogram of Latency Time for a request /histogram",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	})
)

func index(w http.ResponseWriter, r *http.Request) {
	COUNTER.Inc()
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func gauge(w http.ResponseWriter, r *http.Request) {
	GAUGE.Inc()
	defer GAUGE.Dec()
	time.Sleep(10 * time.Second)
	fmt.Fprintf(w, "Gauge, %q", html.EscapeString(r.URL.Path))
}

func summary(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer SUMMARY.Observe(float64(time.Now().Sub(start)))
	fmt.Fprintf(w, "Summary, %q", html.EscapeString(r.URL.Path))
}

func histogram(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer HISTOGRAM.Observe(float64(time.Now().Sub(start)))
	fmt.Fprintf(w, "Histogram, %q", html.EscapeString(r.URL.Path))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/gauge", gauge)
	http.HandleFunc("/summary", summary)
	http.HandleFunc("/histogram", histogram)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
