// Package metrics defines the shared metric contracts for the Algoryn ecosystem.
// All Algoryn tools (Pulse, Relay, Beacon) use these types to ensure
// a consistent representation of observability data across the stack.
package metrics

import "time"

// LatencyStats contains aggregate latency measurements for a set of requests.
type LatencyStats struct {
	Min  time.Duration
	Mean time.Duration
	P50  time.Duration
	P90  time.Duration
	P95  time.Duration
	P99  time.Duration
	Max  time.Duration
}

// MetricSnapshot represents aggregated metrics over a time window.
// It is the core unit of observability data shared across Algoryn tools.
// Pulse produces snapshots from load test runs.
// Relay produces snapshots from live traffic.
// Beacon consumes snapshots to evaluate alert conditions.
type MetricSnapshot struct {
	// Source identifies which tool produced this snapshot.
	// Known values: "pulse", "relay".
	Source string

	// Service is the name of the service being observed.
	Service string

	// Timestamp is when the observation window started.
	Timestamp time.Time

	// Window is the duration of the observation window.
	Window time.Duration

	// Total is the number of requests observed.
	Total int64

	// Failed is the number of requests that resulted in an error.
	Failed int64

	// RPS is the average requests per second over the window.
	RPS float64

	// Latency contains aggregate latency measurements.
	Latency LatencyStats

	// StatusCodes maps HTTP status codes to their occurrence count.
	StatusCodes map[int]int64

	// Errors maps normalized error categories to their occurrence count.
	// Known categories: "http_status_error", "deadline_exceeded",
	// "context_canceled", "unknown_error".
	Errors map[string]int64

	// Labels contains tool-specific or user-defined metadata.
	// Pulse uses: {"phase": "spike"}.
	// Relay uses: {"route": "/api/users", "method": "GET"}.
	Labels map[string]string
}

// ThresholdResult records whether a single threshold condition was met.
type ThresholdResult struct {
	// Description is a human-readable representation of the threshold.
	// Example: "error_rate < 0.01"
	Description string

	// Pass is true when the threshold condition was satisfied.
	Pass bool
}

// RunEvent represents the complete result of a single Pulse test run.
// It is emitted via ResultHook and can be forwarded to Relay or any
// other Algoryn tool for storage and visualization.
type RunEvent struct {
	// ID is a unique identifier for this run.
	ID string

	// Source identifies the tool that produced this event.
	// Always "pulse" for events produced by Pulse.
	Source string

	// StartedAt is when the run began.
	StartedAt time.Time

	// EndedAt is when the run completed.
	EndedAt time.Time

	// Snapshot contains the aggregated metrics for the entire run.
	Snapshot MetricSnapshot

	// Thresholds contains the evaluation result for each configured threshold.
	Thresholds []ThresholdResult

	// Passed is true when all configured thresholds were met.
	Passed bool
}

// Sources defines the known source identifiers for Algoryn tools.
const (
	SourcePulse  = "pulse"
	SourceRelay  = "relay"
	SourceBeacon = "beacon"
)
