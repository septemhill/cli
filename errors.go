package cli

import (
	"errors"
	"fmt"
)

var (
	// ErrCommandExisted = errors.New("command already existed")
	ErrNoCommand = errors.New("no command provided")
)

func ErrUnknownFlag(flag string) error {
	return fmt.Errorf("unknown flag: %s", flag)
}

func ErrUnknownCommand(cmd string) error {
	return fmt.Errorf("unknown command: %s", cmd)
}

func ErrCommandExisted(cmd string) error {
	return fmt.Errorf("command already existed: %s", cmd)
}
