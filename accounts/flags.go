package accounts

import (
	"github.com/urfave/cli"
)

var (
	// NumAccountsFlag the number of accounts want to generate
	NumAccountsFlag = cli.IntFlag{
		Name:  "num",
		Usage: "Number of accounts want to generate",
		Value: 4,
	}
	// SeedFlag to generate private key account
	SeedFlag = cli.StringFlag{
		Name:  "seed",
		Usage: "Seed to generate private key account",
		Value: "evrynet",
	}
)

// NewAccountsFlags return flags to generate accounts
func NewAccountsFlags() []cli.Flag {
	return []cli.Flag{NumAccountsFlag, SeedFlag}
}
