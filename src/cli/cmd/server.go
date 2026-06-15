package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage the NebulaOS orchestration engine",
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the NebulaOS stack",
	Long: `Boots the full NebulaOS orchestration plane via Docker Compose.
Uses deploy/local/docker-compose.yml relative to the project root.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot := findProjectRoot()
		composeFile := projectRoot + "/deploy/local/docker-compose.yml"
		if _, err := os.Stat(composeFile); os.IsNotExist(err) {
			return fmt.Errorf("compose file not found at %s. Run from project root or set NEBULA_HOME", composeFile)
		}
		fmt.Println("✦  Starting NebulaOS orchestration plane...")
		c := exec.Command("docker-compose", "-f", composeFile, "up", "-d")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("docker-compose up failed: %w", err)
		}
		fmt.Println("\n✓  NebulaOS is running.")
		fmt.Println("   Dashboard: http://nebula.local")
		fmt.Println("   API:       http://api.nebula.local:8000")
		return nil
	},
}

var serverStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the NebulaOS stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot := findProjectRoot()
		composeFile := projectRoot + "/deploy/local/docker-compose.yml"
		fmt.Println("✦  Stopping NebulaOS...")
		c := exec.Command("docker-compose", "-f", composeFile, "down")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("docker-compose down failed: %w", err)
		}
		fmt.Println("✓  NebulaOS stopped.")
		return nil
	},
}

func findProjectRoot() string {
	if h := os.Getenv("NEBULA_HOME"); h != "" {
		return h
	}
	dir, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(dir + "/go.work"); err == nil {
			return dir
		}
		dir = dir + "/.."
	}
	return "."
}

func init() {
	serverCmd.AddCommand(serverStartCmd)
	serverCmd.AddCommand(serverStopCmd)
	rootCmd.AddCommand(serverCmd)
}
