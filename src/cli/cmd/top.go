package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

type topModel struct {
	spinner    spinner.Model
	loading    bool
	containers []containerInfo
	err        error
	tick       time.Time
}

func initialTopModel() topModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return topModel{
		spinner: s,
		loading: true,
	}
}

func (m topModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tickTop, fetchTopData)
}

func tickTop() tea.Msg {
	time.Sleep(3 * time.Second)
	return topTickMsg{}
}

func fetchTopData() tea.Msg {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errMsg{err}
	}
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return errMsg{err}
	}
	var infos []containerInfo
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = c.Names[0][1:]
		}
		infos = append(infos, containerInfo{
			name:   name,
			image:  c.Image,
			status: c.Status,
			state:  c.State,
		})
	}
	return containersMsg(infos)
}

type topTickMsg struct{}

func (m topModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil
	case containersMsg:
		m.containers = []containerInfo(msg)
		m.loading = false
		return m, tea.Batch(tickTop, fetchTopData)
	case topTickMsg:
		return m, fetchTopData
	case errMsg:
		m.err = msg.err
		m.loading = false
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m topModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Padding(0, 1).
		Render("✦ NebulaOS — Live Resource Monitor (auto-refresh 3s)")

	cols := []table.Column{
		{Title: "Service", Width: 28},
		{Title: "Image", Width: 36},
		{Title: "Status", Width: 40},
		{Title: "State", Width: 12},
	}

	var rows []table.Row
	for _, c := range m.containers {
		stateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
		if c.state != "running" {
			stateStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		}
		rows = append(rows, table.Row{
			c.name,
			c.image,
			c.status,
			stateStyle.Render(c.state),
		})
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(rows) + 2),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)

	content := lipgloss.NewStyle().Padding(1, 2).Render(t.View())

	updated := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(fmt.Sprintf("Last updated: %s", time.Now().Format("15:04:05")))

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press q to quit • auto-refreshes every 3s")

	return lipgloss.JoinVertical(lipgloss.Top,
		header,
		"",
		content,
		updated,
		"",
		footer,
	)
}

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Live resource consumption monitor",
	Long:  "Auto-refreshing TUI showing real-time container resource usage across the NebulaOS plane.",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(initialTopModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(topCmd)
}
