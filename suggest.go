package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type suggestion struct {
	alias   string
	command string
}

var categories = map[string][]suggestion{
	"navigation": {
		{"..", "cd .."},
		{"...", "cd ../.."},
		{"ll", "ls -la"},
		{"la", "ls -A"},
		{"lt", "ls -ltr"},
		{"home", "cd ~"},
	},
	"git": {
		{"gs", "git status"},
		{"ga", "git add"},
		{"gc", "git commit"},
		{"gp", "git push"},
		{"gl", "git pull"},
		{"gd", "git diff"},
		{"gco", "git checkout"},
		{"gb", "git branch"},
		{"glg", "git log --oneline --graph --decorate"},
		{"gca", "git commit --amend"},
		{"gst", "git stash"},
	},
	"system": {
		{"dfh", "df -h"},
		{"duh", "du -sh"},
		{"psg", "ps aux | grep"},
		{"mnt", "mount | column -t"},
		{"ports", "netstat -tulanp"},
		{"myip", "curl ifconfig.me"},
		{"cls", "clear"},
	},
	"utility": {
		{"untar", "tar -xvf"},
		{"extract", "tar -xzvf"},
		{"mcd", "mkdir -p"},
		{"path", "echo $PATH | tr ':' '\\n'"},
		{"bench", "hyperfine"},
	},
	"fun": {
		{"matrix", "cmatrix"},
		{"weather", "curl wttr.in"},
		{"joke", "curl https://icanhazdadjoke.com"},
		{"catfact", "curl https://catfact.ninja/fact"},
	},
}

func suggestCommand(args []string) {
	existing, err := loadAliases()
	if err != nil {
		existing = make(map[string]string)
	}

	category := ""
	if len(args) > 0 {
		category = strings.ToLower(args[0])
	}

	if category != "" {
		if _, ok := categories[category]; !ok {
			valid := make([]string, 0, len(categories))
			for c := range categories {
				valid = append(valid, c)
			}
			sort.Strings(valid)
			fmt.Fprintf(os.Stderr, "  ✗ Unknown category: %s\n", category)
			fmt.Fprintf(os.Stderr, "  Available categories: %s\n", strings.Join(valid, ", "))
			os.Exit(1)
		}
	}

	type displayRow struct {
		alias   string
		command string
		status  string
	}
	var rows []displayRow

	cats := make([]string, 0, len(categories))
	for c := range categories {
		cats = append(cats, c)
	}
	sort.Strings(cats)

	for _, c := range cats {
		if category != "" && c != category {
			continue
		}
		for _, s := range categories[c] {
			status := "NEW"
			if _, exists := existing[s.alias]; exists {
				status = "EXISTS"
			}
			rows = append(rows, displayRow{s.alias, s.command, status})
		}
	}

	if len(rows) == 0 {
		fmt.Println("  No suggestions available.")
		return
	}

	maxAliasLen := len("Alias")
	maxCmdLen := len("Command")
	for _, r := range rows {
		if len(r.alias) > maxAliasLen {
			maxAliasLen = len(r.alias)
		}
		if len(r.command) > maxCmdLen {
			maxCmdLen = len(r.command)
		}
	}
	if maxAliasLen > 30 {
		maxAliasLen = 30
	}
	if maxCmdLen > 70 {
		maxCmdLen = 70
	}

	totalWidth := maxAliasLen + maxCmdLen + 10
	if totalWidth > 114 {
		totalWidth = 114
	}

	fmt.Printf("  %-*s  %-*s  Status\n", maxAliasLen, "Alias", maxCmdLen, "Command")
	fmt.Println("  " + strings.Repeat("─", totalWidth))

	for _, r := range rows {
		displayName := r.alias
		displayCmd := r.command
		if len(displayName) > 30 {
			displayName = displayName[:27] + "..."
		}
		if len(displayCmd) > 70 {
			displayCmd = displayCmd[:67] + "..."
		}
		fmt.Printf("  %-*s  %-*s  %s\n", maxAliasLen, displayName, maxCmdLen, displayCmd, r.status)
	}
}
