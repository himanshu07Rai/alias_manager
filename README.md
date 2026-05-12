# alias_manager

A CLI tool to manage shell aliases persistently across sessions. Supports **bash**, **zsh**, and **fish**.

## Install

```sh
go install alias_manager@latest
# or build from source:
git clone <repo>
cd alias_manager
go build -o alias_manager .
sudo mv alias_manager /usr/local/bin/
```

## Quick Start

```sh
# 1. Initialize
alias_manager init

# 2. Register with shell
alias_manager add-to-shell && source ~/.zshrc

# 3. Create aliases
alias_manager create gs "git status"
alias_manager create ll "ls -la"

# 4. Use them immediately (auto-sourced in current session)
gs
ll

# Or use the shortcut alias (configured during init)
am create gp "git push"
```

## Commands

| Command | Description |
|---|---|
| `init` | Detect shell, configure RC file |
| `list [pattern]` | List aliases, filter by name |
| `create <name> <command> [--force]` | Create alias (interactive TUI if no args) |
| `delete <name> [--yes]` | Delete an alias |
| `edit <name> [command]` | Edit alias (interactive TUI if no command) |
| `import <file> [--overwrite]` | Import aliases from JSON |
| `export <file>` | Export aliases to JSON |
| `suggest [category]` | Show alias suggestions |
| `add-to-shell` | Add shell integration hook |
| `rebuild-shell-file` | Rebuild shell file from JSON |
| `help` | Show usage |

## Features

- **Auto-source** — Changes apply to the current shell session immediately; no manual `source` needed
- **Interactive TUI** — `create` and `edit` launch a bubbletea form when run without required args
- **Suggestions** — Built-in alias ideas organized by category (git, navigation, system, utility, fun)
- **Persistent** — Aliases stored in `~/.alias_manager/aliases.json`, loaded on every shell start
- **Atomic writes** — All files written atomically (`.tmp` + rename); no partial writes on crash
- **Shell-agnostic** — Generates correct syntax for bash, zsh, and fish
- **Recovery** — `rebuild-shell-file` recovers the shell file if accidentally deleted
- **No third-party deps (core)** — Core logic uses stdlib only; bubbletea used for interactive TUI

## Storage

- `~/.alias_manager/aliases.json` — Source of truth (JSON map of name→command)
- `~/.alias_manager/aliases.sh` — Derived shell file (bash/zsh)  
- `~/.alias_manager/aliases.fish` — Derived shell file (fish)
- `~/.alias_manager/config.json` — Tool configuration

## Shell Integration

`add-to-shell` adds a shell function wrapper to your RC file that:

1. Runs the alias_manager binary
2. Unaliases previous alias names from the current session
3. Re-sources `aliases.sh`/`aliases.fish` so changes take effect immediately

No manual `source` needed after `create`/`edit`/`delete`/`import`.
