package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/accounts"
	libApp "github.com/evrynet-official/evrynet-tools/lib/app"
)

func main() {
	app := cli.NewApp()
	app.Name = "accounts"
	app.Usage = "The accounts command line interface"
	app.Version = "0.0.1"
	flags := accounts.NewAccountsFlags()
	flags = append(flags, libApp.NewEvrynetNodeFlags()...)
	app.Commands = []cli.Command{
		{
			Action:      accounts.CreateAccounts,
			Name:        "create",
			Usage:       "Create accounts",
			ArgsUsage:   "<create accounts>",
			Flags:       flags,
			Description: `To prepare accounts`,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
