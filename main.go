package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var stdinReader = bufio.NewReader(os.Stdin)

func main() {
	initPaths()

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "init":
		initCommand()
	case "list":
		ensureInitialized()
		pattern := ""
		if len(os.Args) > 2 {
			pattern = os.Args[2]
		}
		listAliases(pattern)
	case "create":
		ensureInitialized()
		createCommand(os.Args[2:])
	case "delete":
		ensureInitialized()
		deleteCommand(os.Args[2:])
	case "edit":
		ensureInitialized()
		editCommand(os.Args[2:])
	case "import":
		ensureInitialized()
		importCommand(os.Args[2:])
	case "export":
		ensureInitialized()
		exportCommand(os.Args[2:])
	case "suggest":
		suggestCommand(os.Args[2:])
	case "add-to-shell":
		addToShellCommand()
	case "rebuild-shell-file":
		rebuildShellFileCommand()
	case "shell-integration":
		shellIntegrationCommand()
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "  ✗ Unknown command: %s\n", os.Args[1])
		fmt.Fprintln(os.Stderr, "  → Run 'alias_manager help' for usage")
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`Usage: alias_manager <command> [options]

Commands:
  init                    Set up alias_manager for your shell
  list [pattern]          List aliases, optionally filter by pattern
  create <name> <cmd>     Create a new alias (use --force to overwrite)
  delete <name>           Delete an alias (use --yes/-y to skip prompt)
  edit <name> [cmd]       Edit an existing alias
  import <file>           Import aliases from JSON (use --overwrite to replace)
  export <file>           Export aliases to JSON file
  suggest [category]      Show alias suggestions by category
  add-to-shell            Add shell integration hook to RC file
  rebuild-shell-file      Rebuild shell aliases file from JSON storage
  shell-integration       Print aliases in shell-compatible format
  help                    Show this help message

Categories for suggest: navigation, git, system, utility, fun
`)
}

func createCommand(args []string) {
	force := false
	var cmdParts []string
	var name string
	nameSet := false

	for _, a := range args {
		if a == "--force" {
			force = true
		} else if !nameSet {
			name = a
			nameSet = true
		} else {
			cmdParts = append(cmdParts, a)
		}
	}

	if (!nameSet || len(cmdParts) == 0) && isInteractive() {
		tuiName, tuiCmd, ok := runCreateTUI()
		if !ok {
			fmt.Println("  Cancelled.")
			return
		}
		createAlias(tuiName, tuiCmd, force)
		return
	}

	if !nameSet || len(cmdParts) == 0 {
		fmt.Fprintln(os.Stderr, "  ✗ Usage: alias_manager create <name> <command> [--force]")
		os.Exit(1)
	}

	command := strings.Join(cmdParts, " ")
	createAlias(name, command, force)
}

func deleteCommand(args []string) {
	var name string
	yes := false
	for _, a := range args {
		if a == "--yes" || a == "-y" {
			yes = true
		} else if name == "" {
			name = a
		}
	}

	if name == "" {
		fmt.Fprintln(os.Stderr, "  ✗ Usage: alias_manager delete <name> [--yes]")
		os.Exit(1)
	}

	deleteAlias(name, yes)
}

func editCommand(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "  ✗ Usage: alias_manager edit <name> [new command]")
		os.Exit(1)
	}

	name := args[0]
	if len(args) >= 2 {
		newCommand := strings.Join(args[1:], " ")
		editAlias(name, &newCommand)
		return
	}

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

	if isInteractive() {
		newCmd, ok := runEditTUI(name, oldCmd)
		if !ok {
			fmt.Println("  Cancelled.")
			return
		}
		editAlias(name, &newCmd)
		return
	}

	editAlias(name, nil)
}

func importCommand(args []string) {
	var filePath string
	overwrite := false
	for _, a := range args {
		if a == "--overwrite" {
			overwrite = true
		} else if filePath == "" {
			filePath = a
		}
	}

	if filePath == "" {
		fmt.Fprintln(os.Stderr, "  ✗ Usage: alias_manager import <file> [--overwrite]")
		os.Exit(1)
	}

	importAliases(filePath, overwrite)
}

func exportCommand(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "  ✗ Usage: alias_manager export <file>")
		os.Exit(1)
	}

	exportAliases(args[0])
}

func promptInput(prompt string) string {
	fmt.Print(prompt)
	input, err := stdinReader.ReadString('\n')
	if err != nil {
		fmt.Println()
		os.Exit(0)
	}
	return strings.TrimSpace(input)
}

func promptConfirm(prompt string) bool {
	input := promptInput(prompt)
	return strings.ToLower(input) == "y" || strings.ToLower(input) == "yes"
}

func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, perm); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

func shortPath(p string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	return strings.Replace(p, home, "~", 1)
}
