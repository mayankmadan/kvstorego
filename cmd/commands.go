package cmd

import (
	"strings"
)

var Commands map[string]ICommand = make(map[string]ICommand)

type ICommand interface {
	GetName() string
	GetDescription() string
	GetNumberOfOperands() int
	Exec(operands []string) (string, bool, error) // true in case connection needs to be closed
	Init(name string, operandCount int, description string)
}

type Command struct {
	name             string
	description      string
	numberOfoperands int
}

func (c *Command) GetName() string {
	return strings.ToUpper(c.name)
}

func (c *Command) Init(name string, operandCount int, description string) {
	c.name = name
	c.description = description
	c.numberOfoperands = operandCount
}

func (c *Command) GetNumberOfOperands() int {
	return c.numberOfoperands
}

func (c *Command) GetDescription() string {
	return c.description
}

func (c *Command) Parse(operands []string) bool {
	return true
}

func GetCommand(name string) ICommand {
	if command, ok := Commands[strings.ToUpper(name)]; ok {
		return command
	}
	return nil
}

func Register(command ICommand) {
	if _, ok := Commands[command.GetName()]; ok {
		return
	}
	Commands[command.GetName()] = command
}
