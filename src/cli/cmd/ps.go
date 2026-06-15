package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/wendelmax/nebulaos/src/cli/internal/docker"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running NebulaOS services",
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return fmt.Errorf("docker client: %w", err)
		}
		containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
		if err != nil {
			return fmt.Errorf("list containers: %w", err)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tIMAGE\tSTATUS\tPORTS")
		for _, c := range containers {
			if !docker.IsNebulaContainer(c) {
				continue
			}
			name := ""
			if len(c.Names) > 0 {
				name = c.Names[0][1:]
			}
			ports := ""
			for i, p := range c.Ports {
				if i > 0 {
					ports += ", "
				}
				ports += fmt.Sprintf("%d→%d/%s", p.PublicPort, p.PrivatePort, p.Type)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, c.Image, c.Status, ports)
		}
		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
