package cli

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"

	"golang.org/x/exp/maps"
)

type Cli struct {
	cmds   map[string]command
	output io.Writer
}

func New() *Cli {
	return &Cli{
		cmds:   make(map[string]command),
		output: os.Stdout,
	}
}

func (c *Cli) AddCommand(cmd command) error {
	if _, ok := c.cmds[cmd.Name()]; ok {
		return errors.New("command already existed")
	}

	c.cmds[cmd.Name()] = cmd
	return nil
}

func (c *Cli) Run(args []string) error {
	if len(args) == 0 {
		c.commandHelp()
		return errors.New("no command provided")
	}

	if _, ok := c.cmds[args[0]]; !ok {
		c.commandHelp()
		return errors.New("command not found")
	}

	return c.cmds[args[0]].Run(args[1:])
}

func (c *Cli) commandHelp() {
	cmds := maps.Keys(c.cmds)

	sort.Slice(cmds, func(i, j int) bool { return cmds[i] < cmds[j] })
	maxLen := len(slices.MaxFunc(cmds, func(a, b string) int {
		return cmp.Compare(len(a), len(b))
	}))

	fmtStr := fmt.Sprintf("  %%-%ds  %%s\n", maxLen+5)

	fmt.Println("Available Commands:")
	for i := 0; i < len(c.cmds); i++ {
		cmd := c.cmds[cmds[i]]
		fmt.Printf(fmtStr, cmd.Name(), cmd.Description())
	}
}
