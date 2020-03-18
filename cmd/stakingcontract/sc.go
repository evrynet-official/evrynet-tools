package main

import (
	"fmt"

	"github.com/urfave/cli"

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

	tx, err := stakingClient.Vote()
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

	tx, err := stakingClient.UnVote()
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

	tx, err := stakingClient.Resign()
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

	tx, err := stakingClient.Register()
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
	fmt.Printf("There were (%v) voters voting for candidate (%s)\n", len(voters), stakingClient.Candidate.Hex())
	return nil
}
