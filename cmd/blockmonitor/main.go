package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/blockmonitor"
	"github.com/evrynet-official/evrynet-tools/lib/node"
)

type Client struct {
	TeleClient *blockmonitor.Telegram
	BlcClient  *blockmonitor.Blockchain
	SendCount  int
}

func main() {
	app := cli.NewApp()
	app.Name = "blcMonitor"
	app.Usage = "sends messages to telegram when node dont increase blocks"
	app.Version = "0.0.1"
	app.Commands = healthCheckCommand()

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func healthCheckCommand() []cli.Command {
	healthCheckCmd := cli.Command{
		Action:      blcMonitor,
		Name:        "start",
		Usage:       "Alert to telegram when block is stuck",
		Description: `Alert to telegram when block is stuck`,
	}
	healthCheckCmd.Flags = blockmonitor.NewTeleClientFlag()
	healthCheckCmd.Flags = append(healthCheckCmd.Flags, blockmonitor.NewBlcClientFlag()...)
	healthCheckCmd.Flags = append(healthCheckCmd.Flags, node.NewEvrynetNodeFlags()...)

	return []cli.Command{healthCheckCmd}
}
