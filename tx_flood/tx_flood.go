package tx_flood

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/ethclient"

	"github.com/evrynet-official/evrynet-tools/accounts"
)

type TxFlood struct {
	NumAcc      int
	NumTxPerAcc int
	Seed        string
	Ethclient   *ethclient.Client
	Accounts    []*accounts.Account
}

func (tf *TxFlood) Start() error {
	var (
		errChan     = make(chan error)
		complete100 = true
	)

	// Start sending tx flood
	var wg sync.WaitGroup
	for _, acc := range tf.Accounts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < tf.NumTxPerAcc; n++ {
				errChan <- tf.sendTx(tf.Ethclient, acc, tf.Accounts)
			}
		}()
	}

	// Using goroutine for bypass stuck
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Get error & print to know which tx fail
	for err := range errChan {
		if err != nil {
			fmt.Println(err)
			complete100 = false
		}
	}

	if complete100 {
		return nil
	}
	return errors.New("fail to send some transactions")
}

func (tf *TxFlood) sendTx(ethClient *ethclient.Client, acc *accounts.Account, accs []*accounts.Account) error {
	rand.Seed(time.Now().UnixNano())
	switch rand.Intn(1) {
	case 0: // Send Evr
		nonce, err := ethClient.PendingNonceAt(context.Background(), acc.Address)
		if err != nil {
			return err
		}
		gasPrice, err := ethClient.SuggestGasPrice(context.Background())
		if err != nil {
			return err
		}

		genesisBlock, err := ethClient.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return err
		}

		randAcc := accs[rand.Intn(len(accs))]
		if !reflect.DeepEqual(acc.Address, randAcc.Address) {
			amount := big.NewInt(rand.Int63n(10) + 1) // Send at least 1 EVR
			transaction := types.NewTransaction(nonce, randAcc.Address, amount, genesisBlock.GasLimit, gasPrice, nil)
			transaction, err = types.SignTx(transaction, types.HomesteadSigner{}, randAcc.PriKey)
			if err != nil {
				return err
			}

			err = ethClient.SendTransaction(context.Background(), transaction)
			infoTx := fmt.Sprintf("Sent %d EVR from %s => %s", amount, acc.Address.Hex(), randAcc.Address.Hex())
			if err != nil {
				return errors.New(fmt.Sprintf("[x] %s\n\tError: %s", infoTx, err.Error()))
			}
			fmt.Printf("[v] %s\n", infoTx)
		}
	default:
		return errors.New("not support for this type")
	}
	return nil
}
