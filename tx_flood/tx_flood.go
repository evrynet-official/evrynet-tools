package tx_flood

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"github.com/evrynet-official/evrynet-client"
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/common/hexutil"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/ethclient"
	"github.com/evrynet-official/evrynet-client/params"

	"github.com/evrynet-official/evrynet-tools/accounts"
)

type TxFlood struct {
	NumAcc        int
	NumTxPerAcc   int
	Seed          string
	FloodMode     FloodMode
	EvrClient     *ethclient.Client
	Accounts      []*accounts.Account
	Continuous    bool
	SleepInterval time.Duration
}

type FloodMode int

const (
	DefaultMode FloodMode = iota
	NormalTxMode
	SmartContractMode
)

var (
	gasPrice = big.NewInt(params.GasPriceConfig)
)

func handleTxErr(errCh chan error) {
	for err := range errCh {
		if err != nil {
			fmt.Printf("failed to send tx, error %s\n", err)
		}
	}
}

func (tf *TxFlood) Start() error {
	var (
		errChan      = make(chan error)
		failed       uint64
		contractAddr = &common.Address{}
	)

	switch tf.FloodMode {
	case DefaultMode, SmartContractMode:
		var err error
		contractAddr, err = tf.prepareNewContract()
		if err != nil {
			return err
		}
	}

	// Start sending tx flood
	var wg sync.WaitGroup
	for _, acc := range tf.Accounts {
		wg.Add(1)
		go func(acc *accounts.Account, contractAddr *common.Address) {
			defer wg.Done()
			nonce, err := tf.EvrClient.PendingNonceAt(context.Background(), acc.Address)
			if err != nil {
				atomic.AddUint64(&failed, uint64(tf.NumTxPerAcc))
				errChan <- err
				return
			}

			manualNonce := big.NewInt(int64(nonce))
			for {
				for n := 0; n < tf.NumTxPerAcc; n++ {
					err := tf.sendTx(acc, manualNonce, contractAddr)
					if err != nil {
						atomic.AddUint64(&failed, uint64(1))
						errChan <- err

					}
				}
				if !tf.Continuous {
					break
				}
				time.Sleep(tf.SleepInterval)
			}
		}(acc, contractAddr)
	}

	go handleTxErr(errChan)

	wg.Wait()
	close(errChan)

	if failed == 0 {
		return nil
	}

	return fmt.Errorf("fail to send %d transactions", failed)
}

func (tf *TxFlood) sendTx(acc *accounts.Account, nonce *big.Int, contractAddr *common.Address) error {
	rand.Seed(time.Now().UnixNano())
	switch tf.FloodMode {
	case DefaultMode:
		switch rand.Intn(2) {
		case 0: // Send Evr
			err := tf.sendNormalTx(acc, nonce)
			if err != nil {
				return err
			}
		case 1: // Send Evr via SC without provider
			err := tf.sendSmartContractTx(acc, nonce, contractAddr)
			if err != nil {
				return err
			}
		}
	case NormalTxMode:
		err := tf.sendNormalTx(acc, nonce)
		if err != nil {
			return err
		}
	case SmartContractMode:
		err := tf.sendSmartContractTx(acc, nonce, contractAddr)
		if err != nil {
			return err
		}
	default:
		return errors.New("not support for this flood mode")
	}

	return nil
}

func (tf *TxFlood) sendNormalTx(acc *accounts.Account, nonce *big.Int) error {
	randAcc := tf.Accounts[rand.Intn(len(tf.Accounts))]
	var (
		estGas uint64 = 30000
		err    error
	)
	if !reflect.DeepEqual(acc.Address, randAcc.Address) {
		amount := big.NewInt(rand.Int63n(10) + 1) // Send at least 1 EVR
		transaction := types.NewTransaction(nonce.Uint64(), randAcc.Address, amount, estGas, gasPrice, nil)
		transaction, err = types.SignTx(transaction, types.HomesteadSigner{}, acc.PriKey)
		if err != nil {
			return err
		}

		err = tf.EvrClient.SendTransaction(context.Background(), transaction)
		if err != nil {
			return errors.Wrapf(err, "failed to send %d EVR from %s nonce %s", amount, acc.Address.Hex(), nonce.String())
		}
		fmt.Printf("Sent %d EVR from %s => %s nonce %s \n", amount, acc.Address.Hex(), randAcc.Address.Hex(), nonce.String())
		nonce = nonce.Add(nonce, common.Big1)

	}
	return nil
}

func (tf *TxFlood) sendSmartContractTx(acc *accounts.Account, nonce *big.Int, contractAddr *common.Address) error {
	randAcc := tf.Accounts[rand.Intn(len(tf.Accounts))]
	if !reflect.DeepEqual(acc.Address, randAcc.Address) {

		var (
			estGas uint64 = 40000
			err    error
		)
		// data to interact with a function of this contract
		dataBytes := []byte("0x3fb5c1cb0000000000000000000000000000000000000000000000000000000000000002")
		tx := types.NewTransaction(nonce.Uint64(), *contractAddr, big.NewInt(0), estGas, gasPrice, dataBytes)
		tx, err = types.SignTx(tx, types.HomesteadSigner{}, acc.PriKey)

		err = tf.EvrClient.SendTransaction(context.Background(), tx)
		if err != nil {
			return errors.Wrapf(err, "failed to send Tx to SC %s from %s nonce %s", contractAddr.Hex(), acc.Address.Hex(), nonce.String())
		}
		nonce = nonce.Add(nonce, common.Big1)
		fmt.Printf("Sent Tx from %s => SC %s\n", acc.Address.Hex(), contractAddr.Hex())
	}
	return nil
}

func (tf *TxFlood) prepareNewContract() (*common.Address, error) {
	acc := tf.Accounts[0]
	nonce, err := tf.EvrClient.PendingNonceAt(context.Background(), acc.Address)
	if err != nil {
		return nil, err
	}

	// payload to create a smart contract
	payload := "0x608060405260d0806100126000396000f30060806040526004361060525763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633fb5c1cb811460545780638381f58a14605d578063f2c9ecd8146081575b005b60526004356093565b348015606857600080fd5b50606f6098565b60408051918252519081900360200190f35b348015608c57600080fd5b50606f609e565b600055565b60005481565b600054905600a165627a7a723058209573e4f95d10c1e123e905d720655593ca5220830db660f0641f3175c1cdb86e0029"
	payLoadBytes, err := hexutil.Decode(payload)
	if err != nil {
		return nil, err
	}

	msg := evrynet.CallMsg{
		From:  acc.Address,
		To:    nil,
		Value: common.Big0,
		Data:  payLoadBytes,
	}
	estGas, err := tf.EvrClient.EstimateGas(context.Background(), msg)
	if err != nil {
		return nil, err
	}
	tx := types.NewContractCreation(nonce, big.NewInt(0), estGas, gasPrice, payLoadBytes)
	tx, err = types.SignTx(tx, types.HomesteadSigner{}, acc.PriKey)

	err = tf.EvrClient.SendTransaction(context.Background(), tx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create SC from %s", acc.Address.Hex())
	}

	// Wait to get SC address
	for i := 0; i < 10; i++ {
		var receipt *types.Receipt
		receipt, err = tf.EvrClient.TransactionReceipt(context.Background(), tx.Hash())
		if err == nil && receipt.Status == uint64(1) {
			return &receipt.ContractAddress, nil
		}
		time.Sleep(1 * time.Second)
	}
	return nil, errors.New("Can not get SC address")
}
