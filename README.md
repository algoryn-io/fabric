# Algoryn Fabric

Shared contracts and interfaces for the Algoryn ecosystem.

---

## What is Algoryn Fabric?

Algoryn Fabric is an open source infrastructure toolkit for Go teams
building reliable products without vendor lock-in or enterprise pricing.

Each tool in the ecosystem works independently. Together, they cover
the full lifecycle of a production service — from load testing to
traffic management, deployment, and alerting.

| Tool       | What it does                          | Status        |
|------------|---------------------------------------|---------------|
| **Pulse**  | Load testing & chaos engineering      | `v0.2.0`      |
| **Relay**  | API Gateway & observability           | `coming soon` |
| **Beacon** | Alerting & on-call                    | `planned`     |
| **Deploy** | Deployment without friction           | `planned`     |
| **Dev**    | Integrated local environment          | `planned`     |

---

## What is this repository?

This repository defines the shared metric and event contracts used
across all Algoryn tools.

If you are building an integration with any Algoryn tool, start here.
```
algoryn.io/fabric/metrics   → MetricSnapshot, RunEvent, LatencyStats
algoryn.io/fabric/events    → Event, EventType, inter-tool payloads
```

---

## Why shared contracts?

Each Algoryn tool is independent — you can use Pulse without Relay,
Relay without Beacon. But when you combine them, they speak the same
language natively.

Without shared contracts:
```
Pulse result  →  custom adapter  →  Relay input
Relay metrics →  custom adapter  →  Beacon input
```

With Algoryn Fabric:
```
Pulse result  →  fabric.RunEvent     →  Relay
Relay metrics →  fabric.MetricSnapshot →  Beacon
```

No adapters. No translation. The same types, the same field names,
the same semantics across the entire stack.

---

## Installation
```bash
go get algoryn.io/fabric
```

---

## Usage

### Consuming a Pulse RunEvent in Relay
```go
import (
    "algoryn.io/fabric/metrics"
    "algoryn.io/fabric/events"
)

// Relay receives this event from Pulse via ResultHook
func handleRunEvent(e events.Event) {
    if e.Type != events.EventTypeRunCompleted {
        return
    }
    // process e.Payload as events.RunCompletedPayload
}
```

### Producing a MetricSnapshot in Relay
```go
snapshot := metrics.MetricSnapshot{
    Source:  metrics.SourceRelay,
    Service: "api-gateway",
    Window:  time.Minute,
    Total:   1500,
    Failed:  3,
    RPS:     25.0,
    Labels:  map[string]string{"route": "/api/users"},
}
```

---

## Compatibility

Fabric follows semantic versioning strictly. While in `v0.x.x`,
minor versions may introduce breaking changes with a deprecation
notice in the CHANGELOG.

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full compatibility policy.

---

## Part of Algoryn Fabric

| Repository | Module |
|------------|--------|
| [algoryn-io/fabric](https://github.com/algoryn-io/fabric) | `algoryn.io/fabric` |
| [algoryn-io/pulse](https://github.com/algoryn-io/pulse) | `algoryn.io/pulse` |

---

## License

MIT