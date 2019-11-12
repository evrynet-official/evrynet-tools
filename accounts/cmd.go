package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

var (
	numAccountsFlag = cli.IntFlag{
		Name:  "num",
		Usage: "Number of accounts want to generate",
		Value: 4,
	}
	seedFlag = cli.StringFlag{
		Name:  "seed",
		Usage: "Seed to generate private key account",
		Value: "evrynet",
	}
)

// NewAccountsFlags return flags to generate accounts
func NewAccountsFlags() []cli.Flag {
	return []cli.Flag{numAccountsFlag, seedFlag}
}

// CreateAccounts will print created accounts & write to accounts.json file
func CreateAccounts(ctx *cli.Context) error {
	num := ctx.Int(numAccountsFlag.Name)
	seed := ctx.String(seedFlag.Name)

	// generate accounts
	accs, err := GenerateAccount(num, seed)
	if err != nil {
		fmt.Println("Fail to generate new account!", "Err:", err)
		return err
	}

	writeAccounts(accs)
	return nil
}

func writeAccounts(accs []*Account) {
	type account struct {
		PriKey  string `json:"private_key"`
		PubKey  string `json:"public_key"`
		Address string `json:"address"`
	}

	var updatedAccs []account
	for _, acc := range accs {
		tempAcc := account{
			PriKey:  acc.PrivateKeyStr(),
			PubKey:  acc.PublicKeyStr(),
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
