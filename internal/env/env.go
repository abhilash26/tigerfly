package env

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var envMap = make(map[string]string)

func processEnvFile(file *os.File) error {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		processLine(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
		return err
	}
	return nil
}

func processLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") { // Skip empty lines and comments
		return
	}

	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		log.Printf("Skipping invalid line: %s", line)
		return
	}

	key := strings.TrimSpace(parts[0])
	value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

	envMap[key] = value
}

func LoadEnvFile(filepath string) error {
	if filepath == "" {
		filepath = ".env"
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return err
	}
	defer file.Close()

	return processEnvFile(file)
}

func getEnv[T any](key string, defaultValue T, parser func(string) (T, error)) T {
	valueStr := GetString(key, "")
	if valueStr == "" {
		return defaultValue
	}
	parsedValue, err := parser(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s: %v. Using default: %v", key, err, defaultValue)
		return defaultValue
	}
	return parsedValue
}

// Get function for different data types
func GetString(key string, defaultValue string) string {
	if value, exists := envMap[key]; exists && value != "" {
		return value
	}
	log.Printf("Environment variable %s not set or empty, using default: %s", key, defaultValue)
	return defaultValue
}

func GetInt(key string, defaultValue int) int {
	return getEnv(key, defaultValue, strconv.Atoi)
}

func GetFloat(key string, defaultValue float64) float64 {
	return getEnv(key, defaultValue, func(value string) (float64, error) {
		return strconv.ParseFloat(value, 64)
	})
}

func GetBool(key string, defaultValue bool) bool {
	return getEnv(key, defaultValue, strconv.ParseBool)
}

func GetDuration(key string, defaultValue time.Duration) time.Duration {
	return getEnv(key, defaultValue, time.ParseDuration)
}

func GetSlice(key string, defaultValue string, separator ...string) []string {
	defaultSeparator := ","
	if len(separator) > 0 && separator[0] != "" {
		defaultSeparator = separator[0]
	}

	value := GetString(key, defaultValue)
	if value == "" {
		return nil
	}
	return strings.Split(value, defaultSeparator)
}
