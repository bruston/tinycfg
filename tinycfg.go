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

type Config struct {
	vals map[string]string
}

func (c Config) Get(key string) string {
	return c.vals[key]
}

func (c Config) Set(key, value string) error {
	if key == "" {
		return errors.New("key cannot be blank")
	}
	if value == "" {
		return errors.New("value cannot be blank")
	}
	if strings.Contains(key, delim) {
		return fmt.Errorf("key cannot contain '%s'", delim)
	}
	if strings.Contains(value, "\n") {
		return errors.New("value cannot contain newlines")
	}
	c.vals[key] = value
	return nil
}

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

func New() Config {
	return Config{make(map[string]string)}
}

func Open(path string, required []string) (Config, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, nil, err
	}
	defer file.Close()
	return Decode(file, required)
}

func Decode(r io.Reader, required []string) (Config, []string, error) {
	cfg := Config{make(map[string]string)}
	scanner := bufio.NewScanner(r)
	var line string
	var lineNum int
	for scanner.Scan() {
		lineNum++
		line = strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, commentPrefix) {
			continue
		}
		args := strings.SplitN(line, delim, 2)
		if args[0] == "" || args[1] == "" {
			return Config{}, nil, fmt.Errorf("no key/value pair found at line %d", lineNum)
		}
		if _, ok := cfg.vals[args[0]]; ok {
			return Config{}, nil, fmt.Errorf("duplicate entry for key %s at line %d", args[0], lineNum)
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
			return Config{}, missing, fmt.Errorf("missing required fields: %s", strings.Join(missing, " "))
		}
	}
	return cfg, nil, nil
}
