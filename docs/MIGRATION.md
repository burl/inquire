# Migrating from inquire v1 to v2

v2 is a new module path with a rewritten terminal layer and API cleanup. v1 (`github.com/burl/inquire`, termbox-based) is frozen.

```bash
go get github.com/burl/inquire/v2
```

## Module and import path

| v1 | v2 |
|----|-----|
| `github.com/burl/inquire` | `github.com/burl/inquire/v2` |
| `github.com/burl/inquire/widget` | `github.com/burl/inquire/v2/widget` |

Update imports and `go.mod` accordingly.

## Session entry

v1 used global state and could panic or call `os.Exit`. v2 is explicit and embeddable:

```go
err := inquire.Query().
    Input(&name, "what is your name", nil).
    YesNo(&ok, "continue", nil).
    Run(ctx)
if errors.Is(err, inquire.ErrInterrupted) {
    // Ctrl+C — session aborted; earlier answers may already be bound
}
if errors.Is(err, inquire.ErrNotTerminal) {
    // stdin or stdout is not a TTY
}
```

## Widget configuration

All widgets use the closure configurators. v2 removed session-level `MenuItem` / `SelectItem`.

```go
// v2
inquire.Query().
    Menu(&quest, "your quest", func(w *widget.Menu) {
        w.Item("grail", "find the grail")
        w.Item("nuts", "find coconuts")
    }).
    Select("colors", func(w *widget.Select) {
        w.Item(&red, "red")
        w.Item(&blue, "blue")
    })
```

## Conditional prompts

Use `When` with the generic `widget.WhenEqual` predicate:

```go
w.When(widget.WhenEqual(&quest, "nuts")) // replaces v1 WhenEqualString
```

Go does not allow generic methods on widgets; `WhenEqual` is a package-level helper that returns a `func() bool`.

## New widgets

| Widget | Purpose |
|--------|---------|
| `Note` | Non-interactive message; Enter to continue |
| `AnyKey` | Waits for any key (except Ctrl+C) |

## YesNo configurators

`YesNo` now accepts an optional configurator like other widgets:

```go
YesNo(&ok, "continue", func(w *widget.YesNo) {
    w.Hint("y/n")
})
```

## Terminal backend

v2 replaces termbox with an owned inline band layer (`internal/termui`) on `golang.org/x/term`. Prompts render at the current cursor and settle to scrollback — no alternate screen.

Unix terminals (Linux, macOS, BSD) are supported. Windows is not supported in v2.0.

## Color

Color follows `NO_COLOR`, `TERM`, and `COLORTERM` by default. Force mono or color with:

```go
inquire.Query(inquire.WithColor(false))
```

## Dependencies

v2 drops vendored termbox. Runtime deps: `golang.org/x/term`, `github.com/mattn/go-runewidth`.

## API documentation

Tagged releases are indexed at [pkg.go.dev/github.com/burl/inquire/v2](https://pkg.go.dev/github.com/burl/inquire/v2).