tinycfg
=======

Package tinycfg provides minimal configuration file support using a simple key/value format.

# Installation

```
$ go get github.com/bruston/tinycfg
```

# Godoc

https://godoc.org/github.com/bruston/tinycfg

# Limitations

`tinycfg` is minimalistic, likely *too* minimalistic for many projects. As a result of this, there are some limitations that need to be pointed out:

- keys and values are strings
- keys cannot contain '=' or newlines
- values cannot contain newlines or be empty (but may contain '=')

If you require more than just strings, look to the [strconv](http://golang.org/pkg/strconv/) package for help converting to other types. Or use one of the other configuration packages out there.

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

func main() {
	cfg, err := tinycfg.Open("example.cfg")
	if err != nil {
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
