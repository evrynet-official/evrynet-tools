package tx_metric

import (
	"time"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-client/ethclient"
)

const (
	rpcEndpointFlag    = "rpcendpoint"
	defaultRPCEndpoint = "http://0.0.0.0:22001"
	startBlockNumber   = "start-block"
	duration           = "duration"
)

// NewTxMetricFlags return flags to tx metric
func NewTxMetricFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "rpcendpoint",
			Usage: "RPC endpoint to send request",
			Value: defaultRPCEndpoint,
		}, cli.IntFlag{
			Name:  startBlockNumber,
			Usage: "Where blocknumber start at",
			Value: 0,
		}, cli.DurationFlag{
			Name:  duration,
			Usage: "Duration to calculate metric",
			Value: 60 * time.Second,
		}}
}

// NewEvrClientFromFlag returns Ethereum client from flag variable, or error if occurs
func NewEvrClientFromFlag(ctx *cli.Context) (*ethclient.Client, error) {
	evrNodeURL := ctx.String(rpcEndpointFlag)
	return ethclient.Dial(evrNodeURL)
}

// NewTxMetricFromFlags will init metric flags
func NewTxMetricFromFlags(ctx *cli.Context) (tm *TxMetric, err error) {
	tm = &TxMetric{
		StartBlockNumber: ctx.Int(startBlockNumber),
		Duration:         ctx.Duration(duration),
	}

	tm.EvrClient, err = NewEvrClientFromFlag(ctx)
	if err != nil {
		return nil, err
	}
	return tm, nil
}
