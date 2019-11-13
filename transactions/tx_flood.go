package transactions

import (
	"context"
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
	RPCEndpoint string
	err         chan error
}

func (tf *TxFlood) floodTx() error {
	tf.err = make(chan error)

	accs, err := accounts.GenerateAccounts(tf.NumAcc, tf.Seed)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(tf.NumAcc * tf.NumTxPerAcc)
	for _, acc := range accs {
		for n := 0; n < tf.NumTxPerAcc; n++ {
			go tf.sendTx(&wg, acc, accs)
		}
	}

	go func() {
		wg.Wait()
		close(tf.err)
	}()

	for err := range tf.err {
		if err != nil {
			return err
		}
	}
	return nil
}

func (tf *TxFlood) sendTx(wg *sync.WaitGroup, acc *accounts.Account, accs []*accounts.Account) {
	defer wg.Done()

	ethClient, err := ethclient.Dial(tf.RPCEndpoint)
	if err != nil {
		tf.err <- err
	}

	rand.Seed(time.Now().UnixNano())
	switch rand.Intn(1) {
	case 0: // Send Evr
		nonce, err := ethClient.PendingNonceAt(context.Background(), acc.Address)
		if err != nil {
			tf.err <- err
		}
		gasPrice, err := ethClient.SuggestGasPrice(context.Background())
		if err != nil {
			tf.err <- err
		}

		genesisBlock, err := ethClient.HeaderByNumber(context.Background(), nil)
		if err != nil {
			tf.err <- err
		}

		randAcc := accs[rand.Intn(len(accs))]
		if !reflect.DeepEqual(acc.Address, randAcc.Address) {
			amount := big.NewInt(rand.Int63n(10))
			transaction := types.NewTransaction(nonce, randAcc.Address, amount, genesisBlock.GasLimit, gasPrice, nil)
			transaction, err = types.SignTx(transaction, types.HomesteadSigner{}, randAcc.PriKey)
			if err != nil {
				tf.err <- err
			}

			err = ethClient.SendTransaction(context.Background(), transaction)
			if err != nil {
				tf.err <- err
			}
			fmt.Printf("Address %s send %d EVR to %s \n", acc.Address.Hex(), amount, randAcc.Address.Hex())
		}
	}
}
