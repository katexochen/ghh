package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/katexochen/ghh/internal/cmd"

	"github.com/spf13/cobra"
)

var (
	version = "0.0.0-dev"
	commit  = "HEAD"
	date    = "unknown"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	cobra.EnableCommandSorting = false
	rootCmd := newRootCmd()
	ctx, cancel := signalContext(context.Background(), os.Interrupt)
	defer cancel()
	return rootCmd.ExecuteContext(ctx)
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "ghh",
		Short:            "GitHub Helper CLI",
		Version:          version,
		PersistentPreRun: preRunRoot,
	}

	rootCmd.SetOut(os.Stdout)
	rootCmd.AddCommand(
		cmd.NewDeleteAllRunsCmd(),
		cmd.NewCreateProjectIssueCmd(),
		cmd.NewSetAuthCmd(),
	)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.InitDefaultVersionFlag()
	rootCmd.SetVersionTemplate(
		fmt.Sprintf("ghh - GitHub helper CLI\n\nversion   %s\ncommit    %s\nbuilt at  %s\n", version, commit, date),
	)

	return rootCmd
}

func signalContext(ctx context.Context, sig os.Signal) (context.Context, context.CancelFunc) {
	sigCtx, stop := signal.NotifyContext(ctx, sig)
	done := make(chan struct{}, 1)
	stopDone := make(chan struct{}, 1)

	go func() {
		defer func() { stopDone <- struct{}{} }()
		defer stop()
		select {
		case <-sigCtx.Done():
			fmt.Println(" Signal caught. Press ctrl+c again to terminate the program immediately.")
		case <-done:
		}
	}()

	cancelFunc := func() {
		done <- struct{}{}
		<-stopDone
	}

	return sigCtx, cancelFunc
}

func preRunRoot(cmd *cobra.Command, _ []string) {
	cmd.SilenceUsage = true
}
