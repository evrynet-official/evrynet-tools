package tx_metric

import (
	"sync"
	"time"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/lib/node"
)

const (
	startBlockNumber = "start-block"
	numberOfBlock    = "num-block"
	duration         = "duration"
)

// NewTxMetricFlags return flags to tx metric
func NewTxMetricFlags() []cli.Flag {
	return []cli.Flag{
		cli.Uint64Flag{
			Name:  startBlockNumber,
			Usage: "Where blocknumber start at",
			Value: 0,
		}, cli.DurationFlag{
			Name:  duration,
			Usage: "Duration to calculate metric",
			Value: 60 * time.Second,
		}, cli.Uint64Flag{
			Name:  numberOfBlock,
			Usage: "Duration to calculate metric",
			Value: 60,
		},
	}
}

// NewTxMetricFromFlags will init metric flags
func NewTxMetricFromFlags(ctx *cli.Context) (tm *TxMetric, err error) {
	tm = &TxMetric{
		StartBlockNumber: ctx.Uint64(startBlockNumber),
		Duration:         ctx.Duration(duration),
		NumBlock:         ctx.Uint64(numberOfBlock),
		mu:               &sync.Mutex{},
		minuteStats:      []int64{},
	}

	tm.EvrClient, err = node.NewEvrynetClientFromFlags(ctx)
	if err != nil {
		return nil, err
	}
	return tm, nil
}
