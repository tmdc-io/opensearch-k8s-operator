package helpers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNormalizeHTTPPath(t *testing.T) {
	cases := map[string]string{
		"":                                 "/",
		"/":                                "/",
		"/metrics":                         "/metrics",
		"/_cluster/health":                 "/_cluster/health",
		"/_plugins/_security/api/roles/abc-def-ghi-jkl-mnop": "/_plugins/_security/api/roles/:id",
		"/foo/123": "/foo/:id",
	}
	for in, want := range cases {
		if got := normalizeHTTPPath(in); got != want {
			t.Errorf("normalizeHTTPPath(%q)=%q want %q", in, got, want)
		}
	}
}

func TestInstrumentHTTPHandler(t *testing.T) {
	HTTPRequestsTotal.Reset()
	HTTPRequestDurationSeconds.Reset()
	HTTPResponseStatusTotal.Reset()

	h := InstrumentHTTPHandler(HTTPComponentOperator, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if testutil.ToFloat64(HTTPRequestsTotal.WithLabelValues(HTTPComponentOperator, "GET", "/metrics")) != 1 {
		t.Fatal("expected http_requests_total=1")
	}
	if testutil.ToFloat64(HTTPResponseStatusTotal.WithLabelValues(HTTPComponentOperator, "GET", "/metrics", "204")) != 1 {
		t.Fatal("expected http_response_status_total=1")
	}
}

func TestInstrumentHTTPRoundTripper(t *testing.T) {
	HTTPRequestsTotal.Reset()
	HTTPRequestDurationSeconds.Reset()
	HTTPResponseStatusTotal.Reset()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "{}")
	}))
	t.Cleanup(server.Close)

	client := &http.Client{Transport: InstrumentHTTPRoundTripper(HTTPComponentCluster, http.DefaultTransport)}
	resp, err := client.Get(server.URL + "/_cluster/health")
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()

	if testutil.ToFloat64(HTTPRequestsTotal.WithLabelValues(HTTPComponentCluster, "GET", "/_cluster/health")) != 1 {
		t.Fatal("expected cluster http_requests_total=1")
	}
	if testutil.CollectAndCount(HTTPRequestDurationSeconds) == 0 {
		t.Fatal("expected duration histogram samples")
	}
	if testutil.ToFloat64(HTTPResponseStatusTotal.WithLabelValues(HTTPComponentCluster, "GET", "/_cluster/health", "200")) != 1 {
		t.Fatal("expected cluster http_response_status_total=1")
	}
}
