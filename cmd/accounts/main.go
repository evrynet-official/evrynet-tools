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
	app.Commands = createCommands()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func createCommands() []cli.Command {
	createAccountsCmd := cli.Command{
		Action:      accounts.CreateAccounts,
		Name:        "create",
		Usage:       "Create accounts",
		Description: `To prepare accounts`,
	}
	createAccountsCmd.Flags = accounts.NewAccountsFlags()
	createAccountsCmd.Flags = append(createAccountsCmd.Flags, libApp.NewEvrynetNodeFlags()...)

	depositCmd := cli.Command{
		Action:      accounts.CreateAccountsAndDeposit,
		Name:        "deposit",
		Usage:       "Deposit EVR to the generated accounts",
		Description: `Deposit EVR to the generated accounts`,
	}
	depositCmd.Flags = accounts.NewDepositFlags()
	depositCmd.Flags = append(depositCmd.Flags, libApp.NewEvrynetNodeFlags()...)

	return []cli.Command{createAccountsCmd, depositCmd}
}
