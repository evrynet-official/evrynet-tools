package main

import (
	"context"
	"errors"
	"time"

	"github.com/urfave/cli"

	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/evrclient"
	"github.com/evrynet-official/evrynet-tools/lib/log"
	sc "github.com/evrynet-official/evrynet-tools/smartcontract"
)

func vote(ctx *cli.Context) error {
	zap, flush, err := log.NewSugaredLogger(ctx)
	if err != nil {
		return err
	}
	defer flush()
	stakingClient, err := sc.NewNewStakingFromFlags(ctx, zap)
	if err != nil {
		zap.Errorw("cannot initialize a staking contract client ", "error", err)
		return err
	}
	tx, err := stakingClient.Vote(nil)
	if err != nil {
		zap.Errorw("votes for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	if err = checkTransStatus(stakingClient.Client, tx); err != nil {
		zap.Errorw("votes for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	zap.Infow("votes for candidate is finished", "candidate", stakingClient.Candidate.Hex())
	return nil
}

func unVote(ctx *cli.Context) error {
	zap, flush, err := log.NewSugaredLogger(ctx)
	if err != nil {
		return err
	}
	defer flush()
	stakingClient, err := sc.NewNewStakingFromFlags(ctx, zap)
	if err != nil {
		zap.Errorw("cannot initialize a staking contract client ", "error", err)
		return err
	}
	tx, err := stakingClient.UnVote(nil)
	if err != nil {
		zap.Errorw("un-votes for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	if err = checkTransStatus(stakingClient.Client, tx); err != nil {
		zap.Errorw("un-votes for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	zap.Infow("un-votes for candidate is finished", "candidate", stakingClient.Candidate.Hex())
	return nil
}

func resign(ctx *cli.Context) error {
	zap, flush, err := log.NewSugaredLogger(ctx)
	if err != nil {
		return err
	}
	defer flush()
	stakingClient, err := sc.NewNewStakingFromFlags(ctx, zap)
	if err != nil {
		zap.Errorw("cannot initialize a staking contract client ", "error", err)
		return err
	}
	tx, err := stakingClient.Resign(nil)
	if err != nil {
		zap.Errorw("resigns for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	if err = checkTransStatus(stakingClient.Client, tx); err != nil {
		zap.Errorw("resigns for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	zap.Infow("resigns for candidate is finished", "candidate", stakingClient.Candidate.Hex())
	return nil
}

func register(ctx *cli.Context) error {
	zap, flush, err := log.NewSugaredLogger(ctx)
	if err != nil {
		return err
	}
	defer flush()
	stakingClient, err := sc.NewNewStakingFromFlags(ctx, zap)
	if err != nil {
		zap.Errorw("cannot initialize a staking contract client ", "error", err)
		return err
	}
	tx, err := stakingClient.Register(nil)
	if err != nil {
		zap.Errorw("registers for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	if err = checkTransStatus(stakingClient.Client, tx); err != nil {
		zap.Errorw("registers for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	zap.Infow("registers for candidate is finished", "candidate", stakingClient.Candidate.Hex())
	return nil
}

func checkTransStatus(client *evrclient.Client, tx *types.Transaction) error {
	var err error
	for i := 0; i < 10; i++ {
		var receipt *types.Receipt
		receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
		if err == nil {
			if receipt.Status != uint64(1) {
				return errors.New("failed to send transaction")
			}
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return err
}
