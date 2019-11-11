package main

import (
	"github.com/urfave/cli"
)

var (
	NumAccountsFlag = cli.IntFlag{
		Name:  "num",
		Usage: "Number of accounts want to generate",
		Value: 4,
	}
	SeedFlag = cli.StringFlag{
		Name:  "seed",
		Usage: "Seed to generate private key account",
		Value: "evrynet",
	}
)
