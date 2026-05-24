package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/audit"
)

func buildAuditCmd(logPath string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "View audit log of envoy actions",
	}
	cmd.AddCommand(buildAuditListCmd(logPath))
	return cmd
}

func buildAuditListCmd(logPath string) *cobra.Command {
	var jsonOut bool
	var filterEnvSet string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List audit log entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := audit.NewLogger(logPath)
			if err != nil {
				return fmt.Errorf("open audit log: %w", err)
			}
			events, err := logger.ReadAll()
			if err != nil {
				return fmt.Errorf("read audit log: %w", err)
			}

			if filterEnvSet != "" {
				var filtered []audit.Event
				for _, e := range events {
					if e.EnvSet == filterEnvSet {
						filtered = append(filtered, e)
					}
				}
				events = filtered
			}

			if jsonOut {
				return json.NewEncoder(os.Stdout).Encode(events)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIMESTAMP\tTYPE\tENV_SET\tSUCCESS\tMESSAGE")
			for _, e := range events {
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n",
					e.Timestamp.Format("2006-01-02 15:04:05"),
					e.Type, e.EnvSet, e.Success, e.Message)
			}
			return w.Flush()
		},
	}

	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")
	cmd.Flags().StringVar(&filterEnvSet, "env-set", "", "Filter by env set name")
	return cmd
}
