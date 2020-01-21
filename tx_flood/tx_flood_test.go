package tx_flood

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Evrynetlabs/evrynet-node/evrclient"
	"github.com/evrynet-official/evrynet-tools/accounts"
)

// Notice: you must run this script to deposit test accounts
// ./build/accounts deposit --num 200 --seed testnet --expectedbalance "1000000000000000000" --senderpk "ce900e4057ef7253ce737dccf3979ec4e74a19d595e8cc30c6c5ea92dfdd37f1" --rpcendpoint "http://0.0.0.0:22001"
func TestTxFlood_floodNormalMode(t *testing.T) {
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
				FloodMode:   NormalTxMode,
			},
		},
		{
			name: "Send random tx with more than 200 account. Not get error: too many connection",
			txFlood: TxFlood{
				NumAcc:      200,
				NumTxPerAcc: 2,
				Seed:        "testnet",
				FloodMode:   NormalTxMode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.txFlood.Accounts, err = accounts.GenerateAccounts(tt.txFlood.NumAcc, tt.txFlood.Seed)
			assert.NoError(t, err)
			tt.txFlood.EvrClient, err = evrclient.Dial("http://0.0.0.0:22001")
			assert.NoError(t, err)
			assert.NoError(t, tt.txFlood.Start())
		})
	}
}

// Notice: you must run this script to deposit test accounts
// ./build/accounts deposit --num 200 --seed testnet --expectedbalance "1000000000000000000" --senderpk "ce900e4057ef7253ce737dccf3979ec4e74a19d595e8cc30c6c5ea92dfdd37f1" --rpcendpoint "http://0.0.0.0:22001"
func TestTxFlood_floodSmartContractMode(t *testing.T) {
	tests := []struct {
		name    string
		txFlood TxFlood
	}{
		{
			name: "Send Tx to Smart Contract",
			txFlood: TxFlood{
				NumAcc:      5,
				NumTxPerAcc: 2,
				Seed:        "testnet",
				FloodMode:   SmartContractMode,
			},
		},
		{
			name: "Send TX with more than 200 accounts. Not get error: too many connection",
			txFlood: TxFlood{
				NumAcc:      200,
				NumTxPerAcc: 2,
				Seed:        "testnet",
				FloodMode:   SmartContractMode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.txFlood.Accounts, err = accounts.GenerateAccounts(tt.txFlood.NumAcc, tt.txFlood.Seed)
			assert.NoError(t, err)
			tt.txFlood.EvrClient, err = evrclient.Dial("http://0.0.0.0:22001")
			assert.NoError(t, err)
			assert.NoError(t, tt.txFlood.Start())
		})
	}
}

// Notice: you must run this script to deposit test accounts
// ./build/accounts deposit --num 200 --seed testnet --expectedbalance "1000000000000000000" --senderpk "ce900e4057ef7253ce737dccf3979ec4e74a19d595e8cc30c6c5ea92dfdd37f1" --rpcendpoint "http://0.0.0.0:22001"
func TestTxFlood_floodDefaultMode(t *testing.T) {
	tests := []struct {
		name    string
		txFlood TxFlood
	}{
		{
			name: "Send Tx to Smart Contract",
			txFlood: TxFlood{
				NumAcc:      5,
				NumTxPerAcc: 2,
				Seed:        "testnet",
				FloodMode:   DefaultMode,
			},
		},
		{
			name: "Send TX with more than 200 accounts. Not get error: too many connection",
			txFlood: TxFlood{
				NumAcc:      200,
				NumTxPerAcc: 2,
				Seed:        "testnet",
				FloodMode:   DefaultMode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.txFlood.Accounts, err = accounts.GenerateAccounts(tt.txFlood.NumAcc, tt.txFlood.Seed)
			assert.NoError(t, err)
			tt.txFlood.EvrClient, err = evrclient.Dial("http://0.0.0.0:22001")
			assert.NoError(t, err)
			assert.NoError(t, tt.txFlood.Start())
		})
	}
}
