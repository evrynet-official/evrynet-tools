package depositor

import (
	"context"
	"math/big"
	"testing"

	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind/backends"
	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/core"
	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	zapLog "github.com/evrynet-official/evrynet-tools/lib/log"

	"github.com/stretchr/testify/assert"
)

const (
	NodePk       = "ce900e4057ef7253ce737dccf3979ec4e74a19d595e8cc30c6c5ea92dfdd37f1"
	testAddr1    = "0xAFc44e49dB9ba3E43643bc95B27F4A9a4edfFa9D"
	testAddr2    = "0x35E340dACdba43Cd9B05cE0Ea7a1950824b37098"
	GasLimit     = 1000000
	testBal1     = 1000000 //1e6
	testBal2     = 2000000 //2e6
	testExpBal   = 3000000
	testGasLimit = 100000000
)

func TestDepositor(t *testing.T) {
	pk, err := crypto.HexToECDSA(NodePk)
	assert.NoError(t, err)
	opt := bind.NewKeyedTransactor(pk)
	opt.Signer = func(signer types.Signer, from common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return types.SignTx(tx, signer, pk)
	}

	var (
		wAddrs = []common.Address{
			common.HexToAddress(testAddr1),
			common.HexToAddress(testAddr2),
			opt.From,
		}
		genAlloc = core.GenesisAlloc{
			wAddrs[0]: core.GenesisAccount{
				Balance: big.NewInt(testBal1),
			},
			wAddrs[1]: core.GenesisAccount{
				Balance: big.NewInt(testBal2),
			},
			wAddrs[2]: core.GenesisAccount{
				Balance: big.NewInt(testBal2 * 100000),
			},
		}
	)

	zapLogger, _, err := zapLog.NewSugaredLogger(nil)
	sim := backends.NewSimulatedBackend(genAlloc, testGasLimit)
	dep := NewDepositor(zapLogger, opt, wAddrs, sim, big.NewInt(testExpBal),
		WithSendETHHook(sim.Commit),
		WithCheckMiningInterval(0),
		WithGasLimit(GasLimit),
	)
	assert.NoError(t, dep.CheckAndDeposit())
	newBalance, err := dep.client.BalanceAt(context.Background(), wAddrs[0], nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(testExpBal), newBalance.Int64())
}
