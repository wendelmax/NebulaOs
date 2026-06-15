package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

type logsModel struct {
	viewport viewport.Model
	content  string
	err      error
	ready    bool
	service  string
}

func initialLogsModel(service string) logsModel {
	return logsModel{service: service}
}

func (m logsModel) Init() tea.Cmd {
	return fetchLogs(m.service)
}

func fetchLogs(service string) tea.Cmd {
	return func() tea.Msg {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return errMsg{err}
		}
		containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
		if err != nil {
			return errMsg{err}
		}
		var target string
		for _, c := range containers {
			for _, n := range c.Names {
				if strings.Contains(strings.ToLower(n), strings.ToLower(service)) {
					target = c.ID
					break
				}
			}
			if target != "" {
				break
			}
		}
		if target == "" {
			return errMsg{fmt.Errorf("no container found matching %q", service)}
		}
		reader, err := cli.ContainerLogs(context.Background(), target, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Tail:       "100",
		})
		if err != nil {
			return errMsg{err}
		}
		defer reader.Close()
		data, _ := io.ReadAll(reader)
		return logsContent(string(data))
	}
}

type logsContent string

func (m logsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width-4, msg.Height-6)
			m.viewport.YPosition = 2
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - 6
		}
		m.viewport.SetContent(m.content)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	case logsContent:
		m.content = string(msg)
		m.viewport.SetContent(m.content)
		return m, nil
	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	}
	return m, nil
}

func (m logsModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	if !m.ready {
		return "Loading logs..."
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Padding(0, 1).
		Render(fmt.Sprintf("✦ Logs: %s", m.service))

	body := lipgloss.NewStyle().
		Padding(0, 2).
		Render(m.viewport.View())

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("↑↓ scroll • q quit")

	return lipgloss.JoinVertical(lipgloss.Top, header, body, footer)
}

var logsCmd = &cobra.Command{
	Use:   "logs [service]",
	Short: "View service logs",
	Long: `Displays logs for a NebulaOS service container.

Examples:
  nebula logs nebula-api    — API server logs
  nebula logs postgres      — database logs
  nebula logs nebula-dashboard — dashboard logs`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]

		if !useTUI {
			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				return fmt.Errorf("docker client: %w", err)
			}
			containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
			if err != nil {
				return fmt.Errorf("list containers: %w", err)
			}
			var target string
			for _, c := range containers {
				for _, n := range c.Names {
					if strings.Contains(strings.ToLower(n), strings.ToLower(service)) {
						target = c.ID
						break
					}
				}
				if target != "" {
					break
				}
			}
			if target == "" {
				return fmt.Errorf("no container found matching %q", service)
			}
			reader, err := cli.ContainerLogs(context.Background(), target, container.LogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Tail:       "100",
			})
			if err != nil {
				return err
			}
			defer reader.Close()
			io.Copy(os.Stdout, reader)
			return nil
		}

		p := tea.NewProgram(initialLogsModel(service), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

var useTUI bool

func init() {
	logsCmd.Flags().BoolVarP(&useTUI, "tui", "t", true, "use TUI viewer (set to false for plain text output)")
	rootCmd.AddCommand(logsCmd)
}
