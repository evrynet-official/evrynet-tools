package depositor

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/evrynet-official/evrynet-client"
	"github.com/evrynet-official/evrynet-client/accounts/abi/bind"
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/evrynet-official/evrynet-tools/accounts"
)

var (
	checkMiningInterval        = time.Duration(2 * time.Second)
	gasPrice                   = big.NewInt(params.GasPriceConfig)
	estGas              uint64 = 30000
)

const (
	txPerturn = 20
)

// ClientInterface
type ClientInterface interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	NonceAt(background context.Context, addresses common.Address, blockNumber *big.Int) (uint64, error)
	SuggestGasPrice(background context.Context) (*big.Int, error)
	SendTransaction(background context.Context, transaction *types.Transaction) error
	TransactionReceipt(background context.Context, hash common.Hash) (*types.Receipt, error)
	BalanceAt(background context.Context, addresses common.Address, blockNumber *big.Int) (*big.Int, error)
}

//Depositor maintains the balance of list of wallet to be above a min level
type Depositor struct {
	sugar               *zap.SugaredLogger
	opt                 *bind.TransactOpts
	address             common.Address
	walletAddresses     []*accounts.Account
	client              ClientInterface
	gasLimit            uint64
	checkMiningInterval time.Duration
	sendEthHook         func()
	expectBalance       *big.Int
	numWorkers          int
	nCoreAccount        int
}

//Option provide initial behaviour of Depositor
type Option func(*Depositor)

//WithGasLimit return an Option to set gas limit for depositor
func WithGasLimit(gasLimit uint64) Option {
	return func(dp *Depositor) {
		dp.gasLimit = gasLimit
	}
}

//WithCheckMiningInterval return an Option to set mining sleep time for depositor
func WithCheckMiningInterval(duration time.Duration) Option {
	return func(dp *Depositor) {
		dp.checkMiningInterval = duration
	}
}

// WithSendETHHook is the function to be call after the transaction is called.
func WithSendETHHook(fn func()) Option {
	return func(dp *Depositor) {
		dp.sendEthHook = fn
	}
}

// WithNumWorkers return numWorker to call the balance concurrently
func WithNumWorkers(numWorkers int) Option {
	return func(dp *Depositor) {
		dp.numWorkers = numWorkers
	}
}

//NewDepositor returns a depositor
func NewDepositor(sugar *zap.SugaredLogger, opt *bind.TransactOpts, address common.Address, walletAddrs []*accounts.Account, ethClient ClientInterface, exp *big.Int, ncore int, opts ...Option) *Depositor {
	depositor := &Depositor{
		sugar:               sugar,
		opt:                 opt,
		address:             address,
		walletAddresses:     walletAddrs,
		client:              ethClient,
		sendEthHook:         func() {},
		expectBalance:       exp,
		checkMiningInterval: checkMiningInterval,
		nCoreAccount:        ncore,
	}
	for _, opt := range opts {
		opt(depositor)
	}
	return depositor
}

//sendEVR will send and wait for transaction receipt before returning
func (dp *Depositor) sendEvrFromDepositor(to common.Address, amount *big.Int, nonce uint64) (common.Hash, error) {
	var (
		logger = dp.sugar.With("func", "sendEVR", "wallet_addr", to.Hex(), "amount", amount)
	)
	logger.Infow("sending evr...", "nonce", nonce)
	gasPrice, err := dp.client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}
	tx := types.NewTransaction(nonce, to, amount, dp.gasLimit, gasPrice, nil)
	signedTx, err := dp.opt.Signer(types.HomesteadSigner{}, dp.opt.From, tx)
	if err != nil {
		return common.Hash{}, err
	}

	if err = dp.client.SendTransaction(context.Background(), signedTx); err != nil {
		return common.Hash{}, err
	}
	dp.sendEthHook()
	return signedTx.Hash(), nil
}

func (dp *Depositor) waitForTx(hash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := dp.client.TransactionReceipt(context.Background(), hash)
		switch err {
		case evrynet.NotFound:
		case nil:
			//This is only applicable for Byzantine forks
			//if receipt.Status != types.ReceiptStatusSuccessful {
			//	logger.Infow("tx failed", "tx", receipt.TxHash.Hex())
			//	return receipt, fmt.Errorf("tx %s failed", receipt.TxHash.Hex())
			//}
			return receipt, nil
		default:
			return receipt, err
		}
		time.Sleep(dp.checkMiningInterval)
	}
}

func (dp *Depositor) CheckForBalances() (map[common.Address]*big.Int, error) {
	var (
		balances = make(map[common.Address]*big.Int)
		gr       = errgroup.Group{}
		logger   = dp.sugar.With("func", "CheckAndDeposit")
		mu       = &sync.Mutex{}
	)
	batchSize := int(math.Floor(float64(len(dp.walletAddresses)) / float64(dp.numWorkers)))
	for workerIndex := 0; workerIndex <= dp.numWorkers; workerIndex++ {
		from := workerIndex * batchSize
		to := (workerIndex + 1) * batchSize
		if workerIndex == dp.numWorkers {
			to = len(dp.walletAddresses)
		}
		gr.Go(func() error {
			for i := from; i < to; i++ {
				addr := dp.walletAddresses[i].Address
				balance, gErr := dp.client.BalanceAt(context.Background(), addr, nil)
				if gErr != nil {
					logger.Errorw("failed to get account balance", "address", addr.Hex(), "error", gErr)
					return gErr
				}
				mu.Lock()
				balances[addr] = balance
				mu.Unlock()
			}
			return nil
		})
	}
	if err := gr.Wait(); err != nil {
		return balances, err
	}
	return balances, nil
}

//CheckAndDeposit check if any of the wallet address is below minBalance,
// if it is, deposit an amount to wallet to reach the expected Balance
func (dp *Depositor) CheckAndDeposit() error {
	if err := dp.DepositCoreAccounts(); err != nil {
		return err
	}
	fmt.Printf("done depositing core acount \n\n\n\n\n")
	return dp.DepositEnMass()
}

func handleTxErr(errCh chan error) {
	for err := range errCh {
		if err != nil {
			fmt.Printf("failed to send tx, error %s\n", err)
		}
	}
}

func (dp *Depositor) sendEvr(acc *accounts.Account, to *accounts.Account, nonce *big.Int) error {
	var (
		err error
	)
	transaction := types.NewTransaction(nonce.Uint64(), to.Address, dp.expectBalance, estGas, gasPrice, nil)
	transaction, err = types.SignTx(transaction, types.HomesteadSigner{}, acc.PriKey)
	if err != nil {
		return err
	}

	err = dp.client.SendTransaction(context.Background(), transaction)
	if err != nil {
		return errors.Wrapf(err, "failed to send %d EVR from %s nonce %s", dp.expectBalance, acc.Address.Hex(), nonce.String())
	}
	fmt.Printf("Sent %d EVR from %s => %s nonce %s \n", dp.expectBalance, acc.Address.Hex(), to.Address.Hex(), nonce.String())
	nonce = nonce.Add(nonce, common.Big1)
	return nil
}

func (dp *Depositor) DepositEnMass() error {
	var (
		wg      = &sync.WaitGroup{}
		errChan = make(chan error)
		failed  uint64

		success = uint64(dp.nCoreAccount)
	)
	if len(dp.walletAddresses) <= dp.nCoreAccount {
		return nil
	}
	for i := 0; i < dp.nCoreAccount; i++ {
		wg.Add(1)
		go func(acc *accounts.Account, index int) {

			var (
				txPerCoreAccount = len(dp.walletAddresses)/dp.nCoreAccount - 1
				from             = dp.nCoreAccount + (index)*txPerCoreAccount
				to               = dp.nCoreAccount + (index+1)*txPerCoreAccount
			)
			if to > len(dp.walletAddresses) {
				to = len(dp.walletAddresses)
			}

			//last core account will have to send to all the rest of the account
			if index == dp.nCoreAccount-1 && to < len(dp.walletAddresses) {
				to = len(dp.walletAddresses)
			}
			fmt.Printf("prepare to send from account %d to account[%d-%d]\n", index, from, to)

			defer wg.Done()
			nonce, err := dp.client.PendingNonceAt(context.Background(), acc.Address)
			if err != nil {
				atomic.AddUint64(&failed, uint64(txPerCoreAccount))
				errChan <- err
				return
			}

			manualNonce := big.NewInt(int64(nonce))
		batchLoop:
			for {
				thisBatchEnd := from + txPerturn - 1
				if thisBatchEnd > to {
					thisBatchEnd = to
				}
				for x := 0; x < txPerturn; x++ {

					if err := dp.sendEvr(acc, dp.walletAddresses[from], manualNonce); err != nil {
						atomic.AddUint64(&failed, uint64(1))
						errChan <- err
					} else {
						atomic.AddUint64(&success, uint64(1))
					}
					from += 1
					if from >= to {
						break batchLoop
					}
				}
				time.Sleep(1 * time.Second)
			}
		}(dp.walletAddresses[i], i)
	}
	go handleTxErr(errChan)
	wg.Wait()
	close(errChan)
	fmt.Printf("success %d failed %d \n", success, failed)
	if failed == 0 {
		return nil
	}

	return fmt.Errorf("fail to send %d transactions", failed)
}

func (dp *Depositor) DepositCoreAccounts() error {
	var (
		logger = dp.sugar.With("func", "CheckAndDeposit")
		gr     = errgroup.Group{}
		upto   = dp.nCoreAccount
	)
	if upto > len(dp.walletAddresses) {
		upto = len(dp.walletAddresses)
	}
	nonce, err := dp.client.PendingNonceAt(context.Background(), dp.address)
	if err != nil {
		return err
	}
	logger.Info("get nonce successfully", "current_nonce", nonce)
	txsCost := big.NewInt(1).Mul(big.NewInt(int64(estGas)), gasPrice)
	diff := big.NewInt(1).Mul(big.NewInt(1).Add(dp.expectBalance, txsCost), big.NewInt(int64((len(dp.walletAddresses)-dp.nCoreAccount)/dp.nCoreAccount+1)))

	for i := 0; i < upto; i++ {

		addr := dp.walletAddresses[i].Address
		logger := logger.With(
			"address", addr.Hex(),
			"expected_balance", diff.String(),
		)
		logger.Infow("depositing funds from bank", "deposit_amount", diff.String(), "nonce", nonce)
		txHash, err := dp.sendEvrFromDepositor(addr, diff, nonce)
		if err != nil {
			logger.Error("failed to deposit", "error", err)
			return err
		}
		nonce++
		gr.Go(func() error {
			_, wErr := dp.waitForTx(txHash)
			if wErr != nil {
				logger.Error("failed to deposit", "error", wErr)
				return wErr
			}
			logger.Infow("deposited funds to wallet account", "tx", txHash.Hex())
			return nil
		})

	}
	if err := gr.Wait(); err != nil {
		return err
	}
	return nil
}
