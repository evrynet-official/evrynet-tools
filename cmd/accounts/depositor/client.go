package depositor

import (
	"context"
	"math/big"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/ethclient"
)

// ClientInterface 
type ClientInterface interface {
	NonceAt(background context.Context, addresses common.Address, blockNumber *big.Int) (uint64, error)
	SuggestGasPrice(background context.Context) (*big.Int, error)
	SendTransaction(background context.Context, transaction *types.Transaction) error
	TransactionReceipt(background context.Context, hash common.Hash) (*types.Receipt, error)
	BalanceAt(background context.Context, addresses common.Address, blockNumber *big.Int) (*big.Int, error)
}

type ClientAPI struct {
	evrClient *ethclient.Client
}


func (client *ClientAPI) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return client.evrClient.NonceAt(ctx, account, blockNumber)
}

func (client *ClientAPI) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return client.evrClient.SuggestGasPrice(ctx)
}

func (client *ClientAPI) SendTransaction(ctx context.Context, transaction *types.Transaction) error {
	return client.evrClient.SendTransaction(ctx, transaction)
}

func (client *ClientAPI) TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	return client.evrClient.TransactionReceipt(ctx, hash)
}

func (client *ClientAPI) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return client.evrClient.BalanceAt(ctx, account, blockNumber)
}
