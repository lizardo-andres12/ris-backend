package client

import (
	"net/http"
	"time"
)

// NewHTTPClientTransport
func NewHTTPClientTransport() *http.Transport {
	// TODO: replace with config
	return &http.Transport{
		MaxIdleConns: 100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout: 90 * time.Second,
		DisableCompression: false,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// NewHTTPClient returns a new http.Client object with the specified configuration
func NewHTTPClient(transport *http.Transport) *http.Client {
	// TODO: replace with config
	timeout := 5

	return &http.Client{
		Transport: transport,
		Timeout: time.Duration(timeout) * time.Second,
	}
}

