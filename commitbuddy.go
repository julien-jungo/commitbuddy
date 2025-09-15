package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

var validTypes = []string{
	"feat", "fix", "docs", "style", "refactor", "perf", "test", "build", "ci", "chore", "revert",
}

var (
	typeFlag  = flag.String("t", "", "Commit type")
	scopeFlag = flag.String("s", "", "Commit scope")
	msgFlag   = flag.String("m", "", "Commit message")
	useConfig = flag.Bool("c", false, "Use config values for commit message (type and scope)")
)

func isValidType(t string) bool {
	return slices.Contains(validTypes, t)
}

func promptInput(prompt string) string {
	fmt.Print(prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return strings.TrimSpace(scanner.Text())
}

func readConfig() (string, string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", ""
	}

	configPath := filepath.Join(home, ".commitbuddy.config")
	f, err := os.Open(configPath)
	if err != nil {
		return "", ""
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	var ctype, scope string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		case strings.HasPrefix(line, "type="):
			ctype = strings.TrimPrefix(line, "type=")
		case strings.HasPrefix(line, "scope="):
			scope = strings.TrimPrefix(line, "scope=")
		}
	}

	return ctype, scope
}

func writeConfig(ctype, scope string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".commitbuddy.config")
	f, err := os.Create(configPath)
	if err != nil {
		return err
	}

	defer f.Close()

	if ctype != "" {
		fmt.Fprintf(f, "type=%s\n", ctype)
	}
	if scope != "" {
		fmt.Fprintf(f, "scope=%s\n", scope)
	}

	return nil
}

func runConfigCommand(args []string) {
	configFlagSet := flag.NewFlagSet("config", flag.ExitOnError)

	ctype := configFlagSet.String("t", "", "Default commit type")
	scope := configFlagSet.String("s", "", "Default commit scope")

	configFlagSet.Parse(args)

	if *ctype == "" && *scope == "" {
		fmt.Println("Usage: commitbuddy config -t <type> -s <scope>")
		os.Exit(1)
	}

	err := writeConfig(*ctype, *scope)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Config updated!")
}

func runCommitCommand() {
	flag.Parse()

	confType, confScope := readConfig()

	ctype := *typeFlag
	scope := *scopeFlag
	commitMsg := *msgFlag

	if *useConfig {
		if *typeFlag != "" || *scopeFlag != "" {
			fmt.Fprintln(os.Stderr, "Error: Cannot use -t or -s flags with -c.")
			os.Exit(1)
		}

		if confType == "" && confScope == "" {
			fmt.Fprintln(os.Stderr, "Error: No config found.")
			os.Exit(1)
		}

		ctype = confType
		scope = confScope
	}

	if ctype == "" {
		ctype = promptInput("Commit type: ")
	}
	for !isValidType(ctype) {
		fmt.Printf("Invalid commit type '%s'. Valid types: %s\n", ctype, strings.Join(validTypes, ", "))
		ctype = promptInput("Commit type: ")
	}

	if scope == "" {
		scope = promptInput("Commit scope: ")
	}

	for commitMsg == "" {
		commitMsg = promptInput("Commit message: ")
	}

	var fullMsg string

	switch scope {
	case "":
		fullMsg = fmt.Sprintf("%s: %s", ctype, commitMsg)
	default:
		fullMsg = fmt.Sprintf("%s(%s): %s", ctype, scope, commitMsg)
	}

	cmd := exec.Command("git", "commit", "-m", fullMsg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running git commit: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "config" {
		runConfigCommand(os.Args[2:])
		return
	}

	runCommitCommand()
}
