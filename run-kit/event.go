package runkit

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const EventSchemaVersion = "2"

//go:embed schema/tool-event.v2.schema.json
var eventSchemaFS embed.FS

// EventStatus is a stable lifecycle state shared by CLI, TUI, and structured output.
type EventStatus string

// ErrorCode is a stable, renderer-independent failure identifier.
type ErrorCode string

// Attempt is a one-based execution attempt number.
type Attempt uint

// ProgressCount is a non-negative progress position or denominator.
type ProgressCount uint

// ProgressUnit names the concrete items counted by progress.
type ProgressUnit string

// Progress records a current position against an explicit denominator.
type Progress struct {
	Current ProgressCount `json:"current"`
	Total   ProgressCount `json:"total"`
	Unit    ProgressUnit  `json:"unit"`
}

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
	ErrorCode     ErrorCode   `json:"errorCode,omitempty"`
	Attempt       Attempt     `json:"attempt,omitempty"`
	Progress      *Progress   `json:"progress,omitempty"`
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
	if e.Status == StatusFailed && strings.TrimSpace(string(e.ErrorCode)) == "" {
		return fmt.Errorf("failed event error code is required")
	}
	if e.ErrorCode != "" && !isKebabCase(string(e.ErrorCode)) {
		return fmt.Errorf("event error code must be kebab-case")
	}
	if e.Status != StatusPlanned && e.Status != StatusSkipped && e.Status != StatusNotApplicable && e.Attempt == 0 {
		return fmt.Errorf("executed event attempt must start at one")
	}
	if (e.Status == StatusPlanned || e.Status == StatusSkipped || e.Status == StatusNotApplicable) && e.Attempt != 0 {
		return fmt.Errorf("unexecuted event cannot have an attempt")
	}
	if e.ElapsedMS < 0 {
		return fmt.Errorf("event elapsed time cannot be negative")
	}
	if e.Progress != nil {
		if e.Progress.Total == 0 {
			return fmt.Errorf("event progress total must be greater than zero")
		}
		if e.Progress.Current > e.Progress.Total {
			return fmt.Errorf("event progress current cannot exceed total")
		}
		if !isKebabCase(string(e.Progress.Unit)) {
			return fmt.Errorf("event progress unit must be kebab-case")
		}
	}
	return nil
}

// EventJSONSchema returns an independent copy of the versioned event JSON Schema.
func EventJSONSchema() []byte {
	schema, err := eventSchemaFS.ReadFile("schema/tool-event.v2.schema.json")
	if err != nil {
		panic(err)
	}
	return schema
}

func isKebabCase(value string) bool {
	if value == "" || value[0] == '-' || value[len(value)-1] == '-' {
		return false
	}
	previousHyphen := false
	for _, character := range value {
		if character == '-' {
			if previousHyphen {
				return false
			}
			previousHyphen = true
			continue
		}
		if character < 'a' || character > 'z' {
			return false
		}
		previousHyphen = false
	}
	return true
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
