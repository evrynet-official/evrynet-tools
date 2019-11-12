package depositor

import (
	"github.com/evrynet-official/evrynet-client/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/ethclient"
	zapLog "github.com/evrynet-official/evrynet-client/log/zap"
)

const (
	RPCEndpoint = "http://0.0.0.0:22001"
)
func TestSendEth(t *testing.T) {
	var (
		adds = []common.Address{}
		toAccount = common.HexToAddress("0xc1d38df8d2342c84faab9623b2d021466fb2844c")
		ammount = common.Big1
		chainId = big.NewInt(15)
		opts = &bind.TransactOpts{
			From: common.HexToAddress("0xc1d38df8d2342c84faab9623b2d021466fb2844c"),
			Signer: 
		}
	)

	zapLogger, _, err := zapLog.NewSugaredLogger(nil)
	assert.NoError(t, err)

	client, err := ethclient.Dial(RPCEndpoint)
	assert.NoError(t, err)
	s := &ClientAPI{
		evrClient: client,
	}

	despositor := NewDepositor(zapLogger, opts, adds, s, nil, nil, chainId)
	result, err := despositor.sendETH(toAccount, ammount)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}