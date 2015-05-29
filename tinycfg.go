package tinycfg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
)

const (
	delim         = "="
	commentPrefix = "//"
)

// A Config stores key, value pairs.
type Config struct {
	mu   sync.RWMutex
	vals map[string]string
}

// Get returns the value for a specified key or an empty string if the key was not found.
func (c *Config) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.vals[key]
}

// Set adds a key, value pair or modifies an existing one. The returned error can be safely
// ignored if you are certain that both the key and value are valid. Keys are invalid if
// they contain '=', newline characters or are empty. Values are invalid if they contain
// newline characters or are empty.
func (c *Config) Set(key, value string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if value == "" {
		return errors.New("value cannot be empty")
	}
	if strings.Contains(key, delim) {
		return fmt.Errorf("key cannot contain '%s'", delim)
	}
	if strings.Contains(value, "\n") {
		return errors.New("value cannot contain newlines")
	}
	if strings.Contains(key, "\n") {
		return errors.New("key cannot contain newlines")
	}
	c.mu.Lock()
	c.vals[key] = value
	c.mu.Unlock()
	return nil
}

// Delete removes a key, value pair.
func (c *Config) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.vals, key)
}

// Encode writes out a Config instance in the correct format to a Writer. Key, value pairs
// are listed in alphabetical order.
func (c *Config) Encode(w io.Writer) error {
	var lines []string
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k, v := range c.vals {
		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Sort(sort.StringSlice(lines))
	for _, v := range lines {
		_, err := fmt.Fprintln(w, v)
		if err != nil {
			return fmt.Errorf("unable to encode line: %s\n%s", v, err)
		}
	}
	return nil
}

// New returns an empty Config instance ready for use.
func New() *Config {
	return &Config{vals: make(map[string]string)}
}

// Open is a convenience function that opens a file at a specified path, passes it to Decode
// then closes the file.
func Open(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Decode(file)
}

// Decode creates a new Config instance from a Reader.
func Decode(r io.Reader) (*Config, error) {
	cfg := &Config{vals: make(map[string]string)}
	scanner := bufio.NewScanner(r)
	for lineNum := 1; scanner.Scan(); lineNum++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, commentPrefix) {
			continue
		}
		args := strings.SplitN(line, delim, 2)
		key, value := strings.TrimSpace(args[0]), strings.TrimSpace(args[1])
		if key == "" || value == "" {
			return cfg, fmt.Errorf("no key/value pair found at line %d", lineNum)
		}
		if _, ok := cfg.vals[key]; ok {
			return cfg, fmt.Errorf("duplicate entry for key %s at line %d", key, lineNum)
		}
		cfg.vals[key] = value
	}
	if scanner.Err() != nil {
		return cfg, scanner.Err()
	}
	return cfg, nil
}

// DecodeWithDefaults creates a new Config instance and allows a map of defaults to be provided.
// After decoding, the default key/value pairs are set if not already present.
func DecodeWithDefaults(r io.Reader, defaults map[string]string) (*Config, error) {
	cfg, err := Decode(r)
	if err != nil {
		return cfg, err
	}
	for k, v := range defaults {
		if cfg.Get(k) == "" {
			if err := cfg.Set(k, v); err != nil {
				return cfg, err
			}
		}
	}
	return cfg, nil
}

// Missing checks for the existence of a slice of keys in a Config instance and returns a slice
// which contains keys that are missing, or nil if there are no missing keys.
func Missing(cfg *Config, required []string) []string {
	var missing []string
	for _, k := range required {
		if v := cfg.Get(k); v == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return missing
	}
	return nil
}

// NewFromEnv returns a new Config instance populated from environment variables.
func NewFromEnv(keys []string) (*Config, error) {
	var buf bytes.Buffer
	for _, k := range keys {
		fmt.Fprintln(&buf, k, "=", os.Getenv(k))
	}
	cfg, err := Decode(&buf)
	return cfg, err
}
