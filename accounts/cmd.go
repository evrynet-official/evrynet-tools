package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/evrynet-official/evrynet-client/accounts/abi/bind"
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-tools/accounts/depositor"
	"github.com/evrynet-official/evrynet-tools/lib/app"
	zapLog "github.com/evrynet-official/evrynet-tools/lib/log"

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
	isSendTokenFlag = cli.IntFlag{
		Name:  "issendtoken",
		Usage: "The flag to send token for accounts or not 1/0",
		Value: 0,
	}
	nodePkFlag = cli.StringFlag{
		Name:  "nodepk",
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
	return []cli.Flag{NumAccountsFlag, SeedFlag, isSendTokenFlag, nodePkFlag, expectedBalanceFlag}
}

// CreateAccounts will print created accounts & write to accounts.json file
func CreateAccounts(ctx *cli.Context) error {
	var (
		num         = ctx.Int(NumAccountsFlag.Name)
		seed        = ctx.String(SeedFlag.Name)
		isSendtoken = ctx.Int(isSendTokenFlag.Name)
		nodePk      = ctx.String(nodePkFlag.Name)
		amount      = ctx.String(expectedBalanceFlag.Name)
		gasLimit    = big.NewInt(1000000).Uint64()
	)

	expectedAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		fmt.Println("Failed to parse expected amount!", "amount:", amount)
		return nil
	}
	// generate accounts
	accs, err := GenerateAccounts(num, seed)
	if err != nil {
		fmt.Println("Fail to generate new account!", "Err:", err)
		return err
	}

	if isSendtoken == 1 {
		err = sendEvr(ctx, accs, nodePk, expectedAmount, gasLimit)
		if err != nil {
			fmt.Println("Fail to send token to accounts!", "Err:", err)
			return err
		}
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

// SendEvr will send evr token
func sendEvr(ctx *cli.Context, accs []*Account, nodePk string, expectedBalance *big.Int, gasLimit uint64) error {
	pk, err := crypto.HexToECDSA(nodePk)
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

	evrClient, err := app.NewEvrynetClientFromFlags(ctx)
	if err != nil {
		return err
	}
	dep := depositor.NewDepositor(zapLogger, opt, wAddrs, evrClient, expectedBalance,
		depositor.WithGasLimit(gasLimit),
	)
	return dep.CheckAndDeposit()
}
