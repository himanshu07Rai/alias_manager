package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ShellRC   string `json:"shell_rc"`
	AliasFile string `json:"alias_file"`
	Shell     string `json:"shell"`
}

var (
	configDir  string
	configPath string
	aliasFile  string
)

func initPaths() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "  ✗ Cannot determine home directory")
		os.Exit(1)
	}
	configDir = filepath.Join(home, ".alias_manager")
	configPath = filepath.Join(configDir, "config.json")
	aliasFile = filepath.Join(configDir, "aliases.json")
}

func ensureConfigDir() {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to create config directory: %v\n", err)
		os.Exit(1)
	}
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("corrupted config file. Try deleting %s", configPath)
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to marshal config: %v\n", err)
		os.Exit(1)
	}
	if err := atomicWriteFile(configPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to save config: %v\n", err)
		os.Exit(1)
	}
}

func ensureInitialized() {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "  ✗ Run 'alias_manager init' first")
		os.Exit(1)
	}
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "bash"
	}
	base := filepath.Base(shell)
	switch base {
	case "zsh":
		return "zsh"
	case "fish":
		return "fish"
	default:
		return "bash"
	}
}

func defaultRCForShell(shell string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	switch shell {
	case "zsh":
		return filepath.Join(home, ".zshrc")
	case "fish":
		return filepath.Join(home, ".config", "fish", "config.fish")
	default:
		return filepath.Join(home, ".bashrc")
	}
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[1:])
		}
	}
	return path
}

func initCommand() {
	ensureConfigDir()

	detectedShell := detectShell()
	defaultRC := defaultRCForShell(detectedShell)

	fmt.Printf("  Detected shell: %s\n", detectedShell)
	fmt.Printf("  Shell RC file: %s\n", shortPath(defaultRC))

	resp := promptInput("  Is this correct? [Y/n]: ")
	if strings.ToLower(resp) == "n" || strings.ToLower(resp) == "no" {
		customRC := promptInput("  Enter shell RC path: ")
		if customRC == "" {
			fmt.Fprintln(os.Stderr, "  ✗ RC file path cannot be empty")
			os.Exit(1)
		}
		defaultRC = expandPath(customRC)

		customShell := promptInput("  Enter shell type (bash/zsh/fish): ")
		if customShell != "" {
			detectedShell = strings.ToLower(customShell)
		}
	}

	cfg := &Config{
		ShellRC:   defaultRC,
		AliasFile: aliasFile,
		Shell:     detectedShell,
	}
	saveConfig(cfg)

	fmt.Println("  ✓ Config saved.")
	fmt.Println("  → Run 'alias_manager add-to-shell' to register with your shell.")
}
