# Zerolog Wrapper
[rs/zerolog][1] is a logger package for [Go][2] which follows [uber/zap][3]'s JSON output format with calls chained for ease of logging.

This package is just a wrapper on initialization of zerolog logger instance which stores and propagates through [context.Context][4].

## Install
```
go get -u github.com/Bhargav-InfraCloud/zerolog-wrapper@latest
```

## Usage
#### Example-1: Initialization and propagation of logger through context variable
```go
package main

import (
	"context"
	"os"

	zlog "github.com/Bhargav-InfraCloud/zerolog-wrapper"
)

func main() {
	ctx, logger := zlog.NewLogger(context.Background(), os.Stdout, zlog.LevelDebug)

	logger.Info().Msg("Message from main")

	greet(ctx)
}

func greet(ctx context.Context) {
	logger := zlog.FromContext(ctx)

	logger.Info().Msg("Message from greet")
}

```

This outputs logs as follows:
```json
{"log-level":"info","timestamp":"2023-03-10T15:32:07+05:30","caller":"main.go:13","log-message":"Message from main"}
{"log-level":"info","timestamp":"2023-03-10T15:32:07+05:30","caller":"main.go:21","log-message":"Message from greet"}
```
#### Example-2: Using zerolog's `With()` to create child logger with added extra fields.
In following example, `With` adds `{"flow": "example-manager"}` to child logger which is then stored in `Manager` instance.

```go
package main

import (
	"context"
	"os"

	zlog "github.com/Bhargav-InfraCloud/zerolog-wrapper"
)

func main() {
	ctx, logger := zlog.NewLogger(context.Background(), os.Stdout, zlog.LevelDebug)

	logger.Info().Msg("Message from main")

	manager := NewManager(ctx)
	manager.Greet()
}

func (m Manager) Greet() {
	m.logger.Info().Msg("Message from Greet")
}

type Manager struct {
	logger zlog.Logger
}

func NewManager(ctx context.Context) *Manager {
	logger := zlog.FromContext(ctx).With().
		Str("flow", "example-manager").Logger()

	return &Manager{
		logger: zlog.FromRawLogger(logger),
	}
}
```

This outputs logs as follows:
```json
{"log-level":"info","timestamp":"2023-03-10T15:41:33+05:30","caller":"main.go:13","log-message":"Message from main"}
{"log-level":"info","flow":"example-manager","timestamp":"2023-03-10T15:41:33+05:30","caller":"main.go:20","log-message":"Message from greet"}
```

> Note: The `zlog.FromRawLogger` creates `zlog.Logger` from `zerolog.Logger` avoiding the need for direct import of [zerolog][1].

[1]: https://github.com/rs/zerolog
[2]: https://go.dev/
[3]: https://github.com/uber-go/zap
[4]: https://pkg.go.dev/context#Context