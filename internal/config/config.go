package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadConfig loads configuration files in order, with higher-priority files first.
// Environment variables already set in the process are NOT overwritten.
func LoadConfig() error {
	paths, err := getConfigFilePaths()
	if err != nil {
		return fmt.Errorf("failed to get config file paths: %w", err)
	}

	for _, path := range paths {
		values, err := readEnvFile(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to load config file %s: %w", path, err)
		}
		for k, v := range values {
			if os.Getenv(k) == "" {
				os.Setenv(k, v) //nolint:errcheck
			}
		}
	}

	return nil
}

// readEnvFile parses a KEY=VALUE file and returns its entries.
func readEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		values[strings.TrimSpace(k)] = strings.Trim(strings.TrimSpace(v), `"'`)
	}
	return values, scanner.Err()
}

// GetDesktopEntryDir resolves the system path for the desktop entry file.
func GetDesktopEntryDir() string {
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		return filepath.Join(os.Getenv("HOME"), ".local", "share")
	}
	return xdgDataHome
}

// GetDistroboxPath returns the path to the current distrobox executable.
func GetDistroboxPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get distrobox executable path: %w", err)
	}
	return path, nil
}

// getConfigFilePaths returns configuration file paths ordered highest-priority first.
func getConfigFilePaths() ([]string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate symlinks for executable path: %w", err)
	}

	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return []string{
		filepath.Join(os.Getenv("HOME"), ".distroboxrc"),
		filepath.Join(xdgConfigHome, "distrobox", "distrobox.conf"),
		"/etc/distrobox/distrobox.conf",
		"/usr/local/share/distrobox/distrobox.conf",
		"/usr/etc/distrobox/distrobox.conf",
		"/usr/share/defaults/distrobox/distrobox.conf",
		"/usr/share/distrobox/distrobox.conf",
		filepath.Join(filepath.Dir(execPath), "..", "share", "distrobox", "distrobox.conf"),
	}, nil
}

