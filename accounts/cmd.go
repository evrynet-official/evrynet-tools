package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/evrynet-official/evrynet-client/accounts/abi/bind"
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-tools/accounts/depositor"
	zapLog "github.com/evrynet-official/evrynet-tools/lib/log"
	"github.com/evrynet-official/evrynet-tools/lib/node"

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
	senderPkFlag = cli.StringFlag{
		Name:  "senderpk",
		Usage: "The private key of sender",
		Value: "ce900e4057ef7253ce737dccf3979ec4e74a19d595e8cc30c6c5ea92dfdd37f1",
	}
	expectedBalanceFlag = cli.StringFlag{
		Name:  "expectedbalance",
		Usage: "The expected balance of each account (wei)",
		Value: "1000000000000000000",
	}
)

// NewAccountsFlags return flags to generate accounts
func NewAccountsFlags() []cli.Flag {
	return []cli.Flag{NumAccountsFlag, SeedFlag}
}

// NewDepositFlags return flags to generate accounts
func NewDepositFlags() []cli.Flag {
	return []cli.Flag{NumAccountsFlag, SeedFlag, senderPkFlag, expectedBalanceFlag}
}

// CreateAccountsAndDeposit will print created accounts & write to accounts.json file and send token to accounts
func CreateAccountsAndDeposit(ctx *cli.Context) error {
	accounts, err := createAccounts(ctx)
	if err != nil {
		return err
	}
	err = sendToAccounts(ctx, accounts)
	if err != nil {
		return err
	}

	writeAccounts(accounts)
	return nil
}

// CreateAccounts will print created accounts & write to accounts.json file
func CreateAccounts(ctx *cli.Context) error {
	accounts, err := createAccounts(ctx)
	if err != nil {
		return err
	}
	writeAccounts(accounts)
	return nil
}

func sendToAccounts(ctx *cli.Context, accs []*Account) error {
	var (
		senderPk = ctx.String(senderPkFlag.Name)
		amount   = ctx.String(expectedBalanceFlag.Name)
		gasLimit = big.NewInt(1000000).Uint64()
	)
	expectedAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		fmt.Println("failed to parse expected amount", "amount:", amount)
		return errors.New("failed to parse expected amount")
	}
	err := sendEvr(ctx, accs, senderPk, expectedAmount, gasLimit)
	if err != nil {
		fmt.Println("Fail to send token to accounts!", "Err:", err)
		return err
	}
	return nil
}

func createAccounts(ctx *cli.Context) ([]*Account, error) {
	var (
		num  = ctx.Int(NumAccountsFlag.Name)
		seed = ctx.String(SeedFlag.Name)
	)

	// generate accounts
	accs, err := GenerateAccounts(num, seed)
	if err != nil {
		fmt.Println("Fail to generate new account!", "Err:", err)
		return nil, err
	}
	return accs, nil
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

// SendEvr will send evr token
func sendEvr(ctx *cli.Context, accs []*Account, senderPk string, expectedBalance *big.Int, gasLimit uint64) error {
	pk, err := crypto.HexToECDSA(senderPk)
	if err != nil {
		return err
	}

	var (
		wAddrs []common.Address
		opt    = bind.NewKeyedTransactor(pk)
	)

	for i := 0; i < len(accs); i++ {
		wAddrs = append(wAddrs, accs[i].Address)
	}

	opt.Signer = func(signer types.Signer, from common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return types.SignTx(tx, signer, pk)
	}
	zapLogger, _, err := zapLog.NewSugaredLogger(nil)
	if err != nil {
		return err
	}

	evrClient, err := node.NewEvrynetClientFromFlags(ctx)
	if err != nil {
		return err
	}
	dep := depositor.NewDepositor(zapLogger, opt, wAddrs, evrClient, expectedBalance,
		depositor.WithGasLimit(gasLimit),
	)
	return dep.CheckAndDeposit()
}
