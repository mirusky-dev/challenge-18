package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootShort = "GolangBoilerplate CLI"
	rootLong  = `GolangBoilerplate CLI

Web Framework Boilerplate

Find more information at:
	https://www.golangboilerplate.com`
)

func NewCmdRoot() *cobra.Command {

	cmd := &cobra.Command{
		Short: rootShort,
		Long:  rootLong,
	}

	cmd.AddCommand(
		newCmdAPI(),
		newCmdBackground(),
	)

	return cmd
}
