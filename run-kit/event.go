package runkit

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const EventSchemaVersion = "1"

//go:embed schema/tool-event.v1.schema.json
var eventSchemaFS embed.FS

// EventStatus is a stable lifecycle state shared by CLI, TUI, and structured output.
type EventStatus string

const (
	StatusPlanned       EventStatus = "planned"
	StatusRunning       EventStatus = "running"
	StatusPassed        EventStatus = "passed"
	StatusFailed        EventStatus = "failed"
	StatusSkipped       EventStatus = "skipped"
	StatusNotApplicable EventStatus = "not-applicable"
)

// Event is a renderer-independent task lifecycle or result record.
type Event struct {
	SchemaVersion string      `json:"schemaVersion"`
	EventID       string      `json:"eventId,omitempty"`
	Tool          string      `json:"tool,omitempty"`
	Command       string      `json:"command,omitempty"`
	PhaseID       string      `json:"phaseId"`
	Status        EventStatus `json:"status"`
	Message       string      `json:"message"`
	Detail        string      `json:"detail,omitempty"`
	ElapsedMS     int64       `json:"elapsedMs,omitempty"`
	ErrorCode     string      `json:"errorCode,omitempty"`
	Attempt       int         `json:"attempt,omitempty"`
	Current       int         `json:"current,omitempty"`
	Total         int         `json:"total,omitempty"`
}

// Validate checks the required event contract without rendering it.
func (e Event) Validate() error {
	if e.SchemaVersion != EventSchemaVersion {
		return fmt.Errorf("unsupported event schema version: %q", e.SchemaVersion)
	}
	if strings.TrimSpace(e.PhaseID) == "" {
		return fmt.Errorf("event phase ID is required")
	}
	if strings.TrimSpace(e.Message) == "" {
		return fmt.Errorf("event message is required")
	}
	switch e.Status {
	case StatusPlanned, StatusRunning, StatusPassed, StatusFailed, StatusSkipped, StatusNotApplicable:
	default:
		return fmt.Errorf("unsupported event status: %q", e.Status)
	}
	if e.Status == StatusFailed && strings.TrimSpace(e.ErrorCode) == "" {
		return fmt.Errorf("failed event error code is required")
	}
	if e.Attempt < 0 || e.Current < 0 || e.Total < 0 || e.ElapsedMS < 0 {
		return fmt.Errorf("event counters and elapsed time cannot be negative")
	}
	if e.Total > 0 && e.Current > e.Total {
		return fmt.Errorf("event progress current cannot exceed total")
	}
	return nil
}

// EventJSONSchema returns an independent copy of the versioned event JSON Schema.
func EventJSONSchema() []byte {
	schema, err := eventSchemaFS.ReadFile("schema/tool-event.v1.schema.json")
	if err != nil {
		panic(err)
	}
	return schema
}

// HumanLine renders the stable status, phase, and message fields as one line.
func (e Event) HumanLine() string {
	return fmt.Sprintf("%s %s — %s", e.Status, e.PhaseID, e.Message)
}

// WriteEventJSON validates and writes one newline-delimited JSON event.
func WriteEventJSON(w io.Writer, event Event) error {
	if err := event.Validate(); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(event)
}
