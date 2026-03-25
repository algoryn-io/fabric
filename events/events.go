// Package events defines the inter-tool communication contracts
// for the Algoryn ecosystem. Tools emit events to notify other
// tools of significant state changes.
package events

import "time"

// EventType identifies the kind of event being emitted.
type EventType string

const (
	// EventTypeRunCompleted is emitted by Pulse when a load test finishes.
	EventTypeRunCompleted EventType = "run.completed"

	// EventTypeServiceRegistered is emitted by Deploy when a new
	// service backend becomes available. Relay listens for this.
	EventTypeServiceRegistered EventType = "service.registered"

	// EventTypeServiceDeregistered is emitted by Deploy when a
	// service backend is removed. Relay listens for this.
	EventTypeServiceDeregistered EventType = "service.deregistered"

	// EventTypeThresholdViolated is emitted by Pulse or Relay when
	// a metric exceeds a configured threshold. Beacon listens for this.
	EventTypeThresholdViolated EventType = "threshold.violated"

	// EventTypeAlertFired is emitted by Beacon when an alert condition
	// is triggered.
	EventTypeAlertFired EventType = "alert.fired"
)

// Event is the envelope for all inter-tool communication in Algoryn.
// Every tool that emits or consumes events uses this type.
type Event struct {
	// ID is a unique identifier for this event.
	ID string

	// Type identifies what happened.
	Type EventType

	// Source is the tool that emitted this event.
	// Use the constants from the metrics package: metrics.SourcePulse, etc.
	Source string

	// Timestamp is when the event occurred.
	Timestamp time.Time

	// Payload contains event-specific data serialized as JSON.
	// Use the typed payload structs below and serialize with encoding/json.
	Payload []byte
}

// RunCompletedPayload is the payload for EventTypeRunCompleted.
// Emitted by Pulse at the end of every test run.
type RunCompletedPayload struct {
	// RunID matches the RunEvent.ID from the metrics package.
	RunID string

	// Service is the name of the service that was tested.
	Service string

	// Passed indicates whether all thresholds were met.
	Passed bool

	// Duration is the total run time.
	Duration time.Duration

	// Summary contains key metrics for quick evaluation.
	// Full metrics are available via the RunEvent.
	Summary RunSummary
}

// RunSummary contains the key metrics from a Pulse run.
// It is intentionally compact — enough for alerting and dashboards
// without carrying the full MetricSnapshot.
type RunSummary struct {
	Total     int64
	Failed    int64
	RPS       float64
	ErrorRate float64
	P99Ms     float64 // P99 latency in milliseconds
}

// ServiceRegisteredPayload is the payload for EventTypeServiceRegistered.
// Emitted by Deploy when a new backend is ready to receive traffic.
type ServiceRegisteredPayload struct {
	// Name is the service identifier used across the ecosystem.
	Name string

	// Address is the network address of the backend (host:port).
	Address string

	// Tags contains arbitrary metadata about the service.
	// Example: {"version": "1.2.0", "region": "us-east-1"}
	Tags map[string]string
}

// ThresholdViolatedPayload is the payload for EventTypeThresholdViolated.
// Emitted by Pulse or Relay. Beacon uses this to trigger alerts.
type ThresholdViolatedPayload struct {
	// Service is the service where the violation occurred.
	Service string

	// Source is the tool that detected the violation.
	Source string

	// Description is the human-readable threshold description.
	// Example: "p99_latency < 200ms"
	Description string

	// Actual is the observed value that violated the threshold.
	// Stored as float64 for interoperability — interpret based on Description.
	Actual float64

	// Limit is the configured threshold value.
	Limit float64
}
