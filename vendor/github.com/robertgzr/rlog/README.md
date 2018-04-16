rlog
====

_rust log or robert's log if you like ;)_


Yet another logging library written in Go...

Mostly written for my own projects, with inspiration taken from https://github.com/inconshreveable/log15 and https://crates.io/crates/env_logger.

![preview](static/preview.png)

## Usage

Really simple, just import the package and use it:

```
package main

import "github.com/robertgzr/rlog"

func main() {
    rlog.Info("hello world!")
}
```

The log level can be set via an environment variable by including `rlog.ParseEnv()` somewhere near the start of the program. By default `GO_LOG` is parsed.

To create a new logger with some context:

```
rl := rlog.With("ctxkey", "ctxvalue")
rl.Debug("logging with context")

$ DBUG| logging with context (ctxkey="ctxvalue")
```
