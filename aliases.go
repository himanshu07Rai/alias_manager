package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func loadAliases() (map[string]string, error) {
	data, err := os.ReadFile(aliasFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return make(map[string]string), nil
	}
	var aliases map[string]string
	if err := json.Unmarshal(data, &aliases); err != nil {
		return nil, fmt.Errorf("corrupted aliases file at %s. Try deleting it", aliasFile)
	}
	if aliases == nil {
		return make(map[string]string), nil
	}
	return aliases, nil
}

func saveAliases(aliases map[string]string) {
	ensureConfigDir()
	data, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to marshal aliases: %v\n", err)
		os.Exit(1)
	}
	if err := atomicWriteFile(aliasFile, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to save aliases: %v\n", err)
		os.Exit(1)
	}
}

func isValidAliasName(name string) bool {
	if name == "" {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\-]+$`, name)
	return matched
}

func createAlias(name, command string, force bool) {
	if !isValidAliasName(name) {
		fmt.Fprintf(os.Stderr, "  ✗ Invalid alias name '%s'. Use only letters, numbers, _, and -.\n", name)
		os.Exit(1)
	}
	if command == "" {
		fmt.Fprintln(os.Stderr, "  ✗ Command cannot be empty")
		os.Exit(1)
	}

	aliases, err := loadAliases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	if _, exists := aliases[name]; exists && !force {
		fmt.Fprintf(os.Stderr, "  ✗ Alias '%s' already exists. Use --force to overwrite.\n", name)
		os.Exit(1)
	}

	aliases[name] = command
	saveAliases(aliases)
	rebuildShellFile()

	fmt.Printf("  ✓ Created: %s='%s'\n", name, command)
}

func deleteAlias(name string, yes bool) {
	aliases, err := loadAliases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	if _, exists := aliases[name]; !exists {
		fmt.Fprintf(os.Stderr, "  ✗ Alias '%s' not found.\n", name)
		os.Exit(1)
	}

	if !yes {
		fmt.Printf("  Delete '%s'? [y/N]: ", name)
		if !promptConfirm("") {
			fmt.Println("  Cancelled.")
			return
		}
	}

	delete(aliases, name)
	saveAliases(aliases)
	rebuildShellFile()

	fmt.Printf("  ✓ Deleted: %s\n", name)
}

func editAlias(name string, newCommand *string) {
	aliases, err := loadAliases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	oldCmd, exists := aliases[name]
	if !exists {
		fmt.Fprintf(os.Stderr, "  ✗ Alias '%s' not found.\n", name)
		os.Exit(1)
	}

	if newCommand != nil {
		if *newCommand == "" {
			fmt.Fprintln(os.Stderr, "  ✗ Command cannot be empty")
			os.Exit(1)
		}
		aliases[name] = *newCommand
		saveAliases(aliases)
		rebuildShellFile()
		fmt.Printf("  ✓ Updated: %s='%s'\n", name, *newCommand)
		return
	}

	fmt.Printf("  Current: %s='%s'\n", name, oldCmd)
	input := promptInput("  New command (empty to cancel): ")
	if input == "" {
		fmt.Println("  Cancelled.")
		return
	}

	aliases[name] = input
	saveAliases(aliases)
	rebuildShellFile()
	fmt.Printf("  ✓ Updated: %s='%s'\n", name, input)
}

func listAliases(pattern string) {
	aliases, err := loadAliases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	if len(aliases) == 0 {
		fmt.Println("  No aliases defined. Run 'alias_manager create' to add one.")
		return
	}

	filtered := make(map[string]string)
	for name, cmd := range aliases {
		if pattern == "" || strings.Contains(strings.ToLower(name), strings.ToLower(pattern)) {
			filtered[name] = cmd
		}
	}

	if len(filtered) == 0 {
		fmt.Printf("  No aliases matching '%s'.\n", pattern)
		return
	}

	maxNameLen := len("Alias")
	maxCmdLen := len("Command")
	for name, cmd := range filtered {
		if len(name) > maxNameLen {
			maxNameLen = len(name)
		}
		if len(cmd) > maxCmdLen {
			maxCmdLen = len(cmd)
		}
	}
	if maxNameLen > 40 {
		maxNameLen = 40
	}
	if maxCmdLen > 80 {
		maxCmdLen = 80
	}

	colWidth := maxNameLen + maxCmdLen + 4
	if colWidth > 124 {
		colWidth = 124
	}

	fmt.Printf("  %-*s  %s\n", maxNameLen, "Alias", "Command")
	fmt.Println("  " + strings.Repeat("─", colWidth))

	for name, cmd := range filtered {
		displayName := name
		displayCmd := cmd
		if len(displayName) > 40 {
			displayName = displayName[:37] + "..."
		}
		if len(displayCmd) > 80 {
			displayCmd = displayCmd[:77] + "..."
		}
		fmt.Printf("  %-*s  %s\n", maxNameLen, displayName, displayCmd)
	}

	fmt.Printf("\n  Total: %d alias(es)\n", len(filtered))
}

func importAliases(filePath string, overwrite bool) {
	data, err := os.ReadFile(expandPath(filePath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to read file: %v\n", err)
		os.Exit(1)
	}

	var imported map[string]string
	if err := json.Unmarshal(data, &imported); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Invalid JSON format: %v\n", err)
		os.Exit(1)
	}

	aliases, err := loadAliases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	importedCount := 0
	skippedCount := 0

	for name, cmd := range imported {
		if _, exists := aliases[name]; exists && !overwrite {
			skippedCount++
			continue
		}
		aliases[name] = cmd
		importedCount++
	}

	if importedCount > 0 {
		saveAliases(aliases)
		rebuildShellFile()
	}

	fmt.Printf("  ✓ Imported %d, skipped %d\n", importedCount, skippedCount)
}

func exportAliases(filePath string) {
	aliases, err := loadAliases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	data, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to marshal aliases: %v\n", err)
		os.Exit(1)
	}

	outPath := expandPath(filePath)
	if err := atomicWriteFile(outPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("  ✓ Exported %d aliases to %s\n", len(aliases), filePath)
}
