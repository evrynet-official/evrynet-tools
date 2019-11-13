package depositor

import (
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/accounts/abi/bind"
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/ethclient"
	zapLog "github.com/evrynet-official/evrynet-client/log/zap"

	"github.com/stretchr/testify/assert"
)

const (
	RPCEndpoint = "http://0.0.0.0:22001"
	NodePk = "ce900e4057ef7253ce737dccf3979ec4e74a19d595e8cc30c6c5ea92dfdd37f1"
	SenderAddr = "0x3CA5f11792Bad2aA50816726b441FA306DdeAb2f"
	GasLimit = 1000000
)


func TestSendEth(t *testing.T) {
	var (
		adds = []common.Address{
			common.HexToAddress("1"),
			common.HexToAddress("2"),
			common.HexToAddress("3"),
			common.HexToAddress("4"),

		}
		amount = big.NewInt(0)
		chainId = big.NewInt(15)
		opts = &bind.TransactOpts{
			From: common.HexToAddress(SenderAddr),
		}
	)
	spk, err := crypto.HexToECDSA(NodePk)
	assert.NoError(t, err)
	opts.Signer = func(signer types.Signer, from common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return types.SignTx(tx, signer, spk)
	}

	zapLogger, _, err := zapLog.NewSugaredLogger(nil)
	assert.NoError(t, err)

	client, err := ethclient.Dial(RPCEndpoint)
	assert.NoError(t, err)
	s := &ClientAPI{
		evrClient: client,
	}

	depositor := NewDepositor(zapLogger, opts, adds, s, nil, nil, chainId, WithGasLimit(GasLimit))

	for _, addr := range adds {
		t.Run(addr.Hex(), func(t *testing.T) {
			result, err := depositor.sendETH(addr, amount)
			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}