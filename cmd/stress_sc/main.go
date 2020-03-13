package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/lib/node"
	sc "github.com/evrynet-official/evrynet-tools/stakingcontract"
)

func main() {
	app := cli.NewApp()
	app.Name = "stress-test tool"
	app.Usage = "testing for staking contract"
	app.Version = "0.0.1"
	app.Commands = stakingCommands()

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func stakingCommands() []cli.Command {
	stressFlags := sc.NewStressTestFlag()
	stressFlags = append(stressFlags, node.NewEvrynetNodeFlags()...)

	stressVotesCmd := cli.Command{
		Action:      stressVoters,
		Name:        "stressvotes",
		Usage:       "sends vote from list voter to a candidate",
		Description: "sends vote from list voter to a candidate",
		Flags:       stressFlags,
	}

	return []cli.Command{stressVotesCmd}
}
