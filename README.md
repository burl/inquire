# inquire

Lightweight, **line-oriented** interactive CLI prompts for Go — an [inquirer.js](https://www.npmjs.com/package/inquirer)-style experience without taking over the full terminal.

Each question renders as a small inline **band** at the current cursor. Answered prompts settle to static `✔ …` scrollback; the next question opens below. No alternate screen, no full-screen repaint.

```bash
go get github.com/burl/inquire/v2
```

Requires **Go 1.26+**. Unix terminals (Linux, macOS, BSD) are supported; Windows is not in v2.0.

**API docs:** [pkg.go.dev/github.com/burl/inquire/v2](https://pkg.go.dev/github.com/burl/inquire/v2) (indexed on tagged releases).

## Quick start

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "os"

    "github.com/burl/inquire/v2"
)

func main() {
    var name string
    var ok bool

    err := inquire.Query().
        Input(&name, "what is your name", nil).
        YesNo(&ok, "continue", nil).
        Run(context.Background())
    if errors.Is(err, inquire.ErrInterrupted) {
        fmt.Fprintln(os.Stderr, "interrupted")
        os.Exit(130)
    }
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    fmt.Println(name, ok)
}
```

## TTY requirements

`Run` requires **both stdin and stdout** to be terminals. Redirecting either stream (e.g. piping, CI without a pseudo-TTY) returns `inquire.ErrNotTerminal`. Stderr may be redirected freely.

Use `inquire.WithInput` / `inquire.WithOutput` when embedding in tests or alternate file descriptors.

## Interrupts

**Ctrl+C** aborts the entire session and returns `inquire.ErrInterrupted`. Answers from prompts already completed remain in bound variables; later prompts are not run. Your application should treat this as a full cancel of the remaining flow.

## Widgets

| Widget | Binds | Notes |
|--------|-------|-------|
| `Input` | `*string` | Defaults, validation, password mask |
| `YesNo` | `*bool` | Arrows, y/n, space |
| `Menu` | `*string` (tag) | Vertical single-select |
| `Select` | `*bool` per item | Multi-select checkboxes |
| `Note` | — | Non-interactive message; Enter to continue |
| `AnyKey` | — | Continues on any key |

Configure every widget in the `more` callback:

```go
inquire.Query().
    Menu(&quest, "your quest", func(w *widget.Menu) {
        w.Item("grail", "find the grail")
        w.When(widget.WhenEqual(&mode, "advanced"))
    })
```

Conditional prompts use `When` with `widget.WhenEqual` predicates.

## Keybindings

| Context | Keys |
|---------|------|
| **All** | Ctrl+C — abort session |
| **Input** | Arrows, Home/End, Backspace/Delete; Ctrl+A/E/K/D/W; Ctrl+B/F (left/right); Emacs navigation |
| **YesNo** | ←/→, ↑/↓, y/n, Space (toggle), Enter (confirm) |
| **Menu** | ↑/↓, ←/→ (move); Space/Enter (confirm); Ctrl+P/N (up/down) |
| **Select** | ↑/↓ (move); Space (toggle item); a (select all); i (invert); Enter (confirm) |
| **Note** | Enter (continue) |
| **AnyKey** | Any key (continue) |

Footer hints are shown on Menu and Select while active.

## Recipes

### Conditional follow-up

```go
inquire.Query().
    Menu(&quest, "your quest", func(w *widget.Menu) {
        w.Item("grail", "find the grail")
        w.Item("nuts", "find coconuts")
    }).
    Input(&weight, "swallow weight", func(w *widget.Input) {
        w.When(widget.WhenEqual(&quest, "nuts"))
        w.Valid(func(s string) string {
            if s == "" { return "required" }
            return ""
        })
    })
```

### Section divider

```go
inquire.Query().
    Note("Database setup", nil).
    Input(&host, "host", nil).
    Input(&port, "port", nil)
```

### Force monochrome output

```go
inquire.Query(inquire.WithColor(false)).
    Input(&name, "name", nil)
```

### Embed with signal handling

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()

err := inquire.Query().Input(&name, "name", nil).Run(ctx)
if errors.Is(err, inquire.ErrInterrupted) {
    // user pressed Ctrl+C
}
```

## Demo

```bash
task demo
```

Runs [`demo/grail.go`](demo/grail.go) — the Monty Python Bridge of Death quest through every widget (Note, AnyKey, Input, Menu, Select, YesNo).

Lower-level band/terminal demo:

```bash
task demo:termui
```

## Development

Requires [Task](https://taskfile.dev). Common targets:

```bash
task check       # vet, test, lint
task build       # inquire-grail for current host → bin.out/
task build:all   # cross-compile demos for linux+darwin (amd64+arm64)
```

### CI and releases

| Workflow | Trigger | What it does |
|----------|---------|--------------|
| [`ci.yml`](.github/workflows/ci.yml) | PRs to `master` / `develop` | `task check` (ubuntu + macos) + `task build:all` |
| [`release.yml`](.github/workflows/release.yml) | Push to `master` | Same gate, then tags `v2.0.N` and publishes a short API guide |

Enable branch protection on `master` and require the **ci** jobs (`check`, `build`) before merge.

## Migration from v1

See [docs/MIGRATION.md](docs/MIGRATION.md).

## How this differs from full-screen TUIs

Libraries like Bubble Tea, huh, and tview assume an **application model**: alternate buffer, cleared screen, global event loop. That is the right fit for full dashboards and rich TUIs.

**inquire** targets a narrower case: ask a few questions in the middle of an existing CLI, leave prior output in scrollback, and return control to the caller via `Run(ctx) error` — no `os.Exit`, no init panics.

## v1

The original module (`github.com/burl/inquire`, termbox-based) is frozen. It remains available via git history and existing tags if you need the legacy API.

## License

[MIT](LICENSE)