package node

import (
	"github.com/evrynet-official/evrynet-client/ethclient"

	"github.com/urfave/cli"
)

const (
	rpcEndpointFlag = "rpcendpoint"
	evrynetEndpoint = "http://0.0.0.0:22001"
)

// EvrynetEndpoint returns configured Evrynet node endpoint.
func EvrynetEndpoint() string {
	return evrynetEndpoint
}

// NewEvrynetNodeFlags return flags to EvrynetNode
func NewEvrynetNodeFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  rpcEndpointFlag,
			Usage: "RPC endpoint to send request",
			Value: EvrynetEndpoint(),
		}}
}

// NewEvrynetClientFromFlags returns Evrynet client from flag variable, or error if occurs
func NewEvrynetClientFromFlags(ctx *cli.Context) (*ethclient.Client, error) {
	evrynetClientURL := ctx.String(rpcEndpointFlag)
	return ethclient.Dial(evrynetClientURL)
}