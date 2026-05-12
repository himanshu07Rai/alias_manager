# TODOS

## Priority

- [ ] **Add `add-to-path` command** — Add `alias_manager` binary directory to `$PATH` in RC file so `alias_manager` and the `am` shortcut work without needing it in PATH
- [ ] **Configurable shortcut name** — Store the chosen shortcut alias name (e.g. `am`) in `config.json` and reuse it in `add-to-shell` instead of hardcoding `am`
- [ ] **Shell completions** — Generate shell completions for bash/zsh/fish

## Polish

- [ ] **`--json` flag on `list`** — Output aliases as JSON for scripting
- [ ] **Bulk operations** — `alias_manager create --batch <file>` to create many at once
- [ ] **Rename command** — `alias_manager rename <old> <new>` for renaming an alias
- [ ] **Verbose mode** — `--verbose` flag for detailed operation logging
- [ ] **Confirm before overwrite in TUI** — When creating an alias that exists, show confirmation in the TUI instead of exiting

## Refactoring

- [ ] **Store shortcut name in config** — The shortcut alias name chosen during `init` should be persisted so `add-to-shell` can use the same name
- [ ] **Extract TUI models** — The `createModel` and `editModel` could share a base or be refactored into a reusable form component
- [ ] **Graceful degradation in login shells** — Ensure the function wrapper works in non-interactive login shells where alias expansion is disabled

## Testing

- [ ] **Unit tests for CRUD** — Test `createAlias`, `deleteAlias`, `editAlias` with mock filesystem
- [ ] **Integration test** — Test `init` → `add-to-shell` → `create` → `list` → `delete` in a temp home
- [ ] **JSON edge cases** — Test corrupted `aliases.json`, empty file, missing file
- [ ] **Shell format** — Verify generated `aliases.sh` and `aliases.fish` syntax with shellcheck
- [ ] **Non-TTY fallback** — Test `create` and `edit` when stdin is not a terminal

## Documentation

- [x] README.md — Quick start and command reference
- [x] SPEC.md — Detailed specification
- [x] TODOS.md — Project tracker
