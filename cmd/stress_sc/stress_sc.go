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
	"github.com/evrynet-official/evrynet-tools/lib/txutil"
	sc "github.com/evrynet-official/evrynet-tools/stakingcontract"
)

func stressVoters(ctx *cli.Context) error {
	zap, flush, err := log.NewSugaredLogger(ctx)
	if err != nil {
		return err
	}
	defer flush()
	stakingClient, err := sc.NewContractClientFromFlags(ctx, zap)
	if err != nil {
		stakingClient.Logger.Errorw("cannot initialize a staking contract client ", "error", err)
		return err
	}

	accounts, err := generateAccounts(stakingClient.Logger, stakingClient.NumVoter)
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

	timeStats(stakingClient)
	return nil
}

func timeStats(stakingClient *sc.ContractClient) {
	stakingClient.Logger.Infow("start getVoters from SC")
	start := time.Now()
	voters, err := stakingClient.GetVoters(nil)
	if err != nil {
		return
	}
	stakingClient.Logger.Infow("getVoters", "number of Voter", len(voters), "elapsed", common.PrettyDuration(time.Since(start)))

	stakingClient.Logger.Infow("start getVoterStakes from SC for array voters")
	start = time.Now()
	_, err = stakingClient.GetVoterStakes(nil, voters)
	if err != nil {
		return
	}
	stakingClient.Logger.Infow("getVoterStakes", "number of Voter", len(voters), "elapsed", common.PrettyDuration(time.Since(start)))

}

func voteForCandidate(contractClient *sc.ContractClient, votes []*accounts.Account, candidate common.Address) error {
	var (
		gr     = errgroup.Group{}
		logger = contractClient.Logger.With("func", "voteForCandidate", "candidate", candidate.Hex())
	)

	batchSize := int(math.Floor(float64(len(votes)) / float64(contractClient.NumWorker)))
	for workerIndex := 0; workerIndex <= contractClient.NumWorker; workerIndex++ {
		var (
			from = workerIndex * batchSize
			to   = (workerIndex + 1) * batchSize
		)

		if workerIndex == contractClient.NumWorker {
			to = len(votes)
		}
		gr.Go(func() error {
			var (
				numberSuccess = int64(0)
				totalTime     = int64(0)
			)
			for i := from; i < to; i++ {
				addr, voterPk := votes[i].Address, votes[i].PriKey
				nonce, err := contractClient.Client.PendingNonceAt(context.Background(), addr)
				if err != nil {
					return err
				}

				optTrans := bind.NewKeyedTransactor(voterPk)
				optTrans.Nonce = new(big.Int).SetUint64(nonce)
				contractClient.TranOps = optTrans

				logger.Infow("begin vote for candidate", "number", (i + 1), "account", addr, "amount", contractClient.Amount)
				start := time.Now()
				tx, err := contractClient.Vote()
				if err != nil {
					logger.Errorw("failed to vote for candidate", "number", (i + 1), "error", err)
				} else {
					wErr := txutil.CheckTransStatus(contractClient.Client, tx)
					if wErr != nil {
						logger.Errorw("failed to checks the voting to candidate", "number", (i + 1), "account", addr, "error", wErr)
					} else {
						numberSuccess = numberSuccess + 1
						totalTime = totalTime + time.Since(start).Milliseconds()
						logger.Infow("account have sent a vote with success", "number", (i + 1), "account", addr)
					}
				}

			}

			if numberSuccess > 0 {
				avgTime := totalTime / numberSuccess
				logger.Infow("************************** summary", "voters", numberSuccess, "total time (ms)", totalTime, "avg (ms)", avgTime)
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
