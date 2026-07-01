package runkit

import (
	"context"
	"io"
	"strings"
	"testing"
)

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
