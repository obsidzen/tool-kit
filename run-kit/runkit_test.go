package runkit

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
)

func TestEventContractAndRenderers(t *testing.T) {
	event := Event{
		SchemaVersion: EventSchemaVersion,
		EventID:       "database-check-planned",
		Tool:          "example-check",
		Command:       "database",
		PhaseID:       "database",
		Status:        StatusPlanned,
		Message:       "Check database schema",
		Current:       1,
		Total:         3,
	}
	if err := event.Validate(); err != nil {
		t.Fatal(err)
	}
	const want = "planned database — Check database schema"
	if got := EventLine(event).DisplayText(); got != want {
		t.Fatalf("display = %q, want %q", got, want)
	}
	var out bytes.Buffer
	if err := WriteEventJSON(&out, event); err != nil {
		t.Fatal(err)
	}
	var decoded Event
	if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Status != StatusPlanned || decoded.PhaseID != "database" || decoded.Total != 3 {
		t.Fatalf("decoded event = %#v", decoded)
	}
}

func TestEventRequiresCurrentSchemaVersion(t *testing.T) {
	for _, version := range []string{"", "2"} {
		event := Event{SchemaVersion: version, PhaseID: "database", Status: StatusPassed, Message: "Database check passed"}
		if err := event.Validate(); err == nil {
			t.Fatalf("schema version %q was accepted", version)
		}
	}
}

func TestEventJSONSchemaIsEmbeddedAndIndependent(t *testing.T) {
	first := EventJSONSchema()
	var schema map[string]any
	if err := json.Unmarshal(first, &schema); err != nil {
		t.Fatal(err)
	}
	if schema["title"] != "Tool Event v1" {
		t.Fatalf("schema title = %#v", schema["title"])
	}
	first[0] = 'x'
	if !json.Valid(EventJSONSchema()) {
		t.Fatal("caller mutation changed the embedded schema")
	}
}

func TestFailedEventRequiresErrorCode(t *testing.T) {
	event := Event{SchemaVersion: EventSchemaVersion, PhaseID: "database", Status: StatusFailed, Message: "Database check failed"}
	if err := event.Validate(); err == nil {
		t.Fatal("failed event without error code was accepted")
	}
}

func TestDiagnosticLinePreservesSourceAndPhase(t *testing.T) {
	line := DiagnosticLine("postgres", "database", "schema mismatch")
	if got, want := line.DisplayText(), "postgres/database — schema mismatch"; got != want {
		t.Fatalf("display = %q, want %q", got, want)
	}
}

type fakeRunner struct {
	output string
	err    error
}

func (r fakeRunner) Run(context.Context, CommandSpec) ([]byte, error) {
	return []byte(r.output), r.err
}

func (r fakeRunner) Stream(context.Context, CommandSpec) (io.ReadCloser, func() error, error) {
	return io.NopCloser(strings.NewReader(r.output)), func() error { return r.err }, nil
}

func TestStreamLines(t *testing.T) {
	lines, err := StreamLines(context.Background(), fakeRunner{output: "one\ntwo\n"}, CommandSpec{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for line := range lines {
		if line.Err != nil {
			t.Fatal(line.Err)
		}
		if line.Done {
			break
		}
		got = append(got, line.Text)
	}
	if strings.Join(got, ",") != "one,two" {
		t.Fatalf("lines = %#v", got)
	}
}

func TestTaskWithArgsAppendsToLastSpec(t *testing.T) {
	task := Task{
		Specs: []CommandSpec{
			{Name: "first", Args: []string{"a"}},
			{Name: "second", Args: []string{"b"}},
		},
	}
	got := TaskWithArgs(task, []string{"c", "d"})
	specs := TaskSpecs(got)
	if len(specs) != 2 {
		t.Fatalf("specs = %#v", specs)
	}
	if strings.Join(specs[1].Args, ",") != "b,c,d" {
		t.Fatalf("args = %#v", specs[1].Args)
	}
}

func TestStreamTaskLinesWithFormatter(t *testing.T) {
	task := Task{Spec: CommandSpec{Name: "echo", Args: []string{"secret"}}}
	lines, err := StreamTaskLinesWithFormatter(context.Background(), fakeRunner{output: "ok\n"}, task, func(CommandSpec) string {
		return "echo ********"
	})
	if err != nil {
		t.Fatal(err)
	}
	first := <-lines
	if first.Text != "$ echo ********" {
		t.Fatalf("first line = %q", first.Text)
	}
}

func TestCommandLine(t *testing.T) {
	if got := CommandLine(CommandSpec{Name: "node", Args: []string{"script.mjs", "--env", "dev"}}); got != "node script.mjs --env dev" {
		t.Fatalf("command line = %q", got)
	}
}

func TestRedactCommandLine(t *testing.T) {
	got := RedactCommandLine(CommandSpec{Name: "node", Args: []string{"script.mjs", "--password", "secret", "--env", "dev"}}, "--password")
	if strings.Contains(got, "secret") {
		t.Fatalf("secret leaked: %q", got)
	}
	if got != "node script.mjs --password ******** --env dev" {
		t.Fatalf("redacted line = %q", got)
	}
}

func TestStreamTaskToLabelsEachStep(t *testing.T) {
	task := Task{Key: "two-steps", Specs: []CommandSpec{
		{Name: "echo", Args: []string{"first"}},
		{Name: "echo", Args: []string{"second"}},
	}}
	var out strings.Builder

	if err := StreamTaskTo(context.Background(), ExecRunner{}, task, &out); err != nil {
		t.Fatalf("StreamTaskTo error: %v", err)
	}

	// Without the headers a failing multi-step run gives no clue which step broke.
	for _, want := range []string{"$ echo first", "first", "$ echo second", "second"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("output missing %q:\n%s", want, out.String())
		}
	}
}
