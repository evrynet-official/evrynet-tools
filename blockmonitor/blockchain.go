package blockmonitor

import (
	"context"
	"errors"
	"math/big"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-client/ethclient"
	"github.com/evrynet-official/evrynet-tools/lib/node"
)

type Blockchain struct {
	Client      *ethclient.Client
	LatestBlock *big.Int
}

func NewBlcClientFromFlags(ctx *cli.Context) (*Blockchain, error) {
	client, err := node.NewEvrynetClientFromFlags(ctx)
	if err != nil {
		return nil, err
	}
	blcClient := &Blockchain{
		Client:      client,
		LatestBlock: new(big.Int).SetUint64(0),
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
