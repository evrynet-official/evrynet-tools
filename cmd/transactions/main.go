package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/transactions"
)

func main() {
	app := cli.NewApp()
	app.Name = "transactions"
	app.Usage = "The transactions command line interface"
	app.Commands = []cli.Command{
		{
			Action:      transactions.SendTxFlood,
			Name:        "flood",
			Usage:       "Send tx flood",
			ArgsUsage:   "<send tx flood>",
			Flags:       transactions.NewTxFloodFlags(),
			Description: `Send tx flood`,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
