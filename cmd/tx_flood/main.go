package main

import (
	"fmt"
	"os"

	"github.com/evrynet-official/evrynet-tools/accounts"
	"github.com/evrynet-official/evrynet-tools/tx_flood"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tx_flood"
	app.Usage = "The tx_flood command line interface"
	app.Version = "0.0.1"
	app.Flags = append(app.Flags, accounts.NewAccountsFlags()...)
	app.Flags = append(app.Flags, tx_flood.NewTxFloodFlags()...)
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	tf, err := tx_flood.NewTxFloodFromFlags(c)
	if err != nil {
		return err
	}

	if err := tf.Start(); err != nil {
		return err
	}
	return nil
}
