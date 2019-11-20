package app

import (
	"github.com/evrynet-official/evrynet-client/ethclient"
	"github.com/urfave/cli"
)

const (
	rpcEndpointFlag    = "rpcendpoint"
	defaultRPCEndpoint = "http://0.0.0.0:22001"
)

// NewEvrynetNodeFlags return flags to EvrynetNode
func NewEvrynetNodeFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  rpcEndpointFlag,
			Usage: "RPC endpoint to send request",
			Value: defaultRPCEndpoint,
		}}
}

// NewEvrynetClientFromFlags returns Evrynet client from flag variable, or error if occurs
func NewEvrynetClientFromFlags(ctx *cli.Context) (*ethclient.Client, error) {
	evrynetClientURL := ctx.String(rpcEndpointFlag)
	return ethclient.Dial(evrynetClientURL)
}
