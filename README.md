tinycfg
=======

**Under development and likely to change.**

Package tinycfg provides minimal configuration file support using a simple key/value format.

# Installation

```
$ go get github.com/bruston/tinycfg
```

# Godoc

https://godoc.org/github.com/bruston/tinycfg

# Format

```
// This is a comment.
server=irc.example.com
port=6667
channel=#example
```

# Decoding

```go
package main

import (
	"fmt"
	"log"

	"github.com/bruston/tinycfg"
)

var required = []string{"server", "port", "channel"}

func main() {
	cfg, missing, err := tinycfg.Open("example.cfg", required)
	if err != nil {
		if missing != nil {
			log.Fatalf("missing required config keys: %v", missing)
		}
		log.Fatalf("unable to decode config file: %s", err)
	}

	fmt.Println(cfg.Get("server"), cfg.Get("port"), cfg.Get("channel"))
	// irc.example.com 6667 #example
}
```

# Encoding

```go
package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/bruston/tinycfg"
)

func main() {
	cfg := tinycfg.New()
	cfg.Set("server", "irc.example.com")
	cfg.Set("port", "6667")
	cfg.Set("channel", "#example")

	var buf bytes.Buffer
	if err := cfg.Encode(&buf); err != nil {
		log.Fatalf("error writing to config: %s", err)
	}

	fmt.Print(buf.String())
	// channel=#example
	// port=6667
	// server=irc.example.com
}
```
