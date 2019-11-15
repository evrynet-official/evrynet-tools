package tx_flood

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/evrynet-official/evrynet-client/common"
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
		go func(acc *accounts.Account) {
			defer wg.Done()
			nonce, err := tf.Ethclient.NonceAt(context.Background(), acc.Address, nil)
			if err != nil {
				errChan <- err
			}

			manualNonce := big.NewInt(int64(nonce))
			for n := 0; n < tf.NumTxPerAcc; n++ {
				errChan <- tf.sendTx(acc, manualNonce)
			}
		}(acc)
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

func (tf *TxFlood) sendTx(acc *accounts.Account, nonce *big.Int) error {
	rand.Seed(time.Now().UnixNano())
	switch rand.Intn(1) {
	case 0: // Send Evr
		gasPrice, err := tf.Ethclient.SuggestGasPrice(context.Background())
		if err != nil {
			return err
		}

		genesisBlock, err := tf.Ethclient.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return err
		}

		randAcc := tf.Accounts[rand.Intn(len(tf.Accounts))]
		if !reflect.DeepEqual(acc.Address, randAcc.Address) {
			amount := big.NewInt(rand.Int63n(10) + 1) // Send at least 1 EVR
			transaction := types.NewTransaction(nonce.Uint64(), randAcc.Address, amount, genesisBlock.GasLimit, gasPrice, nil)
			transaction, err = types.SignTx(transaction, types.HomesteadSigner{}, acc.PriKey)
			if err != nil {
				return err
			}

			err = tf.Ethclient.SendTransaction(context.Background(), transaction)
			if err != nil {
				return errors.Wrapf(err, "failed to send %d EVR from %s", amount, acc.Address.Hex())
			}
			nonce = nonce.Add(nonce, common.Big1)
			fmt.Printf("Sent %d EVR from %s => %s\n", amount, acc.Address.Hex(), randAcc.Address.Hex())
		}
	default:
		return errors.New("not support for this type")
	}
	return nil
}
