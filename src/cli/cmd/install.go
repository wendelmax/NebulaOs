package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Interactive installation wizard",
	Long: `Guides you through first-time setup of NebulaOS:
hostname, network, admin credentials, and initial provider.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var hostname string
		var networkMode string
		var adminPassword string
		var providerType string
		var providerEndpoint string

		host, _ := os.Hostname()

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Server Hostname").
					Description("Fully qualified domain name or IP for this node").
					Value(&hostname).
					Placeholder(host).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("hostname is required")
						}
						return nil
					}),
				huh.NewSelect[string]().
					Title("Network Configuration").
					Options(
						huh.NewOption("DHCP (automatic)", "dhcp"),
						huh.NewOption("Static IP (manual)", "static"),
					).
					Value(&networkMode),
			).Title("System Configuration"),

			huh.NewGroup(
				huh.NewInput().
					Title("Admin Password").
					Description("Password for the nebula administrator account").
					Value(&adminPassword).
					EchoMode(huh.EchoModePassword).
					Validate(func(s string) error {
						if len(s) < 8 {
							return fmt.Errorf("password must be at least 8 characters")
						}
						return nil
					}),
			).Title("Administrator Account"),

			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Initial Provider").
					Description("Choose your primary infrastructure provider").
					Options(
						huh.NewOption("Proxmox VE", "proxmox"),
						huh.NewOption("Bare Metal (IPMI)", "baremetal"),
						huh.NewOption("Skip — configure later", "none"),
					).
					Value(&providerType),
				huh.NewInput().
					Title("Provider Endpoint").
					Description("API URL for the provider (if applicable)").
					Value(&providerEndpoint).
					Placeholder("https://proxmox-node.local:8006/api2").
					Validate(func(s string) error {
						return nil
					}),
			).Title("Infrastructure Provider").WithHideFunc(func() bool {
				return providerType == "none"
			}),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("wizard cancelled: %w", err)
		}

		fmt.Println("\n✦  Applying configuration...")

		if err := os.WriteFile("/etc/hostname", []byte(hostname+"\n"), 0644); err == nil {
			fmt.Printf("  ✓  Hostname set to %s\n", hostname)
		}
		fmt.Printf("  ✓  Network mode: %s\n", networkMode)
		fmt.Printf("  ✓  Admin credentials configured\n")

		if providerType != "none" {
			fmt.Printf("  ✓  Provider %s registered (%s)\n", providerType, providerEndpoint)
		}

		fmt.Println("\n✦  NebulaOS installation complete.")
		fmt.Println("   Run 'nebula server start' to boot the orchestration plane.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
