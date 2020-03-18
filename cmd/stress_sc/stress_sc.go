package main

import (
	"context"
	"math"
	"math/big"
	"time"

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

	accounts, err := generateAccounts(stakingClient.Logger, stakingClient.NumVoters)
	if err != nil {
		return err
	}
	err = sendEvrToken(stakingClient, accounts)
	if err != nil {
		return err
	}
	err = voteForCandidate(stakingClient, accounts, stakingClient.Candidate)
	if err != nil {
		return err
	}

	stakingClient.Logger.Infow("start getVoters from SC")
	start := time.Now()
	voters, err := stakingClient.GetVoters(nil)
	if err != nil {
		return err
	}
	stakingClient.Logger.Infow("getVoters", "number of Voter", len(voters), "elapsed", common.PrettyDuration(time.Since(start)))

	stakingClient.Logger.Infow("start getVoterStake from SC")
	start = time.Now()
	stake, err := stakingClient.GetVoterStake(nil, voters[0])
	if err != nil {
		return err
	}
	stakingClient.Logger.Infow("getVoterStake", "voter", voters[0], "stake", stake, "elapsed", common.PrettyDuration(time.Since(start)))
	return nil
}

func voteForCandidate(contractClient *sc.ContractClient, votes []*accounts.Account, candidate common.Address) error {
	var (
		gr       = errgroup.Group{}
		logger   = contractClient.Logger.With("func", "voteForCandidate", "candidate", candidate.Hex())
		optTrans *bind.TransactOpts
	)

	batchSize := int(math.Floor(float64(len(votes)) / float64(contractClient.NumWorkers)))
	for workerIndex := 0; workerIndex <= contractClient.NumWorkers; workerIndex++ {
		from := workerIndex * batchSize
		to := (workerIndex + 1) * batchSize
		if workerIndex == contractClient.NumWorkers {
			to = len(votes)
		}
		gr.Go(func() error {
			for i := from; i < to; i++ {
				addr, voterPk := votes[i].Address, votes[i].PriKey
				nonce, err := contractClient.Client.PendingNonceAt(context.Background(), addr)
				if err != nil {
					return err
				}

				optTrans = bind.NewKeyedTransactor(voterPk)
				optTrans.GasLimit = contractClient.TranOps.GasLimit
				optTrans.Nonce = new(big.Int).SetUint64(nonce)
				optTrans.Value = contractClient.TranOps.Value

				logger.Infow("begin vote for candidate", "number", (i + 1), "account", addr)
				tx, err := contractClient.Contract.Vote(optTrans, candidate)
				if err != nil {
					logger.Errorw("failed to vote for candidate", "number", (i + 1), "error", err)
					return err
				}
				_, wErr := sc.WaitForTx(contractClient.Client, tx.Hash())
				if wErr != nil {
					logger.Errorw("failed to vote for candidate", "number", "error", "account", addr, wErr)
					return wErr
				}
				logger.Infow("account have sent a vote", "account", addr)
			}

			return nil
		})
	}

	if err := gr.Wait(); err != nil {
		return err
	}
	logger.Infow("all voters have sent votes for candidate", "total_account", len(votes))
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
