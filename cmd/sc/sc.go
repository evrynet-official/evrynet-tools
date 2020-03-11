package main

import (
	"log"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/staking"
)

func vote(ctx *cli.Context) {
	stakingClient, err := staking.NewNewStakingFromFlags(ctx)
	if err != nil {
		log.Printf("can not initialize a staking contract client %s", err.Error())
		return
	}
	_, err = stakingClient.Vote()
	if err != nil {
		log.Printf("un-vote for candidate %s got error: %w", stakingClient.Candidate.Hex(), err)
	}
}

func unVote(ctx *cli.Context) {
	stakingClient, err := staking.NewNewStakingFromFlags(ctx)
	if err != nil {
		log.Printf("can not initialize a staking contract client %s", err.Error())
		return
	}
	_, err = stakingClient.UnVote()
	if err != nil {
		log.Printf("un-vote for candidate %s got error: %w", stakingClient.Candidate.Hex(), err)
	}
}

func resign(ctx *cli.Context) {
	stakingClient, err := staking.NewNewStakingFromFlags(ctx)
	if err != nil {
		log.Printf("can not initialize a staking contract client %s", err.Error())
		return
	}
	_, err = stakingClient.Resign()
	if err != nil {
		log.Printf("un-vote for candidate %s got error: %w", stakingClient.Candidate.Hex(), err)
	}
}

func register(ctx *cli.Context) {
	stakingClient, err := staking.NewNewStakingFromFlags(ctx)
	if err != nil {
		log.Printf("can not initialize a staking contract client %s", err.Error())
		return
	}
	_, err = stakingClient.Register()
	if err != nil {
		log.Printf("un-vote for candidate %s got error: %w", stakingClient.Candidate.Hex(), err)
	}
}
