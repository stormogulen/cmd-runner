# Go TUI Command Runner

A terminal user interface (TUI) tool for running system commands and custom Go functions with style.
Built using [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss), and YAML configuration.

## Features

- **YAML-based configuration**
  Define commands and categories without touching the code.

- **Scrollable command output**
  Command results stream into a Bubble Tea viewport.

- **Easy packaging**
  Build and distribute as a self-extracting archive with `makeself` and a single `Makefile`.

## Getting Started

1. **Clone the repo**
   ```bash
   git clone https://github.com/stormogulen/cmd-runner.git
   cd cmd-runner
   ```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Run it**
```bash
go run .
```

4. **Build binary**
```bash
go build -o cmd-runner
``

## Configuration
Commands are defined in config.yaml. Example:

```yaml
commands:
  - name: List files
    type: exec
    command: ls
    args: ["-l", "-a"]

  - name: Show date
    type: exec
    command: date

  - name: Custom Go Function
    type: go
    func: customFunction
```
* exec → Runs an external system command.

* go → Runs a Go function implemented inside the app.

## Packaging

To create a self-extracting archive:
```bash
make package
```
This produces a command.run script that contains the binary.
