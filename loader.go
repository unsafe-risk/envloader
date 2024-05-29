package envloader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// LoadEnvReader loads environment variables from a reader (e.g., a file) and sets them as environment variables.
func LoadEnvReader(r io.Reader) error {
	// Read the file line by line.
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore comments and empty lines.
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Split the line into key and value.
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid .env file format: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set the environment variable.
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	return nil
}

// LoadEnvFile loads environment variables from a file at the given path.
func LoadEnvFile(path string) error {
	// Open the file at the given path.
	f, err := os.Open(path)
	if err != nil {
		// Return the error if there was an issue opening the file.
		return err
	}
	// Ensure the file is closed when the function exits.
	defer f.Close()
	// Load environment variables from the file reader.
	return LoadEnvReader(f)
}

// LoadAndBindEnvReader loads environment variables from a reader and binds them to a struct.
func LoadAndBindEnvReader(r io.Reader, dst interface{}) error {
	// Load environment variables from the reader.
	err := LoadEnvReader(r)
	if err != nil {
		// Return the error if there was an issue loading environment variables.
		return err
	}

	// Bind environment variables to the struct.
	return BindStruct(dst, func(key string) (string, error) {
		if value, ok := os.LookupEnv(key); ok {
			return value, nil
		}
		return "", fmt.Errorf("environment variable %s not found", key)
	})
}

// LoadAndBindEnvFile loads environment variables from a file at the given path and binds them to a struct.
func LoadAndBindEnvFile(path string, dst interface{}) error {
	// Load environment variables from the file.
	err := LoadEnvFile(path)
	if err != nil {
		// Return the error if there was an issue loading environment variables.
		return err
	}

	// Bind environment variables to the struct.
	return BindStruct(dst, func(key string) (string, error) {
		if value, ok := os.LookupEnv(key); ok {
			return value, nil
		}
		return "", fmt.Errorf("environment variable %s not found", key)
	})
}
