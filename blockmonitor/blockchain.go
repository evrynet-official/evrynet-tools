package blockmonitor

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-client/ethclient"
	"github.com/evrynet-official/evrynet-tools/lib/node"
)

var (
	timeTickerFlag = cli.IntFlag{
		Name:  "duration",
		Usage: "The duration delay for each times to checks the situation of block (seconds)",
		Value: 60,
	}
)

// NewBlcClientFlag returns flags for block-chain
func NewBlcClientFlag() []cli.Flag {
	return []cli.Flag{timeTickerFlag}
}

type Blockchain struct {
	Client      *ethclient.Client
	LatestBlock *big.Int
	Duration    time.Duration
}

func NewBlcClientFromFlags(ctx *cli.Context) (*Blockchain, error) {
	var (
		delay = ctx.Duration(timeTickerFlag.Name)
	)

	client, err := node.NewEvrynetClientFromFlags(ctx)
	if err != nil {
		return nil, err
	}
	blcClient := &Blockchain{
		Client:      client,
		LatestBlock: new(big.Int).SetUint64(0),
		Duration:    delay,
	}
	return blcClient, nil
}

func (blc *Blockchain) GetLastBlock() (*big.Int, error) {
	header, err := blc.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	if header == nil {
		return nil, errors.New("can not get latest block")
	}
	return header.Number, nil
}
