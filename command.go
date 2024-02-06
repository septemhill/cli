package cli

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	cliTag              = "cli"
	flagsClause         = "flags:"
	descriptionClause   = "desc:"
	shortClause         = "short:"
	clauseSeprator      = ";"
	clauseValueSeprator = ","
)

type command interface {
	Name() string
	Run([]string) error
	Description() string
}

type internalFlagsInfo struct {
	Short       string
	Flag        string
	Description string
	Value       reflect.Value
}

type CommandAction[T any] func(T) error

type Command[T any] struct {
	Cmd    string
	Flags  T
	Desc   string
	Action CommandAction[T]
	intl   []*internalFlagsInfo
}

func NewCommand[T any](cmd string, defaultFlags T, action CommandAction[T], desc string) *Command[T] {
	c := &Command[T]{
		Cmd:    cmd,
		Flags:  defaultFlags,
		Action: action,
		Desc:   desc,
		intl:   make([]*internalFlagsInfo, 0),
	}

	c.flagsProcess()

	return c
}

func (c *Command[T]) Name() string {
	return c.Cmd
}

func (c *Command[T]) Description() string {
	return c.Desc
}

func (c *Command[T]) Run(args []string) error {
	if _, err := c.parseFlag(args); err != nil {
		c.commandHelp()
		return err
	}
	return c.Action(c.Flags)
}

func (c *Command[T]) flagsProcess() {
	flagType := reflect.TypeOf(c.Flags).Elem()

	// TODO: asssume c.Flags is ptr here, need to check non-ptr case
	v := reflect.ValueOf(c.Flags).Elem()

	for i := 0; i < flagType.NumField(); i++ {
		intl := &internalFlagsInfo{}
		intl.Value = v.Field(i)
		intl.Flag, intl.Short, intl.Description = c.tagProcess(flagType.Field(i).Tag)
		c.intl = append(c.intl, intl)
	}
}

func (c *Command[T]) tagProcess(tag reflect.StructTag) (string, string, string) {
	tagValues := tag.Get(cliTag)
	clauses := strings.Split(tagValues, clauseSeprator)

	intl := &internalFlagsInfo{}
	for i := 0; i < len(clauses); i++ {
		clause := strings.Trim(clauses[i], " ")
		switch {
		case strings.HasPrefix(clause, flagsClause):
			intl.Flag = c.flagClauseProcess(clause)
		case strings.HasPrefix(clause, shortClause):
			intl.Short = c.shortClauseProcess(clause)
		case strings.HasPrefix(clause, descriptionClause):
			intl.Description = c.descClauseProcess(clause)
		}
	}
	return intl.Flag, intl.Short, intl.Description
}

func (c *Command[T]) flagClauseProcess(clause string) string {
	return "--" + clause[len(flagsClause):]
}

func (c *Command[T]) shortClauseProcess(clause string) string {
	return "-" + clause[len(shortClause):]
}

func (c *Command[T]) descClauseProcess(clause string) string {
	return clause[len(descriptionClause):]
}

func (c *Command[T]) parseFlag(flags []string) ([]string, error) {
	i := 0

	for i = 0; i < len(flags); {
		if !strings.HasPrefix(flags[i], "-") {
			break
		}

		if i+1 >= len(flags) || strings.HasPrefix(flags[i+1], "-") {
			if err := c.flagCheck(flags[i], "true"); err != nil {
				return nil, err
			}
			i += 1
			continue
		}

		if err := c.flagCheck(flags[i], flags[i+1]); err != nil {
			return nil, err
		}
		i += 2
	}

	return flags[i:], nil
}

func (c *Command[T]) flagCheck(flag string, value string) error {
	for i := 0; i < len(c.intl); i++ {
		if c.intl[i].Flag == flag || c.intl[i].Short == flag {
			c.flagSet(c.intl[i].Value, value)
			return nil
		}
	}

	return fmt.Errorf("unknown flag: %s", flag)
}

func (c *Command[T]) flagSet(rv reflect.Value, value string) error {

	switch rv.Kind() {
	case reflect.Bool:
		rv.SetBool(value == "true")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		rv.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		rv.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		rv.SetFloat(f)
	case reflect.String:
		rv.SetString(value)
	}

	return nil
}

func (c *Command[T]) commandHelp() {
	fmt.Println(c.Cmd)
	for i := 0; i < len(c.intl); i++ {
		fmt.Printf("  %s,%-10s 	 %s\n", c.intl[i].Short, c.intl[i].Flag, c.intl[i].Description)
	}
}
