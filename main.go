// Copyright 2015 Adam Stokes <adam.stokes@ubuntu.com>

package main

import (
	"fmt"
	"github.com/juju/cmd"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/cmd/envcmd"
	"github.com/juju/juju/juju"
	"github.com/juju/loggo"
	"github.com/juju/names"
	"launchpad.net/gnuflag"
	"os"
	"time"

	_ "github.com/juju/juju/provider/all"
)

var logger = loggo.GetLogger("juju.plugin.lyaplugin")

// LYAPluginCommand represents the basic command structure
// and a map of known machines from juju's status
type LYAPluginCommand struct {
	envcmd.EnvCommandBase
	out      cmd.Output
	timeout  time.Duration
	machines []string
	services []string
	units    []string
	commands string
}

var doc = "Run a command on target machine(s)"

func (c *LYAPluginCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "lyaplugin",
		Args:    "<commands>",
		Purpose: "Run a command on remote target",
		Doc:     doc,
	}
}

func (c *LYAPluginCommand) SetFlags(f *gnuflag.FlagSet) {
	f.Var(cmd.NewStringsValue(nil, &c.machines), "machine", "one or more machine ids")
}

func (c *LYAPluginCommand) Init(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no commands specified")
	}
	c.commands, args = args[0], args[1:]
	if len(c.machines) == 0 {
		return fmt.Errorf("You must specify a target with --machine")
	}

	for _, machineId := range c.machines {
		if !names.IsValidMachine(machineId) {
			return fmt.Errorf("(%s) not a valid machine id.", machineId)
		}
	}
	return cmd.CheckEmpty(args)
}

func (c *LYAPluginCommand) Run(ctx *cmd.Context) error {
	envName, err := envcmd.GetDefaultEnvironment()
	if err != nil {
		return fmt.Errorf("could not get default environment (%s)", err)
	}
	c.SetEnvName(envName)
	client, err := c.NewAPIClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	var runResults []params.RunResult

	logger.Infof("Running cmd: %s on machine: %s", c.commands, c.machines[0])
	params := params.RunParams{
		Commands: c.commands,
		Timeout:  c.timeout,
		Machines: c.machines,
		Services: c.services,
		Units:    c.units,
	}
	runResults, err = client.Run(params)
	if err != nil {
		fmt.Errorf("An error occurred: %s", err)
	}
	for _, result := range runResults {
		logger.Infof("Run result: (%d) (%s) (%s)", result.Code, result.Stdout, result.Stderr)
	}
	return nil
}

func main() {
	loggo.ConfigureLoggers("<root>=INFO")
	err := juju.InitJujuHome()
	if err != nil {
		panic(err)
	}

	ctx, err := cmd.DefaultContext()
	if err != nil {
		panic(err)
	}

	c := &LYAPluginCommand{}
	cmd.Main(c, ctx, os.Args[1:])
}
