package main

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
	"github.com/evrynet-official/evrynet-tools/lib/log"
	"github.com/evrynet-official/evrynet-tools/lib/txutil"
	sc "github.com/evrynet-official/evrynet-tools/stakingcontract"
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
	transOpts := bind.NewKeyedTransactor(stakingClient.SenderPk)
	transOpts.GasLimit = stakingClient.GasLimit
	transOpts.Value = stakingClient.Amount

	tx, err := stakingClient.Vote(transOpts)
	if err != nil {
		zap.Errorw("votes for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	if err = txutil.CheckTransStatus(stakingClient.Client, tx); err != nil {
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

	transOpts := bind.NewKeyedTransactor(stakingClient.SenderPk)
	transOpts.GasLimit = stakingClient.GasLimit

	tx, err := stakingClient.UnVote(transOpts)
	if err != nil {
		zap.Errorw("un-votes for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	if err = txutil.CheckTransStatus(stakingClient.Client, tx); err != nil {
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

	transOpts := bind.NewKeyedTransactor(stakingClient.SenderPk)
	transOpts.GasLimit = stakingClient.GasLimit

	tx, err := stakingClient.Resign(transOpts)
	if err != nil {
		zap.Errorw("resigns for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	if err = txutil.CheckTransStatus(stakingClient.Client, tx); err != nil {
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

	transOpts := bind.NewKeyedTransactor(stakingClient.SenderPk)
	transOpts.GasLimit = stakingClient.GasLimit

	tx, err := stakingClient.Register(transOpts)
	if err != nil {
		zap.Errorw("registers for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	if err = txutil.CheckTransStatus(stakingClient.Client, tx); err != nil {
		zap.Errorw("registers for candidate got error ", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	zap.Infow("registers for candidate is finished", "candidate", stakingClient.Candidate.Hex())
	return nil
}

func getCandidates(ctx *cli.Context) error {
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
	candidates, err := stakingClient.GetAllCandidates(nil)
	if err != nil {
		zap.Errorw("GetAllCandidates returns error", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}

	sc.PrintCandidates(candidates)
	fmt.Printf("There are (%v) candidates\n", len(candidates))
	return nil
}

func getVoters(ctx *cli.Context) error {
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
	voters, err := stakingClient.GetVoters(nil)
	if err != nil {
		zap.Errorw("GetVoters returns error", "candidate", stakingClient.Candidate.Hex(), "error", err)
		return err
	}
	sc.PrintCandidates(voters)
	fmt.Printf("There are (%v) voters had voted for candidate (%s)\n", len(voters), stakingClient.Candidate.Hex())
	return nil
}
