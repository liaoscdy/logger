package logger

import (
	"fmt"
	"os"
)

type ConsoleExecutor struct{}

func NewConsoleExecutor() *ConsoleExecutor {
	return &ConsoleExecutor{}
}

func (c *ConsoleExecutor) WriteMsg(msg []byte) error {
	_, err := fmt.Fprintln(os.Stdout, string(msg))
	return err
}

func (c *ConsoleExecutor) Flush() {}

func (c *ConsoleExecutor) Close() {}
