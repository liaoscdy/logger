# About

A simple, high-performance, multi-threaded friendly golang logging library

This library provides convenient logging capabilities to stdout and logging to files.

It also provides simple and convenient custom extension capabilities.

# Example

Print logs to standard output / console

```go
package main

import (
	"github.com/liaoscdy/logger"
)

func main() {
	log := logger.NewDefaultLogger()
	log.SetCallDepth(2)
	log.SetLevel(logger.LevelDebug)
	log.Start()
	defer log.Close()

	// ------------------ print log msg
	log.Info("info Msg, value:%d", 1)
	log.Error("error Msg without any format value")
	log.Warn("warn Msg, Value1:%d, Value2:%s", 100, "stringByFmt")
}
```

Print logs to a file

```go
package main

import (
	"github.com/liaoscdy/logger"
)

func main() {
	log := logger.NewLogger()
	// logging to file
	fileExecutor := logger.NewFileExecutor("/tmp/logger_test.log")
	fileExecutor.EnableFileRotate()
	fileExecutor.SetRotateMaxDays(7)
	log.AppendExecutor(fileExecutor)
	log.SetCallDepth(2)
	log.SetLevel(logger.LevelDebug)
	log.Start()
	defer log.Close()

	// ------------------ print log msg
	log.Info("info Msg, value:%d", 1)
	log.Error("error Msg without any format value")
	log.Warn("warn Msg, Value1:%d, Value2:%s", 100, "stringByFmt")
}
```

Print logs to files and standard output/console at the same time

```go
package main

import (
	"github.com/liaoscdy/logger"
)

func main() {
	log := logger.NewLogger()
	log.AppendExecutor(logger.NewConsoleExecutor())
	log.AppendExecutor(logger.NewFileExecutor("/tmp/logger_test.log"))
	log.SetCallDepth(2)
	log.SetLevel(logger.LevelDebug)
	log.Start()
	defer log.Close()

	// ------------------ print log msg
	log.Info("info Msg, value:%d", 1)
	log.Error("error Msg without any format value")
	log.Warn("warn Msg, Value1:%d, Value2:%s", 100, "stringByFmt")
}
```

# Customized

This log library supports custom log formatting and custom log output methods

## Formatter

Implement Formatter to format the log into any format you want.

The definitions of the Formatter interface and LogMsg structure in this library are as follows:

```go
type LogMsg struct {
    Level     LogLevel
    Timestamp time.Time
    Msg       string
    FileName  string
    FileLine  int
}

type Formatter interface {
    FormatMsg(msg *LogMsg) []byte
}
```

Here is a simple code example

```go
import "github.com/liaoscdy/logger"

type CustomerFormatter struct{}

func (f *CustomerFormatter) FormatMsg(msg *logger.LogMsg) []byte {
	// example
	// ......
	return []byte(msg.Msg)
}

func main() {
	log := logger.NewLogger()
	log.SetLevel(logger.LevelInfo)
	log.SetCallDepth(2)
	log.AppendExecutor(logger.NewConsoleExecutor())
	log.SetFormatter(&CustomerFormatter{})
	log.Start()
	defer log.Close()

	log.Info("use the customer formatter")
}
```

## Executor

Executor is responsible for outputting the formatted log information to where you want it.

Implement the Executor interface to customize your output method.

```go
type Executor interface {
	WriteMsg(msg []byte) error
	Flush()
	Close()
}
```

Here is a simple code example

```go
import (
	"github.com/liaoscdy/logger"
)

type CustomerExecutor struct{}

func (f *CustomerExecutor) WriteMsg(msg []byte) error {
	// ... your code here
	return nil
}

func (f *CustomerExecutor) Flush() {
	// ... your code here
	// ... responsible for persisting cached data
}

func (f *CustomerExecutor) Close() {
	// ... your code here
	// It is about to end and resources should be cleaned up
}

func main() {
	log := logger.NewLogger()
	log.SetLevel(logger.LevelInfo)
	log.SetCallDepth(2)
	log.AppendExecutor(&CustomerExecutor{})
	log.Start()
	defer log.Close()

	log.Info("use the customer formatter")
}
```