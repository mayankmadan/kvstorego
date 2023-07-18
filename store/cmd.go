package store

import (
	"errors"
	"kvstore/cmd"
)

type (
	cmdGet struct {
		db IDB
		cmd.Command
	}

	cmdSet struct {
		db IDB
		cmd.Command
	}

	cmdPing struct {
		cmd.Command
	}

	cmdExit struct {
		cmd.Command
	}
)

func (c *cmdGet) Exec(operands []string) (string, bool, error) {
	val, ok := c.db.Get(operands[0])
	if ok {
		return val, false, nil
	}
	return "", false, errors.New("key not found")
}

func (c *cmdSet) Exec(operands []string) (string, bool, error) {
	ok := c.db.Set(operands[0], operands[1])
	if ok {
		return "Key Set Successfully", false, nil
	}
	return "", false, errors.New("error setting key")
}

func (c *cmdPing) Exec(operands []string) (string, bool, error) {
	return "Pong", false, nil
}

func (c *cmdExit) Exec(operands []string) (string, bool, error) {
	return "Exiting!", true, nil
}

func Init(db IDB) {
	get := &cmdGet{db: db}
	get.Init("get", 1, "Get Key from cluster")
	cmd.Register(get)
	set := &cmdSet{db: db}
	set.Init("set", 2, "Set Key Value in cluster")
	cmd.Register(set)
	ping := &cmdPing{}
	ping.Init("ping", 0, "Ping cluster")
	cmd.Register(ping)
	exit := &cmdExit{}
	exit.Init("exit", 0, "Close Connection")
	cmd.Register(exit)
}
