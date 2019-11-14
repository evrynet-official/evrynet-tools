package tx_flood

import (
	"testing"

	"github.com/evrynet-official/evrynet-client/ethclient"
	"github.com/evrynet-official/evrynet-tools/accounts"
	"github.com/stretchr/testify/assert"
)

func TestTxFlood_floodTx(t *testing.T) {
	tests := []struct {
		name    string
		txFlood TxFlood
	}{
		{
			name: "Send random tx",
			txFlood: TxFlood{
				NumAcc:      5,
				NumTxPerAcc: 2,
				Seed:        "testnet",
			},
		},
		{
			name: "Send random tx with more than 200 account. Not get error: too many connection",
			txFlood: TxFlood{
				NumAcc:      200,
				NumTxPerAcc: 2,
				Seed:        "testnet",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.txFlood.Accounts, _ = accounts.GenerateAccounts(tt.txFlood.NumAcc, tt.txFlood.Seed)
			tt.txFlood.Ethclient, _ = ethclient.Dial("http://0.0.0.0:22001")
			assert.NoError(t, tt.txFlood.Start())
		})
	}
}
