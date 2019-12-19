package tx_metric

import (
	"time"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/lib/node"
)

const (
	startBlockNumber = "start-block"
	duration         = "duration"
)

// NewTxMetricFlags return flags to tx metric
func NewTxMetricFlags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:  startBlockNumber,
			Usage: "Where blocknumber start at",
			Value: 0,
		}, cli.DurationFlag{
			Name:  duration,
			Usage: "Duration to calculate metric",
			Value: 60 * time.Second,
		}}
}

// NewTxMetricFromFlags will init metric flags
func NewTxMetricFromFlags(ctx *cli.Context) (tm *TxMetric, err error) {
	tm = &TxMetric{
		StartBlockNumber: ctx.Int(startBlockNumber),
		Duration:         ctx.Duration(duration),
	}

	tm.EvrClient, err = node.NewEvrynetClientFromFlags(ctx)
	if err != nil {
		return nil, err
	}
	return tm, nil
}
