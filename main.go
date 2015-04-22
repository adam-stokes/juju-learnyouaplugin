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
	envName string
	description bool
}

var doc = `Run a command on target machine(s)

This example plugin mimics what "juju run" does.

eg.

juju lyaplugin -m 1 -e local "touch /tmp/testfile"
`

func (c *LYAPluginCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "lyaplugin",
		Args:    "<commands>",
		Purpose: "Run a command on remote target",
		Doc:     doc,
	}
}

func (c *LYAPluginCommand) SetFlags(f *gnuflag.FlagSet) {
	f.BoolVar(&c.description, "description", false, "Plugin Description")
	f.Var(cmd.NewStringsValue(nil, &c.machines), "machine", "one or more machine ids")
	f.Var(cmd.NewStringsValue(nil, &c.machines), "m", "")
	f.StringVar(&c.envName, "e", "local", "Juju environment")
	f.StringVar(&c.envName, "environment", "local", "")
}

func (c *LYAPluginCommand) Init(args []string) error {
	if c.description {
		fmt.Println(doc)
		os.Exit(0)
	}
	if len(args) == 0 {
		return fmt.Errorf("no commands specified")
	}
	if c.envName == "" {
		return fmt.Errorf("Juju environment must be specified.")
	}
	c.commands, args = args[0], args[1:]
	if len(c.machines) == 0 {
		return fmt.Errorf("You must specify a target with --machine, -m")
	}

	for _, machineId := range c.machines {
		if !names.IsValidMachine(machineId) {
			return fmt.Errorf("(%s) not a valid machine id.", machineId)
		}
	}
	return cmd.CheckEmpty(args)
}

func (c *LYAPluginCommand) Run(ctx *cmd.Context) error {
	c.SetEnvName(c.envName)
	client, err := c.NewAPIClient()
	if err != nil {
		return fmt.Errorf("Failed to load api client: %s", err)
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
	if len(runResults) == 1 {
		result := runResults[0]
		logger.Infof("Result: out(%s), err(%s), code(%d)", result.Stdout, result.Stderr, result.Code)
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
