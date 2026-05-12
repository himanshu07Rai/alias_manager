# alias_manager — Specification

## Overview

CLI tool for managing shell aliases persistently across sessions.
Supports bash, zsh, and fish shells.

## Project Structure

```
alias_manager/
├── go.mod
├── main.go          Entry point, command routing, utilities
├── config.go        Config load/save, init logic, path management
├── aliases.go       CRUD operations on aliases
├── shell.go         Shell integration, RC file injection, rebuild
├── suggest.go       Built-in alias suggestions by category
├── tui.go           Bubbletea interactive TUI for create/edit
├── README.md
├── SPEC.md
└── TODOS.md
```

## Storage

| File | Purpose |
|---|---|
| `~/.alias_manager/config.json` | Tool configuration (shell type, RC path) |
| `~/.alias_manager/aliases.json` | Source of truth — JSON map of aliases |
| `~/.alias_manager/aliases.sh` | Bash/zsh shell file (derived from JSON) |
| `~/.alias_manager/aliases.fish` | Fish shell file (derived from JSON) |

### Config schema

```json
{
  "shell_rc": "/home/user/.zshrc",
  "alias_file": "/home/user/.alias_manager/aliases.json",
  "shell": "bash"
}
```

### Storage invariants

- `aliases.json` is the source of truth; shell files are derived
- After every mutation (create/delete/edit/import), shell files are regenerated
- Shell files are written atomically (write to `.tmp`, rename)
- If shell file is deleted/corrupted, `rebuild-shell-file` recovers from JSON

## Commands

### `init`

- Detect shell from `$SHELL` env var
- Display detected shell and default RC path
- Prompt user to confirm or override (shell type, RC path)
- Save config to `~/.alias_manager/config.json`
- Optionally create a shortcut alias (e.g. `am`) in the RC file

### `list [pattern]`

- Table output: `Alias | Command`
- Optional case-insensitive substring filter on name
- Show total count at bottom
- Friendly message if no aliases exist

### `create <name> <command> [--force]`

- Multi-word commands supported (join remaining non-flag args)
- Reject names with spaces or special chars (only `[a-zA-Z0-9_-]`)
- If name exists and no `--force`, print error with hint
- If no args provided and running in a TTY, launch bubbletea TUI form
- Print success with source hint

### `delete <name> [--yes]`

- Prompt confirmation `[y/N]` unless `--yes`/`-y` passed
- Error if alias doesn't exist
- Print success

### `edit <name> [command]`

- If command arg provided: update directly
- If no command arg and running in a TTY: launch bubbletea TUI with current value pre-filled
- If no command arg and not a TTY: show current value, prompt with plain text input
- Cancel on empty input or Ctrl+C

### `import <file> [--overwrite]`

- Accept JSON file in same format as `aliases.json`
- Skip duplicates unless `--overwrite`
- Report: `Imported N, skipped M`

### `export <file>`

- Write aliases to given path as JSON
- Print count

### `suggest [category]`

- Table: `Alias | Command | Status` (NEW or EXISTS)
- Categories: navigation, git, system, utility, fun
- If no category, show all
- Unknown category prints available categories

### `add-to-shell`

- Read RC path from config (error if not initialized)
- Generate shell function wrapper and append to RC file
- Detect and replace existing integration blocks
- Print success with source command

### `rebuild-shell-file`

- Read aliases from JSON, rewrite shell file
- Print "Rebuilt aliases.sh with N aliases"

### `shell-integration`

- Print alias definitions in shell-compatible format
- Used internally by shell hooks

### `help`

- Print full usage with all commands and flags

## Auto-Source (Shell Function Wrapper)

`add-to-shell` adds this function to the RC file:

```
alias_manager() {
    local f="$HOME/.alias_manager/aliases.sh"
    local old_names
    if [ -f "$f" ]; then
        old_names=$(sed -n 's/^alias \([^=]*\)=.*/\1/p' "$f")
    fi
    /path/to/alias_manager "$@"
    if [ -f "$f" ]; then
        while IFS= read -r name; do
            [ -n "$name" ] && unalias "$name" 2>/dev/null
        done <<< "$old_names"
        source "$f"
    fi
}
```

The function:
1. Saves current alias names from shell file
2. Runs the binary
3. Unaliases all previously-known names (handles deletions)
4. Sources the shell file (handles creates/updates)

## Output Style

| Prefix | Type | Fd |
|---|---|---|
| `✓` | Success | stdout |
| `✗` | Error | stderr |
| `!` | Warning | stderr |
| `→` | Info/next steps | stdout |

## Edge Cases

- `init` not run → `✗ Run 'alias_manager init' first` on commands needing config
- Corrupted JSON → Clear error with suggestion to delete the file
- Empty name/command → Rejected
- Name with spaces → Rejected
- Ctrl+C during prompt → Clean exit, no partial writes
- RC file doesn't exist → Warning, don't crash
- Unsupported shell → Warning, default to bash behavior
- Non-TTY environment → Fall back to plain text prompts instead of TUI

## Dependencies

- **Core**: Go stdlib only (no third-party deps for core logic)
- **TUI**: github.com/charmbracelet/bubbletea, github.com/charmbracelet/bubbles

## Build

```sh
go build -o alias_manager .
```
