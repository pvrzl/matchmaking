package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	sqlDir = "sql"
)

type Command string

const (
	Create Command = "create"
	Up     Command = "up"
	Down   Command = "down"
	Fix    Command = "fix"
)

var commandFlag Command

func (c *Command) String() string {
	return string(*c)
}

func (c *Command) Set(value string) error {
	switch strings.ToLower(value) {
	case string(Create), string(Up), string(Down), string(Fix):
		*c = Command(value)
		return nil
	default:
		return errors.New("must be one of: create, up, down and fix")
	}
}

func main() {
	flag.Var(&commandFlag, "command", "available command: create, up, down and fix")
	flag.Parse()

	switch commandFlag {
	case Create:
		create()
	case Up:
		up()
	case Down:
		down()
	case Fix:
		fix()
	default:
		fmt.Fprintln(os.Stderr, "error: the -command flag is required")
		os.Exit(1)
	}
}
