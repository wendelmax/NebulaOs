package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the NebulaOS control plane",
	Long: `Creates the directory structure and default configuration
for the NebulaOS control plane at /var/lib/nebula/.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dirs := []string{
			"/var/lib/nebula",
			"/var/lib/nebula/data",
			"/var/lib/nebula/logs",
			"/var/lib/nebula/ssl",
			"/etc/nebula",
		}
		for _, d := range dirs {
			if err := os.MkdirAll(d, 0755); err != nil {
				return fmt.Errorf("creating %s: %w", d, err)
			}
			fmt.Printf("✓  %s\n", d)
		}
		cfg := `# NebulaOS configuration
host: 0.0.0.0
port: 8000
db_path: /var/lib/nebula/data/nebula.db
log_dir: /var/lib/nebula/logs
ssl_dir: /var/lib/nebula/ssl
`
		cfgPath := filepath.Join("/etc/nebula", "nebula.yaml")
		if err := os.WriteFile(cfgPath, []byte(cfg), 0644); err != nil {
			return fmt.Errorf("writing config: %w", err)
		}
		fmt.Printf("✓  %s\n", cfgPath)
		fmt.Println("\nNebulaOS control plane initialized. Run 'nebula install' to configure.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
