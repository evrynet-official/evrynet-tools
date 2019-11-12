package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/accounts"
)

func main() {
	app := cli.NewApp()
	app.Name = "accounts"
	app.Usage = "The accounts command line interface"
	app.Flags = append(app.Flags, accounts.NewAccountsFlags()...)
	app.Commands = []cli.Command{
		{
			Action:      accounts.CreateAccounts,
			Name:        "create",
			Usage:       "Create accounts",
			ArgsUsage:   "<create accounts>",
			Flags:       accounts.NewAccountsFlags(),
			Description: `To prepare accounts`,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
