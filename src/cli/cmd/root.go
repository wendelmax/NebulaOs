package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "nebula",
	Short: "NebulaOS — Sovereign Cloud Orchestrator",
	Long: `NebulaOS CLI: manage, deploy, and monitor your sovereign cloud infrastructure.

  nebula init          — initialize the control plane
  nebula server start  — start the orchestration engine
  nebula install       — interactive installation wizard
  nebula ps            — list running services
  nebula logs          — view service logs
  nebula status        — dashboard overview
  nebula top           — live resource consumption`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/etc/nebula/nebula.yaml", "config file path")
}
