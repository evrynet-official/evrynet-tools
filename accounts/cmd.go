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
	"github.com/evrynet-official/evrynet-client/ethclient"
	depositor "github.com/evrynet-official/evrynet-tools/accounts/depositor"
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
	expectedAmountFlag = cli.StringFlag{
		Name:  "expectedamount",
		Usage: "The amount sends to accounts (wei)",
		Value: "1000000000000000000",
	}
	rpcEndpointFlag = cli.StringFlag{
		Name:  "rpcendpoint",
		Usage: "RPC endpoint to send request",
		Value: "http://0.0.0.0:22001",
	}
)

// NewAccountsFlags return flags to generate accounts
func NewAccountsFlags() []cli.Flag {
	return []cli.Flag{NumAccountsFlag, SeedFlag, isSendTokenFlag, nodePkFlag, expectedAmountFlag, rpcEndpointFlag}
}

// CreateAccounts will print created accounts & write to accounts.json file
func CreateAccounts(ctx *cli.Context) error {
	var (
		num         = ctx.Int(NumAccountsFlag.Name)
		seed        = ctx.String(SeedFlag.Name)
		isSendtoken = ctx.Int(isSendTokenFlag.Name)
		rpcEndpoint = ctx.String(rpcEndpointFlag.Name)
		nodePk      = ctx.String(nodePkFlag.Name)
		amount      = ctx.String(expectedAmountFlag.Name)
		gasLimit    = big.NewInt(1000000).Uint64()
		chainID     = big.NewInt(15)
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
		err = sendEvr(accs, nodePk, expectedAmount, chainID, gasLimit, rpcEndpoint)
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
func sendEvr(accs []*Account, nodePk string, expectedAmount *big.Int, chainID *big.Int, gasLimit uint64, rpcEndpoint string) error {
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
		return types.SignTx(tx, types.HomesteadSigner{}, pk)
	}
	zapLogger, _, err := zapLog.NewSugaredLogger(nil)
	if err != nil {
		return err
	}

	evrClient, err := ethclient.Dial(rpcEndpoint)
	if err != nil {
		return err
	}
	dep := depositor.NewDepositor(zapLogger, opt, wAddrs, evrClient, expectedAmount, chainID,
		depositor.WithGasLimit(gasLimit),
	)
	return dep.CheckAndDeposit()
}
