package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

//
// Command system
//

type CommandAction interface {
	Run() (string, error)
}

type FuncCommand struct {
	Output string
}

func (f FuncCommand) Run() (string, error) {
	return f.Output, nil
}

type ExecCommand struct {
	Name string
	Args []string
}

func (e ExecCommand) Run() (string, error) {
	out, err := exec.Command(e.Name, e.Args...).CombinedOutput()
	return string(out), err
}

type Command struct {
	Name     string
	Category string
	Action   CommandAction
}

func (c Command) Title() string       { return c.Name }
func (c Command) Description() string { return fmt.Sprintf("Category: %s", c.Category) }
func (c Command) FilterValue() string { return c.Name }

//
// Config loader
//

type commandYAML struct {
	Name     string   `yaml:"name"`
	Category string   `yaml:"category"`
	Type     string   `yaml:"type"`
	Code     string   `yaml:"code"`
	Command  string   `yaml:"command"`
	Args     []string `yaml:"args"`
}

type config struct {
	Commands []commandYAML `yaml:"commands"`
}

func loadCommands(path string) ([]Command, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	var cmds []Command
	for _, c := range cfg.Commands {
		switch c.Type {
		case "func":
			cmds = append(cmds, Command{
				Name:     c.Name,
				Category: c.Category,
				Action:   FuncCommand{Output: c.Code},
			})
		case "exec":
			cmds = append(cmds, Command{
				Name:     c.Name,
				Category: c.Category,
				Action:   ExecCommand{Name: c.Command, Args: c.Args},
			})
		default:
			return nil, fmt.Errorf("unknown command type: %s", c.Type)
		}
	}
	return cmds, nil
}

//
// TUI
//

var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("9")).
	Bold(true)

type model struct {
	list       list.Model
	viewport   viewport.Model
	showOutput bool
	output     string
	err        error
}

type resultMsg struct {
	output string
	err    error
}

func initialModel(commands []Command) model {
	items := make([]list.Item, len(commands))
	for i, cmd := range commands {
		items[i] = cmd
	}

	const defaultWidth = 30
	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, 14)
	l.Title = "Select a command to run"

	vp := viewport.New(50, 14)
	vp.Style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

	return model{list: l, viewport: vp}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.showOutput {
				if selected, ok := m.list.SelectedItem().(Command); ok {
					return m, func() tea.Msg {
						out, err := selected.Action.Run()
						return resultMsg{out, err}
					}
				}
			} else {
				m.showOutput = false
			}
		case "esc", "q", "ctrl+c":
			if m.showOutput {
				m.showOutput = false
			} else {
				return m, tea.Quit
			}
		}
	case resultMsg:
		m.output = msg.output
		m.err = msg.err
		m.viewport.SetContent(m.output)
		m.showOutput = true
	}

	if m.showOutput {
		m.viewport, cmd = m.viewport.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.showOutput {
		var b strings.Builder
		b.WriteString("Output:\n")
		b.WriteString(m.viewport.View())
		if m.err != nil {
			b.WriteString("\n")
			b.WriteString(errorStyle.Render(m.err.Error()))
		}
		b.WriteString("\n\n[enter/esc to return]")
		return b.String()
	}
	return m.list.View() + "\n\n[q to quit]"
}

func main() {
	commands, err := loadCommands("commands.yaml")
	if err != nil {
		fmt.Println("Error loading commands:", err)
		os.Exit(1)
	}

	if err := tea.NewProgram(initialModel(commands)).Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
