package depositor

import (
	"fmt"
	"math/big"

	"github.com/evrynet-official/evrynet-client/accounts/abi/bind"
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/evrynet-official/evrynet-tools/accounts"
	"github.com/evrynet-official/evrynet-tools/lib/node"
)

var (
	senderPkFlag = cli.StringFlag{
		Name:  "senderpk",
		Usage: "The private key of sender",
		Value: "ce900e4057ef7253ce737dccf3979ec4e74a19d595e8cc30c6c5ea92dfdd37f1",
	}
	expectedBalanceFlag = cli.StringFlag{
		Name:  "expectedbalance",
		Usage: "The expected balance of each account (wei)",
		Value: "1000000000000000000",
	}
	numberOfWorkerFlag = cli.IntFlag{
		Name:  "nworkers",
		Usage: "The number of worker for the program",
		Value: 1,
	}
)

// NewDepositFlags return flags to create a depositor
func NewDepositFlags() []cli.Flag {
	return []cli.Flag{accounts.NumAccountsFlag, accounts.SeedFlag, senderPkFlag, expectedBalanceFlag, numberOfWorkerFlag}
}

// NewDepositFlags return a ready-to-use depositor from cli
func NewDepositorFromFlag(ctx *cli.Context, logger *zap.SugaredLogger) (*Depositor, error) {
	var (
		senderPk = ctx.String(senderPkFlag.Name)
		amount   = ctx.String(expectedBalanceFlag.Name)
		gasLimit = big.NewInt(1000000).Uint64()
		nworker  = ctx.Int(numberOfWorkerFlag.Name)
	)

	expectedAmount, ok := new(big.Int).SetString(amount, 10)

	if !ok {
		return nil, fmt.Errorf("failed to parse expected amount from input %s", amount)
	}

	pk, err := crypto.HexToECDSA(senderPk)

	if err != nil {
		return nil, err
	}

	accs, err := accounts.GenerateAccountsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var (
		wAddrs []common.Address
		opt    = bind.NewKeyedTransactor(pk)
	)

	for i := 0; i < len(accs); i++ {
		wAddrs = append(wAddrs, accs[i].Address)
	}

	opt.Signer = func(signer types.Signer, from common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return types.SignTx(tx, signer, pk)
	}

	evrClient, err := node.NewEvrynetClientFromFlags(ctx)
	if err != nil {
		return nil, err
	}
	dep := NewDepositor(logger, opt, wAddrs, evrClient, expectedAmount,
		WithGasLimit(gasLimit), WithNumWorkers(nworker),
	)
	return dep, nil

}
