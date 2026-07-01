// Package clikit provides shared Cobra helpers for obsidzen CLI tools.
package clikit

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Command = cobra.Command
type ShellCompDirective = cobra.ShellCompDirective
type CompletionFunc = cobra.CompletionFunc
type PositionalArgs = cobra.PositionalArgs

const ShellCompDirectiveNoFileComp = cobra.ShellCompDirectiveNoFileComp

var (
	NoArgs        = cobra.NoArgs
	ArbitraryArgs = cobra.ArbitraryArgs
)

func ExactArgs(n int) PositionalArgs { return cobra.ExactArgs(n) }
func RangeArgs(min, max int) PositionalArgs {
	return cobra.RangeArgs(min, max)
}

type Action struct {
	Key                string
	Description        string
	Detail             string
	Aliases            []string
	Args               PositionalArgs
	DisableFlagParsing bool
	RunE               func(context.Context) error
	RunArgsE           func(context.Context, []string) error
}

type VersionInfo struct {
	Name    string
	Version string
	Commit  string
	Date    string
}

func (v VersionInfo) String() string {
	name := strings.TrimSpace(v.Name)
	version := strings.TrimSpace(v.Version)
	commit := strings.TrimSpace(v.Commit)
	date := strings.TrimSpace(v.Date)
	if version == "" {
		version = "dev"
	}
	if commit == "" {
		commit = "unknown"
	}
	if date == "" {
		date = "unknown"
	}
	if name == "" {
		return fmt.Sprintf("%s (commit %s, built %s)", version, commit, date)
	}
	return fmt.Sprintf("%s %s (commit %s, built %s)", name, version, commit, date)
}

type CommandOptions struct {
	Use               string
	Short             string
	Long              string
	Aliases           []string
	Args              PositionalArgs
	ValidArgsFunction CompletionFunc
	RunE              func(context.Context, []string) error
	Configure         func(*Command)
}

type RootOptions struct {
	Use               string
	Short             string
	Long              string
	Args              PositionalArgs
	ValidArgsFunction CompletionFunc
	RunE              func(context.Context, []string) error
	Menu              func(context.Context) error
	Actions           []Action
	Version           VersionInfo
	Configure         func(*Command)
}

func NewRoot(opts RootOptions) *Command {
	root := &cobra.Command{
		Use:               opts.Use,
		Short:             opts.Short,
		Long:              opts.Long,
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              opts.args(),
		ValidArgsFunction: opts.ValidArgsFunction,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.RunE != nil {
				return opts.RunE(cmd.Context(), args)
			}
			if opts.Menu == nil {
				return cmd.Help()
			}
			return opts.Menu(cmd.Context())
		},
	}
	if opts.Configure != nil {
		opts.Configure(root)
	}
	if opts.Menu != nil {
		root.AddCommand(&cobra.Command{
			Use:   "menu",
			Short: "Open interactive menu",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				return opts.Menu(cmd.Context())
			},
		})
	}
	for _, action := range opts.Actions {
		action := action
		root.AddCommand(&cobra.Command{
			Use:                action.Key,
			Short:              action.Description,
			Long:               action.Long(),
			Aliases:            action.Aliases,
			Args:               action.args(),
			DisableFlagParsing: action.DisableFlagParsing,
			RunE: func(cmd *cobra.Command, args []string) error {
				if action.RunArgsE != nil {
					return action.RunArgsE(cmd.Context(), args)
				}
				if action.RunE == nil {
					return fmt.Errorf("action is not implemented: %s", action.Key)
				}
				return action.RunE(cmd.Context())
			},
		})
	}
	root.AddCommand(CompletionCommand())
	if opts.Version.Name != "" || opts.Version.Version != "" || opts.Version.Commit != "" || opts.Version.Date != "" {
		version := opts.Version
		if version.Name == "" {
			version.Name = opts.Use
		}
		root.Version = version.String()
		root.SetVersionTemplate("{{.Version}}\n")
		root.AddCommand(VersionCommand(version))
	}
	return root
}

func NewCommand(opts CommandOptions) *Command {
	cmd := &cobra.Command{
		Use:               opts.Use,
		Short:             opts.Short,
		Long:              opts.Long,
		Aliases:           opts.Aliases,
		Args:              opts.args(),
		ValidArgsFunction: opts.ValidArgsFunction,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.RunE == nil {
				return nil
			}
			return opts.RunE(cmd.Context(), args)
		},
	}
	if opts.Configure != nil {
		opts.Configure(cmd)
	}
	return cmd
}

func (opts RootOptions) args() PositionalArgs {
	if opts.Args != nil {
		return opts.Args
	}
	return cobra.NoArgs
}

func (action Action) args() PositionalArgs {
	if action.Args != nil {
		return action.Args
	}
	return cobra.NoArgs
}

func (action Action) Long() string {
	if strings.TrimSpace(action.Detail) != "" {
		return action.Detail
	}
	return action.Description
}

func (opts CommandOptions) args() PositionalArgs {
	if opts.Args != nil {
		return opts.Args
	}
	return cobra.NoArgs
}

func Execute(root *Command) {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "✗", err)
		os.Exit(1)
	}
}

func CompletionCommand() *Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion script",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(os.Stdout)
			case "zsh":
				return root.GenZshCompletion(os.Stdout)
			case "fish":
				return root.GenFishCompletion(os.Stdout, true)
			case "powershell":
				return root.GenPowerShellCompletion(os.Stdout)
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
			}
		},
	}
}

func VersionCommand(version VersionInfo) *Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), version.String())
			return nil
		},
	}
}

func EnumCompletion(values ...string) CompletionFunc {
	return func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return values, cobra.ShellCompDirectiveNoFileComp
	}
}

func NoFileCompletion(fn func() []string) CompletionFunc {
	return func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return fn(), cobra.ShellCompDirectiveNoFileComp
	}
}

func MustRegisterFlagCompletion(cmd *Command, name string, fn CompletionFunc) {
	if err := cmd.RegisterFlagCompletionFunc(name, fn); err != nil {
		panic(err)
	}
}

func RequireExactlyOne(values map[string]string) error {
	var set []string
	for name, value := range values {
		if value != "" {
			set = append(set, name)
		}
	}
	if len(set) == 1 {
		return nil
	}
	if len(set) == 0 {
		return fmt.Errorf("exactly one of %s is required", joinedFlagNames(values))
	}
	return fmt.Errorf("only one of %s can be used", joinedFlagNames(values))
}

func joinedFlagNames(values map[string]string) string {
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, "--"+name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func ChangedFlags(cmd *Command) map[string]bool {
	changed := map[string]bool{}
	visit := func(flag *pflag.Flag) {
		changed[flag.Name] = true
	}
	cmd.Flags().Visit(visit)
	cmd.PersistentFlags().Visit(visit)
	cmd.InheritedFlags().Visit(visit)
	return changed
}
