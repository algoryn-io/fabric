# Changelog

All notable changes to Algoryn Fabric will be documented in this file.

---

## [v0.1.0] — 2026-03-25

Initial release of the shared contract definitions for the Algoryn ecosystem.

### Added

**metrics package**
- `MetricSnapshot` — aggregated metrics over a time window
- `LatencyStats` — latency percentiles (min, mean, p50, p95, p99, max)
- `ThresholdResult` — single threshold evaluation outcome
- `RunEvent` — complete result of a Pulse test run
- `SourcePulse`, `SourceRelay`, `SourceBeacon` — known source identifiers

**events package**
- `Event` — envelope for all inter-tool communication
- `EventType` — typed event identifiers
- `RunCompletedPayload` — emitted by Pulse on run completion
- `RunSummary` — compact metrics for alerting and dashboards
- `ServiceRegisteredPayload` — emitted by Deploy on backend registration
- `ThresholdViolatedPayload` — emitted when a threshold is violated
- Known event types: `run.completed`, `service.registered`,
  `service.deregistered`, `threshold.violated`, `alert.fired`