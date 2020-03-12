package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/lib/node"
	"github.com/evrynet-official/evrynet-tools/smartcontract"
)

func main() {
	app := cli.NewApp()
	app.Name = "stakingClient"
	app.Usage = "sends a vote/ unvote/ register/ resign for a candidate to a node"
	app.Version = "0.0.1"
	app.Commands = stakingCommands()

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func stakingCommands() []cli.Command {
	flags := sc.NewStakingFlag()
	flags = append(flags, node.NewEvrynetNodeFlags()...)

	voteCmd := cli.Command{
		Action:      vote,
		Name:        "vote",
		Usage:       "Sends a vote for a candidate",
		Description: `Sends a vote for a candidate`,
		Flags:       flags,
	}

	unVoteCmd := cli.Command{
		Action:      unVote,
		Name:        "unvote",
		Usage:       "Sends a un-vote for a candidate",
		Description: `Sends a un-vote for a candidate`,
		Flags:       flags,
	}

	resignCmd := cli.Command{
		Action:      resign,
		Name:        "resign",
		Usage:       "resign a candidate, only called by owner of that candidate",
		Description: `resign a candidate, only called by owner of that candidate`,
		Flags:       flags,
	}

	registerCmd := cli.Command{
		Action:      register,
		Name:        "register",
		Usage:       "register a new candidate",
		Description: `register a new candidate`,
		Flags:       flags,
	}

	return []cli.Command{voteCmd, unVoteCmd, resignCmd, registerCmd}
}
