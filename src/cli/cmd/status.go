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

type statusModel struct {
	spinner   spinner.Model
	loading   bool
	containers []containerInfo
	err       error
	width     int
	height    int
}

type containerInfo struct {
	name   string
	image  string
	status string
	state  string
}

func initialStatusModel() statusModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return statusModel{
		spinner: s,
		loading: true,
	}
}

func (m statusModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchContainers)
}

func fetchContainers() tea.Msg {
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

type containersMsg []containerInfo
type errMsg struct{ err error }

func (m statusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil
	case containersMsg:
		m.containers = []containerInfo(msg)
		m.loading = false
		return m, nil
	case errMsg:
		m.err = msg.err
		m.loading = false
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m statusModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit.\n", m.err)
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Padding(0, 1).
		Render("✦ NebulaOS — Control Plane Status")

	sub := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))

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

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Press q to quit • Press r to refresh")

	return lipgloss.JoinVertical(lipgloss.Top,
		header,
		sub,
		"",
		content,
		"",
		footer,
	)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Interactive dashboard overview",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(initialStatusModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
