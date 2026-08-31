package helpers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	HTTPComponentOperator = "operator"
	HTTPComponentCluster  = "cluster"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"component", "method", "path"},
	)
	HTTPRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"component", "method", "path"},
	)
	HTTPResponseStatusTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_total",
			Help: "Total number of HTTP responses by status code.",
		},
		[]string{"component", "method", "path", "status"},
	)
)

func httpCollectors() []prometheus.Collector {
	return []prometheus.Collector{HTTPRequestsTotal, HTTPRequestDurationSeconds, HTTPResponseStatusTotal}
}

func observeHTTP(component, method, path string, status int, duration time.Duration) {
	method = strings.ToUpper(method)
	path = normalizeHTTPPath(path)
	statusLabel := strconv.Itoa(status)
	if status == 0 {
		statusLabel = "error"
	}
	HTTPRequestsTotal.WithLabelValues(component, method, path).Inc()
	HTTPRequestDurationSeconds.WithLabelValues(component, method, path).Observe(duration.Seconds())
	HTTPResponseStatusTotal.WithLabelValues(component, method, path, statusLabel).Inc()
}

func normalizeHTTPPath(path string) string {
	if path == "" {
		return "/"
	}
	if i := strings.IndexByte(path, '?'); i >= 0 {
		path = path[:i]
	}
	parts := strings.Split(path, "/")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		if isPathID(part) {
			out = append(out, ":id")
			continue
		}
		out = append(out, part)
	}
	if len(out) == 0 {
		return "/"
	}
	return "/" + strings.Join(out, "/")
}

func isPathID(part string) bool {
	if part == "" {
		return false
	}
	digits := 0
	for _, r := range part {
		if r >= '0' && r <= '9' {
			digits++
		}
	}
	if digits == len(part) {
		return true
	}
	return len(part) >= 8 && strings.Count(part, "-") >= 4
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func InstrumentHTTPHandler(component string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(rec, req)
		observeHTTP(component, req.Method, req.URL.Path, rec.status, time.Since(start))
	})
}

type instrumentedRoundTripper struct {
	component string
	next      http.RoundTripper
}

func InstrumentHTTPRoundTripper(component string, next http.RoundTripper) http.RoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &instrumentedRoundTripper{component: component, next: next}
}

func (t *instrumentedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := t.next.RoundTrip(req)
	status := 0
	if resp != nil {
		status = resp.StatusCode
	}
	observeHTTP(t.component, req.Method, req.URL.Path, status, time.Since(start))
	return resp, err
}
