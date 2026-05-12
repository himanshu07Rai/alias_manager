package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func shellEscape(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}

func generateShellContent(aliases map[string]string, shellType string) string {
	var b strings.Builder
	now := time.Now().Format("2006-01-02 15:04:05")

	b.WriteString("# Managed by alias_manager — do not edit manually\n")
	fmt.Fprintf(&b, "# Last updated: %s\n\n", now)

	keys := make([]string, 0, len(aliases))
	for k := range aliases {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		cmd := aliases[name]
		if shellType == "fish" {
			fmt.Fprintf(&b, "alias %s '%s'\n", name, shellEscape(cmd))
		} else {
			fmt.Fprintf(&b, "alias %s='%s'\n", name, shellEscape(cmd))
		}
	}

	return b.String()
}

func writeShellFile(aliases map[string]string, shellType string) {
	ensureConfigDir()

	var fileName string
	if shellType == "fish" {
		fileName = "aliases.fish"
	} else {
		fileName = "aliases.sh"
	}
	shellFilePath := filepath.Join(configDir, fileName)

	content := generateShellContent(aliases, shellType)

	if err := atomicWriteFile(shellFilePath, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to write shell file: %v\n", err)
		os.Exit(1)
	}
}

func rebuildShellFile() {
	aliases, err := loadAliases()
	if err != nil {
		return
	}

	shellType := "bash"
	cfg, err := loadConfig()
	if err == nil {
		shellType = cfg.Shell
	}

	writeShellFile(aliases, shellType)
}

func rebuildShellFileCommand() {
	ensureInitialized()

	aliases, err := loadAliases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	writeShellFile(aliases, cfg.Shell)

	shellName := "aliases.sh"
	if cfg.Shell == "fish" {
		shellName = "aliases.fish"
	}
	fmt.Printf("  ✓ Rebuilt %s with %d aliases\n", shellName, len(aliases))
}

func shellIntegrationCommand() {
	aliases, err := loadAliases()
	if err != nil {
		return
	}

	shellType := "bash"
	cfg, err := loadConfig()
	if err == nil {
		shellType = cfg.Shell
	}

	keys := make([]string, 0, len(aliases))
	for k := range aliases {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		cmd := aliases[name]
		if shellType == "fish" {
			fmt.Printf("alias %s '%s'\n", name, shellEscape(cmd))
		} else {
			fmt.Printf("alias %s='%s'\n", name, shellEscape(cmd))
		}
	}
}

func addToShellCommand() {
	ensureInitialized()

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	rcFile := cfg.ShellRC
	shellType := cfg.Shell

	if _, err := os.Stat(rcFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "  ! Shell RC file not found: %s\n", shortPath(rcFile))
	}

	var sourceLine string
	var marker string
	if shellType == "fish" {
		marker = "$HOME/.alias_manager/aliases.fish"
		sourceLine = fmt.Sprintf("if test -f %s\n    source %s\nend", marker, marker)
	} else {
		marker = "$HOME/.alias_manager/aliases.sh"
		sourceLine = fmt.Sprintf(`[ -f "%s" ] && source "%s"`, marker, marker)
	}

	data, err := os.ReadFile(rcFile)
	if err == nil && strings.Contains(string(data), ".alias_manager/aliases.") {
		fmt.Printf("  ! Already integrated into %s\n", shortPath(rcFile))
		return
	}

	f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to open RC file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "\n# Alias manager integration\n%s\n", sourceLine); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to write to RC file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("  ✓ Added to %s. Run: source %s\n", shortPath(rcFile), shortPath(rcFile))
}
