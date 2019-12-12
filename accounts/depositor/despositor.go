package depositor

import (
	"context"
	"math"
	"math/big"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/evrynet-official/evrynet-client"
	"github.com/evrynet-official/evrynet-client/accounts/abi/bind"
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core/types"
)

var (
	checkMiningInterval = time.Duration(2 * time.Second)
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
	walletAddresses     []common.Address
	client              ClientInterface
	gasLimit            uint64
	checkMiningInterval time.Duration
	sendEthHook         func()
	expectBalance       *big.Int
	numWorkers          int
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
func NewDepositor(sugar *zap.SugaredLogger, opt *bind.TransactOpts, address common.Address, walletAddrs []common.Address, ethClient ClientInterface, exp *big.Int, opts ...Option) *Depositor {
	depositor := &Depositor{
		sugar:               sugar,
		opt:                 opt,
		address:             address,
		walletAddresses:     walletAddrs,
		client:              ethClient,
		sendEthHook:         func() {},
		expectBalance:       exp,
		checkMiningInterval: checkMiningInterval,
	}
	for _, opt := range opts {
		opt(depositor)
	}
	return depositor
}

//sendEVR will send and wait for transaction receipt before returning
func (dp *Depositor) sendEVR(to common.Address, amount *big.Int, nonce uint64) (common.Hash, error) {
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

//CheckAndDeposit check if any of the wallet address is below minBalance,
// if it is, deposit an amount to wallet to reach the expected Balance
func (dp *Depositor) CheckAndDeposit() error {
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
				addr := dp.walletAddresses[i]
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
		return err
	}

	return dp.Deposit(balances)
}

func (dp *Depositor) Deposit(balances map[common.Address]*big.Int) error {
	var (
		logger = dp.sugar.With("func", "CheckAndDeposit")
		gr     = errgroup.Group{}
	)

	nonce, err := dp.client.PendingNonceAt(context.Background(), dp.address)
	if err != nil {
		return err
	}
	logger.Info("get nonce successfully", "current_nonce", nonce)
	for addr, bal := range balances {
		logger := logger.With(
			"address", addr.Hex(),
			"balance", bal.String(),
			"expected_balance", dp.expectBalance.String(),
		)
		if bal.Cmp(dp.expectBalance) < 0 {
			diff := big.NewInt(0).Sub(dp.expectBalance, bal)
			logger.Infow("wallet balance is insufficient, depositing funds from bank", "deposit_amount", diff.String(), "nonce", nonce)
			txHash, err := dp.sendEVR(addr, diff, nonce)
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
	}
	if err := gr.Wait(); err != nil {
		return err
	}
	return nil
}
