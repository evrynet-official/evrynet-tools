package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/accounts"
	"github.com/evrynet-official/evrynet-tools/accounts/depositor"

	"github.com/evrynet-official/evrynet-tools/lib/node"
)

func main() {
	app := cli.NewApp()
	app.Name = "accounts"
	app.Usage = "The accounts command line interface"
	app.Version = "0.0.1"
	app.Commands = accountsCommand()

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func accountsCommand() []cli.Command {
	createAccountsCmd := cli.Command{
		Action:      generate,
		Name:        "generate",
		Usage:       "generate a number accounts based on a seed",
		Description: `To prepare accounts`,
	}
	createAccountsCmd.Flags = accounts.NewAccountsFlags()

	depositCmd := cli.Command{
		Action:      deposit,
		Name:        "deposit",
		Usage:       "Deposit EVR to the generated accounts",
		Description: `Deposit EVR to the generated accounts`,
	}
	depositCmd.Flags = depositor.NewDepositFlags()
	depositCmd.Flags = append(depositCmd.Flags, node.NewEvrynetNodeFlags()...)

	return []cli.Command{createAccountsCmd, depositCmd}
}
