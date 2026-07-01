package clikit

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func TestNewRootAddsMenuActionAndCompletion(t *testing.T) {
	root := NewRoot(RootOptions{
		Use:   "tool",
		Short: "test tool",
		Menu:  func(context.Context) error { return nil },
		Actions: []Action{
			{Key: "status", Description: "show status", Aliases: []string{"st"}, RunE: func(context.Context) error { return nil }},
		},
	})

	names := commandNames(root)
	want := []string{"completion", "menu", "status"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("commands mismatch\nwant: %#v\n got: %#v", want, names)
	}
	status, _, err := root.Find([]string{"st"})
	if err != nil {
		t.Fatal(err)
	}
	if status.Name() != "status" {
		t.Fatalf("alias resolved to %s", status.Name())
	}
}

func TestNewRootAddsVersionCommand(t *testing.T) {
	root := NewRoot(RootOptions{
		Use:     "tool",
		Short:   "test tool",
		Version: VersionInfo{Name: "tool", Version: "1.2.3", Commit: "abc", Date: "today"},
	})
	root.SetArgs([]string{"version"})
	var out strings.Builder
	root.SetOut(&out)
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(out.String()); got != "tool 1.2.3 (commit abc, built today)" {
		t.Fatalf("version output = %q", got)
	}
}

func TestActionRunArgsE(t *testing.T) {
	var got []string
	root := NewRoot(RootOptions{
		Use: "tool",
		Actions: []Action{
			{
				Key:                "deploy",
				Description:        "deploy",
				Args:               ArbitraryArgs,
				DisableFlagParsing: true,
				RunArgsE: func(_ context.Context, args []string) error {
					got = args
					return nil
				},
			},
		},
	})
	root.SetArgs([]string{"deploy", "--server-only"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []string{"--server-only"}) {
		t.Fatalf("args = %#v", got)
	}
}

func TestActionDetailBecomesCommandLongHelp(t *testing.T) {
	root := NewRoot(RootOptions{
		Use: "tool",
		Actions: []Action{
			{Key: "status", Description: "show status", Detail: "Show detailed service status.", RunE: func(context.Context) error { return nil }},
		},
	})
	cmd, _, err := root.Find([]string{"status"})
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Long != "Show detailed service status." {
		t.Fatalf("long help = %q", cmd.Long)
	}
}

func TestChangedFlags(t *testing.T) {
	cmd := &Command{Use: "test"}
	cmd.Flags().String("local", "", "")
	cmd.PersistentFlags().String("global", "", "")
	if err := cmd.ParseFlags([]string{"--local", "a", "--global", "b"}); err != nil {
		t.Fatal(err)
	}
	got := ChangedFlags(cmd)
	if !got["local"] || !got["global"] {
		t.Fatalf("changed flags missing values: %#v", got)
	}
}

func TestNewCommandAndCompletionHelpers(t *testing.T) {
	called := false
	cmd := NewCommand(CommandOptions{
		Use:   "status",
		Short: "show status",
		RunE: func(context.Context, []string) error {
			called = true
			return nil
		},
	})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("expected RunE to be called")
	}

	values, directive := EnumCompletion("debug", "release")(nil, nil, "")
	if !reflect.DeepEqual(values, []string{"debug", "release"}) || directive != ShellCompDirectiveNoFileComp {
		t.Fatalf("unexpected enum completion: %#v %v", values, directive)
	}

	values, directive = NoFileCompletion(func() []string { return []string{"one"} })(nil, nil, "")
	if !reflect.DeepEqual(values, []string{"one"}) || directive != ShellCompDirectiveNoFileComp {
		t.Fatalf("unexpected no-file completion: %#v %v", values, directive)
	}
}

func TestRequireExactlyOne(t *testing.T) {
	if err := RequireExactlyOne(map[string]string{"execute": "op", "restore": ""}); err != nil {
		t.Fatal(err)
	}
	if err := RequireExactlyOne(map[string]string{"execute": "", "restore": ""}); err == nil {
		t.Fatal("expected missing value error")
	}
	if err := RequireExactlyOne(map[string]string{"execute": "op", "restore": "snap"}); err == nil {
		t.Fatal("expected conflict error")
	}
}

func commandNames(cmd *Command) []string {
	commands := cmd.Commands()
	names := make([]string, len(commands))
	for i, command := range commands {
		names[i] = command.Name()
	}
	return names
}
