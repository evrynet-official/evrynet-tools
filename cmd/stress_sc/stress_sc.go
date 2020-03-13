package main

import (
	"context"
	"math/big"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	"golang.org/x/sync/errgroup"

	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/core/types"

	"github.com/evrynet-official/evrynet-tools/accounts"
	"github.com/evrynet-official/evrynet-tools/accounts/depositor"
	"github.com/evrynet-official/evrynet-tools/lib/log"
	sc "github.com/evrynet-official/evrynet-tools/stakingcontract"
)

func stressVoters(ctx *cli.Context) error {
	zap, flush, err := log.NewSugaredLogger(ctx)
	if err != nil {
		return err
	}
	defer flush()
	stakingClient, err := sc.NewNewStakingFromFlags(ctx, zap)
	if err != nil {
		stakingClient.Logger.Errorw("cannot initialize a staking contract client ", "error", err)
		return err
	}

	voters, err := generateAccounts(stakingClient.Logger, stakingClient.NumVoter)
	if err != nil {
		return err
	}
	err = sendEvrToken(stakingClient, voters)
	if err != nil {
		return err
	}
	err = voteForCandidate(stakingClient, voters, stakingClient.Candidate)
	if err != nil {
		return err
	}

	stakingClient.Logger.Infow("all voters have sent votes for candidate", "candidate", stakingClient.Candidate.Hex())
	return nil
}

func voteForCandidate(contractClient *sc.ContractClient, votes []*accounts.Account, candidate common.Address) error {
	var (
		gr       = errgroup.Group{}
		logger   = contractClient.Logger.With("func", "voteForCandidate", "candidate", candidate.Hex())
		optTrans = bind.NewKeyedTransactor(contractClient.SenderPk)
		gasLimit = uint64(8000000)
		amount   = new(big.Int).SetUint64(50)
	)

	for i := 0; i < len(votes); i++ {
		addr, voterPk := votes[i].Address, votes[i].PriKey
		nonce, err := contractClient.Client.PendingNonceAt(context.Background(), addr)
		if err != nil {
			return err
		}

		optTrans = bind.NewKeyedTransactor(voterPk)
		optTrans.GasLimit = gasLimit
		optTrans.Nonce = new(big.Int).SetUint64(nonce)
		optTrans.Value = amount

		logger.Infow("begin vote for candidate", "account", addr)
		tx, err := contractClient.Contract.Vote(optTrans, candidate)
		if err != nil {
			logger.Errorw("failed to vote for candidate", "error", err)
			return err
		}
		nonce++
		gr.Go(func() error {
			_, wErr := sc.WaitForTx(contractClient.Client, tx.Hash())
			if wErr != nil {
				logger.Errorw("failed to vote for candidate", "error", wErr)
				return wErr
			}
			logger.Infow("account have sent a vote", "account", addr)
			return nil
		})
	}
	if err := gr.Wait(); err != nil {
		return err
	}
	return nil
}

func sendEvrToken(stakingClient *sc.ContractClient, voters []*accounts.Account) error {
	var (
		gasLimit       = uint64(1000000)
		expectedAmount = new(big.Int).Exp(new(big.Int).SetUint64(10), new(big.Int).SetUint64(18), nil)
	)

	optTrans := bind.NewKeyedTransactor(stakingClient.SenderPk)
	optTrans.Signer = func(signer types.Signer, from common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return types.SignTx(tx, signer, stakingClient.SenderPk)
	}
	dep := depositor.NewDepositor(stakingClient.Logger, optTrans, optTrans.From, voters, stakingClient.Client, expectedAmount, len(voters),
		depositor.WithGasLimit(gasLimit))

	return dep.DepositCoreAccounts()
}

func generateAccounts(logger *zap.SugaredLogger, numVoters int) ([]*accounts.Account, error) {
	// generate accounts
	accs, err := accounts.GenerateAccounts(numVoters, "evr")
	if err != nil {
		logger.Errorw("Fail to generate new account!", "Err:", err)
		return nil, err
	}
	return accs, nil
}
