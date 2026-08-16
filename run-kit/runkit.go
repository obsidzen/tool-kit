// Package runkit provides injectable process runners for obsidzen tools.
package runkit

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

type CommandSpec struct {
	Dir  string
	Name string
	Args []string
	Env  []string
}

type Task struct {
	Key         string
	Description string
	Detail      string
	Spec        CommandSpec
	Specs       []CommandSpec
	Stream      func(context.Context) (<-chan Line, error)
}

type Runner interface {
	Run(context.Context, CommandSpec) ([]byte, error)
	Stream(context.Context, CommandSpec) (io.ReadCloser, func() error, error)
}

type Line struct {
	Text string
	Err  error
	Done bool
}

type CommandLineFormatter func(CommandSpec) string

type ExecRunner struct{}

func MergeEnv(overrides map[string]string) []string {
	env := os.Environ()
	for key, value := range overrides {
		env = append(env, key+"="+value)
	}
	return env
}

func StreamTo(ctx context.Context, runner Runner, spec CommandSpec, writer io.Writer) error {
	stdout, wait, err := runner.Stream(ctx, spec)
	if err != nil {
		return err
	}
	defer stdout.Close()
	if _, err := io.Copy(writer, stdout); err != nil {
		return err
	}
	return wait()
}

func RunTask(ctx context.Context, runner Runner, task Task) error {
	if task.Stream != nil {
		lines, err := task.Stream(ctx)
		if err != nil {
			return err
		}
		for line := range lines {
			if line.Err != nil {
				return line.Err
			}
		}
		return nil
	}
	for _, spec := range TaskSpecs(task) {
		if _, err := runner.Run(ctx, spec); err != nil {
			return err
		}
	}
	return nil
}

func StreamTaskTo(ctx context.Context, runner Runner, task Task, writer io.Writer) error {
	if task.Stream != nil {
		lines, err := task.Stream(ctx)
		if err != nil {
			return err
		}
		for line := range lines {
			if line.Err != nil {
				return line.Err
			}
			if line.Done {
				return nil
			}
			if _, err := io.WriteString(writer, line.Text+"\n"); err != nil {
				return err
			}
		}
		return nil
	}
	for _, spec := range TaskSpecs(task) {
		// TUI 경로(StreamTaskLinesWithFormatter)와 같은 머리글을 단다. 여러 단계를
		// 잇는 작업에서 이것이 없으면 출력이 한 덩어리로 쏟아져, 실패했을 때 어느
		// 단계였는지 로그만 보고는 알 수 없다.
		if _, err := io.WriteString(writer, "$ "+CommandLine(spec)+"\n"); err != nil {
			return err
		}
		if err := StreamTo(ctx, runner, spec, writer); err != nil {
			return err
		}
	}
	return nil
}

func StreamLines(ctx context.Context, runner Runner, spec CommandSpec) (<-chan Line, error) {
	stdout, wait, err := runner.Stream(ctx, spec)
	if err != nil {
		return nil, err
	}
	ch := make(chan Line)
	go func() {
		defer close(ch)
		defer stdout.Close()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			ch <- Line{Text: scanner.Text()}
		}
		if err := scanner.Err(); err != nil {
			ch <- Line{Err: err, Done: true}
			_ = wait()
			return
		}
		if err := wait(); err != nil {
			ch <- Line{Err: err, Done: true}
			return
		}
		ch <- Line{Done: true}
	}()
	return ch, nil
}

func StreamTaskLines(ctx context.Context, runner Runner, task Task) (<-chan Line, error) {
	return StreamTaskLinesWithFormatter(ctx, runner, task, CommandLine)
}

func StreamTaskLinesWithFormatter(ctx context.Context, runner Runner, task Task, formatter CommandLineFormatter) (<-chan Line, error) {
	if task.Stream != nil {
		return task.Stream(ctx)
	}
	if formatter == nil {
		formatter = CommandLine
	}
	out := make(chan Line)
	go func() {
		defer close(out)
		for _, spec := range TaskSpecs(task) {
			out <- Line{Text: "$ " + formatter(spec)}
			lines, err := StreamLines(ctx, runner, spec)
			if err != nil {
				out <- Line{Err: err, Done: true}
				return
			}
			for line := range lines {
				if line.Done {
					if line.Err != nil {
						out <- line
						return
					}
					break
				}
				out <- line
			}
		}
		out <- Line{Done: true}
	}()
	return out, nil
}

func TaskSpecs(task Task) []CommandSpec {
	specs := task.Specs
	if task.Spec.Name != "" {
		specs = append([]CommandSpec{task.Spec}, specs...)
	}
	return specs
}

func TaskWithArgs(task Task, args []string) Task {
	if len(args) == 0 {
		return task
	}
	specs := TaskSpecs(task)
	if len(specs) == 0 {
		return task
	}
	specs = append([]CommandSpec(nil), specs...)
	last := len(specs) - 1
	specs[last].Args = append(append([]string(nil), specs[last].Args...), args...)
	task.Spec = CommandSpec{}
	task.Specs = specs
	return task
}

func CommandLine(spec CommandSpec) string {
	parts := append([]string{spec.Name}, spec.Args...)
	return strings.TrimSpace(strings.Join(parts, " "))
}

func RedactCommandLine(spec CommandSpec, sensitiveFlags ...string) string {
	redacted := spec
	redacted.Args = append([]string{}, spec.Args...)
	sensitive := map[string]struct{}{}
	for _, flag := range sensitiveFlags {
		sensitive[flag] = struct{}{}
	}
	for i, arg := range redacted.Args {
		if _, ok := sensitive[arg]; ok && i+1 < len(redacted.Args) {
			redacted.Args[i+1] = "********"
		}
	}
	return CommandLine(redacted)
}

func (ExecRunner) Run(ctx context.Context, spec CommandSpec) ([]byte, error) {
	cmd := exec.CommandContext(ctx, spec.Name, spec.Args...)
	cmd.Dir = spec.Dir
	if len(spec.Env) > 0 {
		cmd.Env = spec.Env
	}
	return cmd.CombinedOutput()
}

func (ExecRunner) Stream(ctx context.Context, spec CommandSpec) (io.ReadCloser, func() error, error) {
	cmd := exec.CommandContext(ctx, spec.Name, spec.Args...)
	cmd.Dir = spec.Dir
	if len(spec.Env) > 0 {
		cmd.Env = spec.Env
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	return stdout, cmd.Wait, nil
}
