package tinycfg

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

const (
	delim         = "="
	commentPrefix = "//"
)

// A Config stores key, value pairs.
type Config struct {
	vals map[string]string
}

// Get returns the value for a specified key or an empty string if the key was not found.
func (c Config) Get(key string) string {
	return c.vals[key]
}

// Set adds a key, value pair or modifies an existing one. The returned error can be safely
// ignored if you are certain that both the key and value are valid. Keys are invalid if
// they contain '=', newline characters or are empty. Values are invalid if they contain
// newline characters or are empty.
func (c Config) Set(key, value string) error {
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
	c.vals[key] = value
	return nil
}

// Delete removes a key, value pair.
func (c Config) Delete(key string) {
	delete(c.vals, key)
}

func (c Config) Encode(w io.Writer) error {
	var lines []string
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
func New() Config {
	return Config{make(map[string]string)}
}

// Open is a convenience function that opens a file at a specified path, passes it to Encode
// then closes the file.
func Open(path string, required []string) (Config, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, nil, err
	}
	defer file.Close()
	return Decode(file, required)
}

// Decode creates a new Config instance from a Reader. Required keys can be specified by passing
// in a string slice, or nil if there are no required keys. If there are missing required keys
// they are returned in a string slice along with an error.
func Decode(r io.Reader, required []string) (Config, []string, error) {
	cfg := Config{make(map[string]string)}
	scanner := bufio.NewScanner(r)
	for lineNum := 1; scanner.Scan(); lineNum++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, commentPrefix) {
			continue
		}
		args := strings.SplitN(line, delim, 2)
		if args[0] == "" || args[1] == "" {
			return cfg, nil, fmt.Errorf("no key/value pair found at line %d", lineNum)
		}
		if _, ok := cfg.vals[args[0]]; ok {
			return cfg, nil, fmt.Errorf("duplicate entry for key %s at line %d", args[0], lineNum)
		}
		cfg.vals[strings.TrimSpace(args[0])] = strings.TrimSpace(args[1])
	}
	if scanner.Err() != nil {
		return cfg, nil, scanner.Err()
	}
	if required != nil {
		var missing []string
		for _, v := range required {
			if val := cfg.Get(v); val == "" {
				missing = append(missing, v)
			}
		}
		if len(missing) > 0 {
			return cfg, missing, fmt.Errorf("missing required keys")
		}
	}
	return cfg, nil, nil
}
