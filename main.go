// Copyright 2015 Adam Stokes <adam.stokes@ubuntu.com>

package main

import (
	"fmt"
	"github.com/juju/juju/juju"
	"github.com/juju/juju/juju/osenv"
	"github.com/juju/juju/cmd/envcmd"
	"github.com/juju/juju/state"

	_ "github.com/juju/juju/provider/all"
)

type LYAPluginCommand struct {
	envcmd.EnvCommandBase
	MachineMap map[string]*state.Machine
}

func (c *LYAPluginCommand) Info(envName string) error {
	client, err := juju.NewAPIClientFromName(envName)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	c.MachineMap = make(map[string]*state.Machine)

	status, err := client.Status([]string{})
	if err != nil {
		panic(err)
	}
	for _, m := range status.Machines {
		fmt.Printf("Found machine (%s)", m.Id)
	}
	return nil
}

func main() {
	fmt.Println("Starting learnyouaplugin!")
	err := juju.InitJujuHome()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Using JUJU_HOME: (%s)\n", osenv.JujuHome())
	fmt.Println("Grabbing State information.")
	c := &LYAPluginCommand{}
	c.Info("local")
}
