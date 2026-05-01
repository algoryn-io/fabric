// Package fabric provides helpers bridging legacy hand-written contracts (events, metrics)
// and the generated Protocol Buffer messages under gen/go/fabric/v1.
package fabric

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"algoryn.io/fabric/events"
	"algoryn.io/fabric/metrics"
	fabricv1 "algoryn.io/fabric/gen/go/fabric/v1"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var protoJSONUnmarshal = protojson.UnmarshalOptions{DiscardUnknown: true}
var protoJSONMarshal = protojson.MarshalOptions{
	UseProtoNames:   true,
	EmitUnpopulated: false,
}

// --- metrics: LatencyStats ---

func LatencyStatsToProto(m metrics.LatencyStats) *fabricv1.LatencyStats {
	return &fabricv1.LatencyStats{
		Min:  durationpb.New(m.Min),
		Mean: durationpb.New(m.Mean),
		P50:  durationpb.New(m.P50),
		P90:  durationpb.New(m.P90),
		P95:  durationpb.New(m.P95),
		P99:  durationpb.New(m.P99),
		Max:  durationpb.New(m.Max),
	}
}

func LatencyStatsFromProto(pb *fabricv1.LatencyStats) metrics.LatencyStats {
	if pb == nil {
		return metrics.LatencyStats{}
	}
	return metrics.LatencyStats{
		Min:  durationFromProto(pb.GetMin()),
		Mean: durationFromProto(pb.GetMean()),
		P50:  durationFromProto(pb.GetP50()),
		P90:  durationFromProto(pb.GetP90()),
		P95:  durationFromProto(pb.GetP95()),
		P99:  durationFromProto(pb.GetP99()),
		Max:  durationFromProto(pb.GetMax()),
	}
}

// --- metrics: MetricSnapshot ---

func MetricSnapshotToProto(m metrics.MetricSnapshot) *fabricv1.MetricSnapshot {
	status := make(map[int32]int64, len(m.StatusCodes))
	for code, n := range m.StatusCodes {
		status[int32(code)] = n
	}
	return &fabricv1.MetricSnapshot{
		Source:      m.Source,
		Service:     m.Service,
		Timestamp:   timestamppb.New(m.Timestamp),
		Window:      durationpb.New(m.Window),
		Total:       m.Total,
		Failed:      m.Failed,
		Rps:         m.RPS,
		Latency:     LatencyStatsToProto(m.Latency),
		StatusCodes: status,
		Errors:      cloneStringInt64Map(m.Errors),
		Labels:      cloneStringStringMap(m.Labels),
	}
}

func MetricSnapshotFromProto(pb *fabricv1.MetricSnapshot) metrics.MetricSnapshot {
	if pb == nil {
		return metrics.MetricSnapshot{}
	}
	sc := make(map[int]int64, len(pb.GetStatusCodes()))
	for code, n := range pb.GetStatusCodes() {
		sc[int(code)] = n
	}
	return metrics.MetricSnapshot{
		Source:      pb.GetSource(),
		Service:     pb.GetService(),
		Timestamp:   timestampFromProto(pb.GetTimestamp()),
		Window:      durationFromProto(pb.GetWindow()),
		Total:       pb.GetTotal(),
		Failed:      pb.GetFailed(),
		RPS:         pb.GetRps(),
		Latency:     LatencyStatsFromProto(pb.GetLatency()),
		StatusCodes: sc,
		Errors:      cloneStringInt64Map(pb.GetErrors()),
		Labels:      cloneStringStringMap(pb.GetLabels()),
	}
}

// --- metrics: ThresholdResult ---

func ThresholdResultToProto(m metrics.ThresholdResult) *fabricv1.ThresholdResult {
	return &fabricv1.ThresholdResult{
		Description: m.Description,
		Pass:        m.Pass,
	}
}

func ThresholdResultFromProto(pb *fabricv1.ThresholdResult) metrics.ThresholdResult {
	if pb == nil {
		return metrics.ThresholdResult{}
	}
	return metrics.ThresholdResult{
		Description: pb.GetDescription(),
		Pass:        pb.GetPass(),
	}
}

// --- metrics: RunEvent ---

func RunEventToProto(m metrics.RunEvent) *fabricv1.RunEvent {
	ths := make([]*fabricv1.ThresholdResult, 0, len(m.Thresholds))
	for i := range m.Thresholds {
		ths = append(ths, ThresholdResultToProto(m.Thresholds[i]))
	}
	return &fabricv1.RunEvent{
		Id:         m.ID,
		Source:     m.Source,
		StartedAt:  timestamppb.New(m.StartedAt),
		EndedAt:    timestamppb.New(m.EndedAt),
		Snapshot:   MetricSnapshotToProto(m.Snapshot),
		Thresholds: ths,
		Passed:     m.Passed,
	}
}

func RunEventFromProto(pb *fabricv1.RunEvent) metrics.RunEvent {
	if pb == nil {
		return metrics.RunEvent{}
	}
	ths := make([]metrics.ThresholdResult, 0, len(pb.GetThresholds()))
	for _, t := range pb.GetThresholds() {
		ths = append(ths, ThresholdResultFromProto(t))
	}
	return metrics.RunEvent{
		ID:         pb.GetId(),
		Source:     pb.GetSource(),
		StartedAt:  timestampFromProto(pb.GetStartedAt()),
		EndedAt:    timestampFromProto(pb.GetEndedAt()),
		Snapshot:   MetricSnapshotFromProto(pb.GetSnapshot()),
		Thresholds: ths,
		Passed:     pb.GetPassed(),
	}
}

// --- events: EventType ---

// EventTypeToProto maps legacy string-typed EventType constants to the protobuf enum.
func EventTypeToProto(t events.EventType) fabricv1.EventType {
	switch t {
	case events.EventTypeRunCompleted:
		return fabricv1.EventType_EVENT_TYPE_RUN_COMPLETED
	case events.EventTypeServiceRegistered:
		return fabricv1.EventType_EVENT_TYPE_SERVICE_REGISTERED
	case events.EventTypeServiceDeregistered:
		return fabricv1.EventType_EVENT_TYPE_SERVICE_DEREGISTERED
	case events.EventTypeThresholdViolated:
		return fabricv1.EventType_EVENT_TYPE_THRESHOLD_VIOLATED
	case events.EventTypeAlertFired:
		return fabricv1.EventType_EVENT_TYPE_ALERT_FIRED
	default:
		return fabricv1.EventType_EVENT_TYPE_UNSPECIFIED
	}
}

// EventTypeFromProto maps the protobuf enum to legacy events.EventType strings.
func EventTypeFromProto(t fabricv1.EventType) events.EventType {
	switch t {
	case fabricv1.EventType_EVENT_TYPE_RUN_COMPLETED:
		return events.EventTypeRunCompleted
	case fabricv1.EventType_EVENT_TYPE_SERVICE_REGISTERED:
		return events.EventTypeServiceRegistered
	case fabricv1.EventType_EVENT_TYPE_SERVICE_DEREGISTERED:
		return events.EventTypeServiceDeregistered
	case fabricv1.EventType_EVENT_TYPE_THRESHOLD_VIOLATED:
		return events.EventTypeThresholdViolated
	case fabricv1.EventType_EVENT_TYPE_ALERT_FIRED:
		return events.EventTypeAlertFired
	default:
		return events.EventType("")
	}
}

// --- events: typed payloads ---

func RunCompletedPayloadToProto(p *events.RunCompletedPayload) *fabricv1.RunCompletedPayload {
	if p == nil {
		return nil
	}
	return &fabricv1.RunCompletedPayload{
		RunId:    p.RunID,
		Service:  p.Service,
		Passed:   p.Passed,
		Duration: durationpb.New(p.Duration),
		Summary:  runSummaryToProto(&p.Summary),
	}
}

func RunCompletedPayloadFromProto(pb *fabricv1.RunCompletedPayload) events.RunCompletedPayload {
	if pb == nil {
		return events.RunCompletedPayload{}
	}
	return events.RunCompletedPayload{
		RunID:    pb.GetRunId(),
		Service:  pb.GetService(),
		Passed:   pb.GetPassed(),
		Duration: durationFromProto(pb.GetDuration()),
		Summary:  runSummaryFromProto(pb.GetSummary()),
	}
}

func runSummaryToProto(s *events.RunSummary) *fabricv1.RunSummary {
	if s == nil {
		return nil
	}
	return &fabricv1.RunSummary{
		Total:     s.Total,
		Failed:    s.Failed,
		Rps:       s.RPS,
		ErrorRate: s.ErrorRate,
		P99Ms:     s.P99Ms,
	}
}

func runSummaryFromProto(pb *fabricv1.RunSummary) events.RunSummary {
	if pb == nil {
		return events.RunSummary{}
	}
	return events.RunSummary{
		Total:     pb.GetTotal(),
		Failed:    pb.GetFailed(),
		RPS:       pb.GetRps(),
		ErrorRate: pb.GetErrorRate(),
		P99Ms:     pb.GetP99Ms(),
	}
}

func ServiceRegisteredPayloadToProto(p *events.ServiceRegisteredPayload) *fabricv1.ServiceRegisteredPayload {
	if p == nil {
		return nil
	}
	return &fabricv1.ServiceRegisteredPayload{
		Name:    p.Name,
		Address: p.Address,
		Tags:    cloneStringStringMap(p.Tags),
	}
}

func ServiceRegisteredPayloadFromProto(pb *fabricv1.ServiceRegisteredPayload) events.ServiceRegisteredPayload {
	if pb == nil {
		return events.ServiceRegisteredPayload{}
	}
	return events.ServiceRegisteredPayload{
		Name:    pb.GetName(),
		Address: pb.GetAddress(),
		Tags:    cloneStringStringMap(pb.GetTags()),
	}
}

func ThresholdViolatedPayloadToProto(p *events.ThresholdViolatedPayload) *fabricv1.ThresholdViolatedPayload {
	if p == nil {
		return nil
	}
	return &fabricv1.ThresholdViolatedPayload{
		Service:     p.Service,
		Source:      p.Source,
		Description: p.Description,
		Actual:      p.Actual,
		Limit:       p.Limit,
	}
}

func ThresholdViolatedPayloadFromProto(pb *fabricv1.ThresholdViolatedPayload) events.ThresholdViolatedPayload {
	if pb == nil {
		return events.ThresholdViolatedPayload{}
	}
	return events.ThresholdViolatedPayload{
		Service:     pb.GetService(),
		Source:      pb.GetSource(),
		Description: pb.GetDescription(),
		Actual:      pb.GetActual(),
		Limit:       pb.GetLimit(),
	}
}

// AlertFiredPayloadClone returns a deep copy of the protobuf alert payload (no legacy struct yet).
func AlertFiredPayloadClone(pb *fabricv1.AlertFiredPayload) *fabricv1.AlertFiredPayload {
	if pb == nil {
		return nil
	}
	return proto.Clone(pb).(*fabricv1.AlertFiredPayload)
}

// --- events: Event envelope ---

// EventToProto decodes JSON payloads from the legacy Event shape and sets the protobuf oneof.
func EventToProto(e events.Event) (*fabricv1.Event, error) {
	out := &fabricv1.Event{
		Id:        e.ID,
		Type:      EventTypeToProto(e.Type),
		Source:    e.Source,
		Timestamp: timestamppb.New(e.Timestamp),
	}

	switch e.Type {
	case events.EventTypeRunCompleted:
		var p events.RunCompletedPayload
		if err := json.Unmarshal(e.Payload, &p); err != nil {
			return nil, fmt.Errorf("run.completed payload: %w", err)
		}
		out.Payload = &fabricv1.Event_RunCompleted{RunCompleted: RunCompletedPayloadToProto(&p)}

	case events.EventTypeServiceRegistered:
		var p events.ServiceRegisteredPayload
		if err := json.Unmarshal(e.Payload, &p); err != nil {
			return nil, fmt.Errorf("service.registered payload: %w", err)
		}
		out.Payload = &fabricv1.Event_ServiceRegistered{ServiceRegistered: ServiceRegisteredPayloadToProto(&p)}

	case events.EventTypeServiceDeregistered:
		p := &fabricv1.ServiceDeregisteredPayload{}
		if len(e.Payload) > 0 {
			if err := protoJSONUnmarshal.Unmarshal(e.Payload, p); err != nil {
				var w serviceDeregisteredWire
				if err2 := json.Unmarshal(e.Payload, &w); err2 != nil {
					return nil, fmt.Errorf("service.deregistered payload: protojson: %v; json: %w", err, err2)
				}
				p.Name = w.Name
				p.Address = w.Address
				p.Tags = cloneStringStringMap(w.Tags)
			}
		}
		out.Payload = &fabricv1.Event_ServiceDeregistered{ServiceDeregistered: p}

	case events.EventTypeThresholdViolated:
		var p events.ThresholdViolatedPayload
		if err := json.Unmarshal(e.Payload, &p); err != nil {
			return nil, fmt.Errorf("threshold.violated payload: %w", err)
		}
		out.Payload = &fabricv1.Event_ThresholdViolated{ThresholdViolated: ThresholdViolatedPayloadToProto(&p)}

	case events.EventTypeAlertFired:
		p := &fabricv1.AlertFiredPayload{}
		if len(e.Payload) > 0 {
			if err := protoJSONUnmarshal.Unmarshal(e.Payload, p); err != nil {
				return nil, fmt.Errorf("alert.fired payload: %w", err)
			}
		}
		out.Payload = &fabricv1.Event_AlertFired{AlertFired: p}

	default:
		return nil, fmt.Errorf("unsupported events.EventType %q", e.Type)
	}

	return out, nil
}

// EventFromProto builds the legacy Event with JSON-encoded Payload bytes.
func EventFromProto(pb *fabricv1.Event) (events.Event, error) {
	if pb == nil {
		return events.Event{}, errors.New("nil Event")
	}
	out := events.Event{
		ID:        pb.GetId(),
		Type:      EventTypeFromProto(pb.GetType()),
		Source:    pb.GetSource(),
		Timestamp: timestampFromProto(pb.GetTimestamp()),
	}

	var (
		payload []byte
		err     error
	)

	switch p := pb.GetPayload().(type) {
	case *fabricv1.Event_RunCompleted:
		native := RunCompletedPayloadFromProto(p.RunCompleted)
		payload, err = json.Marshal(native)
	case *fabricv1.Event_ServiceRegistered:
		native := ServiceRegisteredPayloadFromProto(p.ServiceRegistered)
		payload, err = json.Marshal(native)
	case *fabricv1.Event_ServiceDeregistered:
		payload, err = protoJSONMarshal.Marshal(p.ServiceDeregistered)
	case *fabricv1.Event_ThresholdViolated:
		native := ThresholdViolatedPayloadFromProto(p.ThresholdViolated)
		payload, err = json.Marshal(native)
	case *fabricv1.Event_AlertFired:
		payload, err = protoJSONMarshal.Marshal(p.AlertFired)
	case nil:
		payload = nil
	default:
		return events.Event{}, fmt.Errorf("unknown protobuf payload type %T", p)
	}
	if err != nil {
		return events.Event{}, err
	}
	out.Payload = payload
	return out, nil
}

// --- helpers ---

func durationFromProto(d *durationpb.Duration) time.Duration {
	if d == nil {
		return 0
	}
	return d.AsDuration()
}

func timestampFromProto(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func cloneStringStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func cloneStringInt64Map(in map[string]int64) map[string]int64 {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]int64, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

type serviceDeregisteredWire struct {
	Name    string            `json:"name"`
	Address string            `json:"address"`
	Tags    map[string]string `json:"tags,omitempty"`
}
