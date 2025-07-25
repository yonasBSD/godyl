// Package download provides file‑download utilities with timeout,
// retry, and optional SSL‑verification control.
package download

import (
	"time"

	"github.com/hashicorp/go-getter/v2"
)

// Option defines a functional option for configuring a Downloader.
type Option func(*Downloader)

// WithProgress returns an option that sets the progress tracker.
func WithProgress(progressTracker getter.ProgressTracker) Option {
	return func(d *Downloader) {
		d.progressListener = progressTracker
	}
}

// WithContextTimeout returns an option that sets the context timeout.
func WithContextTimeout(timeout time.Duration) Option {
	return func(d *Downloader) {
		d.contextTimeout = timeout
	}
}

// WithReadTimeout returns an option that sets the read timeout.
func WithReadTimeout(timeout time.Duration) Option {
	return func(d *Downloader) {
		d.readTimeout = timeout
	}
}

// WithHeadTimeout returns an option that sets the head timeout.
func WithHeadTimeout(timeout time.Duration) Option {
	return func(d *Downloader) {
		d.headTimeout = timeout
	}
}

// WithInsecureSkipVerify returns an option that sets whether to skip TLS verification.
func WithInsecureSkipVerify() Option {
	return func(d *Downloader) {
		d.insecureSkipVerify = true
	}
}

// WithMaxRetries returns an option that sets the maximum number of retry attempts.
func WithMaxRetries(retries int) Option {
	return func(d *Downloader) {
		d.maxRetries = retries
	}
}

// WithRetryWaits returns an option that sets the min and max retry wait durations.
func WithRetryWaits(min, max time.Duration) Option {
	return func(d *Downloader) {
		d.retryWaitMin = min
		d.retryWaitMax = max
	}
}
