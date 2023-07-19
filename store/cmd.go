package store

import (
	"errors"
	"kvstore/cmd"
)

const (
	Exit string = "exit"
	Ping string = "ping"
	Set  string = "set"
	Get  string = "get"
	Del  string = "del"
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

	cmdDelete struct {
		db IDB
		cmd.Command
	}

	cmdExit struct {
		cmd.Command
	}
)

func (c *cmdGet) Exec(operands []string) *cmd.Result {
	val, ok := c.db.Get(operands[0])
	if ok {
		return &cmd.Result{Data: []cmd.Element{{ElementType: cmd.String, Value: val}}, Close: false, Err: nil}
	}
	return &cmd.Result{Data: nil, Close: false, Err: errors.New("key not found")}
}

func (c *cmdSet) Exec(operands []string) *cmd.Result {
	ok := c.db.Set(operands[0], operands[1])
	if ok {
		resultEl := cmd.Element{ElementType: cmd.String, Value: "Key Set Successfully"}
		return &cmd.Result{Data: []cmd.Element{resultEl}, Close: false, Err: nil}
	}
	return &cmd.Result{Data: nil, Close: false, Err: errors.New("error setting key")}
}

func (c *cmdDelete) Exec(operands []string) *cmd.Result {
	ok := c.db.Delete(operands[0])
	if ok {
		return &cmd.Result{Data: []cmd.Element{{ElementType: cmd.String, Value: "Key Deleted Successfully"}}, Close: false, Err: nil}
	}
	return &cmd.Result{Data: nil, Close: false, Err: errors.New("error deleting key")}
}

func (c *cmdPing) Exec(operands []string) *cmd.Result {
	return &cmd.Result{Data: []cmd.Element{{ElementType: cmd.String, Value: "PONG"}}, Close: false, Err: nil}
}

func (c *cmdExit) Exec(operands []string) *cmd.Result {
	return &cmd.Result{Data: []cmd.Element{{ElementType: cmd.String, Value: "Exiting"}}, Close: true, Err: nil}
}

func Init(db IDB) {
	get := &cmdGet{db: db}
	get.Init(Exit, 1, "Get Key from cluster")
	cmd.Register(get)
	set := &cmdSet{db: db}
	set.Init(Set, 2, "Set Key Value in cluster")
	cmd.Register(set)
	del := &cmdDelete{db: db}
	del.Init(Del, 1, "Delete Key From cluster")
	cmd.Register(del)
	ping := &cmdPing{}
	ping.Init(Ping, 0, "Ping cluster")
	cmd.Register(ping)
	exit := &cmdExit{}
	exit.Init(Exit, 0, "Close Connection")
	cmd.Register(exit)
}
