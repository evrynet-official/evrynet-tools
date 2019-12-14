package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/tx_metric"
)

func main() {
	app := cli.NewApp()
	app.Name = "tx_metric"
	app.Usage = "The tx_metric command line interface"
	app.Version = "0.0.1"
	app.Flags = append(app.Flags, tx_metric.NewTxMetricFlags()...)
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	tm, err := tx_metric.NewTxMetricFromFlags(c)
	if err != nil {
		return err
	}

	if err := tm.Start(); err != nil {
		return err
	}
	return nil
}
