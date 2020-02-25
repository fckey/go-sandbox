package retry

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(t *testing.T, fn RoundTripFunc) *http.Client {
	t.Helper()
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestRunWithHTTPRequestRetry_SuccessAtFive(t *testing.T) {
	ctx := context.Background()
	n := 0

	client := newTestClient(t, func(req *http.Request) *http.Response {
		n++
		var statusCode int
		if n < 5 {
			statusCode = http.StatusBadGateway
		} else {
			statusCode = http.StatusOK
		}
		return &http.Response{
			StatusCode: statusCode,
			Body:       ioutil.NopCloser(bytes.NewBuffer(nil)),
			Header:     make(http.Header),
		}
	})

	backoff := NewConfig(1*time.Second, 32*time.Second, 2, 10)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8080", bytes.NewBuffer(nil))

	resp, err := DoWithRetry(ctx, backoff, client, req, http.StatusBadGateway)
	if err != nil {
		t.Error("unexpected error")
	}
	if n != 5 {
		t.Errorf("n: got %d, want %d", n, 5)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("unexpected response")
	}
}

func TestRunWithHTTPRequestRetry_FailWithMaxRetry(t *testing.T) {
	ctx := context.Background()
	n := 0

	client := newTestClient(t, func(req *http.Request) *http.Response {
		n++
		var statusCode int
		if n < 5 {
			statusCode = http.StatusBadGateway
		} else {
			statusCode = http.StatusOK
		}
		return &http.Response{
			StatusCode: statusCode,
			Body:       ioutil.NopCloser(bytes.NewBuffer(nil)),
			Header:     make(http.Header),
		}
	})

	backoff := NewConfig(100*time.Microsecond, 10*time.Second, 2, 3)
	req, _ := http.NewRequest(http.MethodGet, "localhost:8080", bytes.NewBuffer(nil))

	resp, err := DoWithRetry(ctx, backoff, client, req, http.StatusBadGateway)
	if strings.Contains(err.Error(), "traial") {
		t.Error("unexpected error")
	}
	if n != 3 {
		t.Errorf("n: got %d, want %d", n, 5)
	}
	if resp.StatusCode != http.StatusBadGateway {
		t.Error("unexpected response")
	}
}
