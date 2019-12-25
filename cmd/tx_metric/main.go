package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/lib/node"
	"github.com/evrynet-official/evrynet-tools/tx_metric"
)

func main() {
	app := cli.NewApp()
	app.Name = "tx_metric"
	app.Usage = "The tx_metric command line interface"
	app.Version = "0.0.1"
	app.Commands = metricsCommand()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func byTime(c *cli.Context) error {
	tm, err := tx_metric.NewTxMetricFromFlags(c)
	if err != nil {
		return err
	}

	if err := tm.MetricByTime(); err != nil {
		return err
	}
	return nil
}

func byBlock(c *cli.Context) error {
	tm, err := tx_metric.NewTxMetricFromFlags(c)
	if err != nil {
		return err
	}

	if err := tm.MetricByBlock(); err != nil {
		return err
	}
	return nil
}

func metricsCommand() []cli.Command {
	byBlockCommand := cli.Command{
		Action:      byBlock,
		Name:        "byblock",
		Usage:       "calculate metrics by startBlock + number of block",
		Description: "The total metrics of the network will be aggregated by the first block that is higher than input block and has more than 1 transaction",
	}

	bytimeCommand := cli.Command{
		Action:      byTime,
		Name:        "bytime",
		Usage:       "calculate metrics by startBlock up until a duration",
		Description: "The total metrics of the network will be aggregated by the first block that is higher than input block and has more than 1 transaction",
	}
	flags := append(tx_metric.NewTxMetricFlags(), node.NewEvrynetNodeFlags()...)
	byBlockCommand.Flags = flags
	bytimeCommand.Flags = flags

	return []cli.Command{byBlockCommand, bytimeCommand}
}
