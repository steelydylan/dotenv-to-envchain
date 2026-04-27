package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: env-to-envchain <namespace> <envfile>")
		fmt.Fprintln(os.Stderr, "  Example: env-to-envchain myapp .env")
		os.Exit(1)
	}

	namespace := os.Args[1]
	envFile := os.Args[2]

	entries, err := parseEnvFile(envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", envFile, err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Fprintln(os.Stderr, "No environment variables found in the file.")
		os.Exit(1)
	}

	for _, e := range entries {
		if err := setEnvchain(namespace, e.key, e.value); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting %s: %v\n", e.key, err)
			os.Exit(1)
		}
		fmt.Printf("Set %s in namespace %q\n", e.key, namespace)
	}

	fmt.Printf("\nDone! %d variable(s) saved to namespace %q.\n", len(entries), namespace)
	fmt.Printf("Run: envchain %s env\n", namespace)
}

type envEntry struct {
	key   string
	value string
}

func parseEnvFile(path string) ([]envEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []envEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Skip export prefix
		line = strings.TrimPrefix(line, "export ")

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// Remove surrounding quotes
		value = unquote(value)

		entries = append(entries, envEntry{key: key, value: value})
	}

	return entries, scanner.Err()
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func setEnvchain(namespace, key, value string) error {
	cmd := exec.Command("envchain", "--set", namespace, key)
	cmd.Stdin = strings.NewReader(value + "\n")
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
