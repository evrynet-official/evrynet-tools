package main

import (
	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/accounts/depositor"
	"github.com/evrynet-official/evrynet-tools/lib/log"
)

func deposit(ctx *cli.Context) error {
	zap, flush, err := log.NewSugaredLogger(ctx)
	if err != nil {
		return err
	}
	defer flush()
	dp, err := depositor.NewDepositorFromFlag(ctx, zap)
	if err != nil {
		zap.Error("cannot create depositor", "error", err)
	}
	return dp.CheckAndDeposit()
}
