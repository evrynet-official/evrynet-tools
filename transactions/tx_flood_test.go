package transactions

import (
	"testing"

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
				RPCEndpoint: "http://0.0.0.0:22001",
			},
		},
		{
			name: "Send random tx with more than 200 account. Not get error: too many connection",
			txFlood: TxFlood{
				NumAcc:      200,
				NumTxPerAcc: 2,
				Seed:        "testnet",
				RPCEndpoint: "http://0.0.0.0:22001",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, tt.txFlood.floodTx())
		})
	}
}
