package utils

import (
	"bufio"
	"os"
	"strings"
)

func LoadEnv(envPaths ...string) error {
	envPath := ".env"
	if len(envPaths) > 0 {
		envPath = envPaths[0]
	}
	file, err := os.Open(envPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) < 2 {
			continue
		}
		key := parts[0]
		values := parts[1:]
		value := strings.Join(values, "=")
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
