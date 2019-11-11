package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/evrynet-official/evrynet-tools/accounts"
	"github.com/urfave/cli"
)

var (
	prepareCommand = cli.Command{
		Action:    createAccounts,
		Name:      "create",
		Usage:     "Create accounts",
		ArgsUsage: "<create accounts>",
		Flags: []cli.Flag{
			NumAccountsFlag,
			SeedFlag,
		},
		Description: `To prepare accounts`,
	}
)

func createAccounts(ctx *cli.Context) error {
	num := ctx.Int(NumAccountsFlag.Name)
	seed := ctx.String(SeedFlag.Name)

	// generate accounts
	var prepareAccs []*accounts.Account
	for i := 0; i < num; i++ {
		acc, err := accounts.GenerateAccount(seed + strconv.Itoa(i))
		if err != nil {
			fmt.Println("Fail to generate new account!", "Err:", err)
			return err
		}
		prepareAccs = append(prepareAccs, acc)
	}

	writeAccounts(prepareAccs)
	return nil
}

func writeAccounts(accs []*accounts.Account) {
	type account struct {
		PriKey  string `json:"private_key"`
		PubKey  string `json:"public_key"`
		Address string `json:"address"`
	}

	var updatedAccs []account
	for _, acc := range accs {
		tempAcc := account{
			PriKey:  acc.PrivateKey(),
			PubKey:  acc.PublicKey(),
			Address: acc.Address.Hex(),
		}
		updatedAccs = append(updatedAccs, tempAcc)
		accMarsal, _ := json.MarshalIndent(tempAcc, "", "\t")
		fmt.Println(string(accMarsal))
	}

	accsMarshal, err := json.MarshalIndent(updatedAccs, "", "\t")
	if err != nil {
		fmt.Println("Failed to json Marshal accounts, err: ", err)
		return
	}
	if err := ioutil.WriteFile("accounts.json", accsMarshal, os.ModePerm); err != nil {
		fmt.Println("Failed to write file, err: ", err)
		return
	}
	return
}
