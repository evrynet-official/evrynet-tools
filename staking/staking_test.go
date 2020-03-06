package staking

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/evrclient"
)

const (
	TestNodeEndpoint = "http://52.220.52.16:22001"
	StakingScAddress = "0x2d5bd25efa0ab97aaca4e888c5fbcb4866904e46"
	AdminPk          = "85af6fd1be0b4314fc00e8da30091541fff1a6a7159032ba9639fea4449e86cc"
	CandidatePk      = "8989232d6c283502ae4fc928324d15369a4a973701aee1fcd5792ca2b5fed153"
	Candidate        = "0x29ADC9eFC670F453AF8C17b6bB6181D91fd748c8"
	EpochTime        = 2 * 40 //seconds
)

func TestContractClient_Register(t *testing.T) {
	var (
		candidate = common.HexToAddress(Candidate)
	)
	client, err := evrclient.Dial(TestNodeEndpoint)
	if err != nil {
		panic(err)
	}
	senderPk, err := crypto.HexToECDSA(AdminPk)
	if err != nil {
		t.Error("private key invalid", "private key", senderPk)
	}

	type fields struct {
		Client    *evrclient.Client
		StakingSc common.Address
		SenderPk  *ecdsa.PrivateKey
		Candidate common.Address
		Owner     common.Address
		GasLimit  uint64
		Amount    *big.Int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []common.Address
		wantErr bool
	}{
		{
			name: "test register",
			fields: fields{
				Client:    client,
				StakingSc: common.HexToAddress(StakingScAddress),
				Candidate: candidate,
				Owner:     candidate,
				GasLimit:  uint64(8000000),
				Amount:    new(big.Int).SetUint64(0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ContractClient{
				Client:    tt.fields.Client,
				StakingSc: tt.fields.StakingSc,
				AdminPk:   senderPk,
				Candidate: tt.fields.Candidate,
				Owner:     tt.fields.Owner,
				GasLimit:  tt.fields.GasLimit,
				Amount:    tt.fields.Amount,
			}
			candidates1, err := c.GetAllCandidates()
			if err != nil {
				t.Errorf("GetAllCandidates() error = %v", err)
				return
			}
			fmt.Println("Current candidates:")
			printCandidates(candidates1)
			tx, err := c.Register()
			if err != nil && tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}

			txSucceed := false
			for i := 0; i < 10; i++ {
				var receipt *types.Receipt
				receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
				if err == nil {
					txSucceed = receipt.Status == uint64(1)
					break
				}
				time.Sleep(1 * time.Second)
			}

			if !txSucceed {
				t.Error("Register is not effected")
				return
			}

			time.Sleep(EpochTime * time.Second)
			candidates2, err := c.GetAllCandidates()
			if err != nil {
				t.Errorf("GetAllCandidates() error = %v", err)
				return
			}
			if len(candidates2) == len(candidates1) {
				t.Errorf("Register is not effected, new candidates = %vs", len(candidates2))
				return
			}

			fmt.Println("=========================================")
			fmt.Println("new candidates:")
			printCandidates(candidates2)
		})
	}
}

func TestContractClient_Vote(t *testing.T) {
	client, err := evrclient.Dial(TestNodeEndpoint)
	if err != nil {
		panic(err)
	}

	senderPk, err := crypto.HexToECDSA(CandidatePk)
	if err != nil {
		t.Error("private key invalid", "private key", senderPk)
	}

	type fields struct {
		Client    *evrclient.Client
		StakingSc common.Address
		SenderPk  *ecdsa.PrivateKey
		Candidate common.Address
		Owner     common.Address
		GasLimit  uint64
		Amount    *big.Int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []common.Address
		wantErr bool
	}{
		{
			name: "test vote",
			fields: fields{
				Client:    client,
				StakingSc: common.HexToAddress(StakingScAddress),
				Candidate: common.HexToAddress(Candidate),
				Owner:     common.HexToAddress(Candidate),
				GasLimit:  uint64(8000000),
				Amount:    new(big.Int).SetUint64(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ContractClient{
				Client:    tt.fields.Client,
				StakingSc: tt.fields.StakingSc,
				AdminPk:   senderPk,
				Candidate: tt.fields.Candidate,
				Owner:     tt.fields.Owner,
				GasLimit:  tt.fields.GasLimit,
				Amount:    tt.fields.Amount,
			}
			stakeData1, err := c.GetCandidateData()
			if err != nil {
				t.Errorf("GetListCandidatesWithCurrentData() error = %v", err)
				return
			}
			t.Logf("Current stakes = %v", stakeData1.LatestTotalStakes.Int64())

			tx, err := c.Vote()
			if err != nil && tt.wantErr {
				t.Errorf("Vote() error = %v, wantErr %v", err, tt.wantErr)
			}

			txSucceed := false
			for i := 0; i < 10; i++ {
				var receipt *types.Receipt
				receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
				if err == nil {
					txSucceed = receipt.Status == uint64(1)
					break
				}
				time.Sleep(1 * time.Second)
			}

			if !txSucceed {
				t.Error("Vote is not effected")
				return
			}

			time.Sleep(EpochTime * time.Second)
			stakeData2, err := c.GetCandidateData()
			if err != nil {
				t.Errorf("GetListCandidatesWithCurrentData() error = %v", err)
				return
			}
			if stakeData2.LatestTotalStakes.Int64() <= stakeData1.LatestTotalStakes.Int64() {
				t.Errorf("Register is not effected, new stakes = %vs", stakeData2.LatestTotalStakes.Int64())
			}
		})
	}
}

func TestContractClient_UnVote(t *testing.T) {
	client, err := evrclient.Dial(TestNodeEndpoint)
	if err != nil {
		panic(err)
	}

	senderPk, err := crypto.HexToECDSA(CandidatePk)
	if err != nil {
		t.Error("private key invalid", "private key", senderPk)
	}

	type fields struct {
		Client    *evrclient.Client
		StakingSc common.Address
		SenderPk  *ecdsa.PrivateKey
		Candidate common.Address
		Owner     common.Address
		GasLimit  uint64
		Amount    *big.Int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []common.Address
		wantErr bool
	}{
		{
			name: "test vote",
			fields: fields{
				Client:    client,
				StakingSc: common.HexToAddress(StakingScAddress),
				Candidate: common.HexToAddress(Candidate),
				Owner:     common.HexToAddress(Candidate),
				GasLimit:  uint64(8000000),
				Amount:    new(big.Int).SetUint64(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ContractClient{
				Client:    tt.fields.Client,
				StakingSc: tt.fields.StakingSc,
				AdminPk:   senderPk,
				Candidate: tt.fields.Candidate,
				Owner:     tt.fields.Owner,
				GasLimit:  tt.fields.GasLimit,
				Amount:    tt.fields.Amount,
			}
			stakeData1, err := c.GetCandidateData()
			if err != nil {
				t.Errorf("GetListCandidatesWithCurrentData() error = %v", err)
				return
			}
			t.Logf("Current stakes = %v", stakeData1.LatestTotalStakes.Int64())

			tx, err := c.UnVote()
			if err != nil && tt.wantErr {
				t.Errorf("Vote() error = %v, wantErr %v", err, tt.wantErr)
			}

			txSucceed := false
			for i := 0; i < 10; i++ {
				var receipt *types.Receipt
				receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
				if err == nil {
					txSucceed = receipt.Status == uint64(1)
					break
				}
				time.Sleep(1 * time.Second)
			}

			if !txSucceed {
				t.Error("Vote is not effected")
				return
			}

			time.Sleep(EpochTime * time.Second)
			stakeData2, err := c.GetCandidateData()
			if err != nil {
				t.Errorf("GetListCandidatesWithCurrentData() error = %v", err)
				return
			}
			if stakeData2.LatestTotalStakes.Int64() <= stakeData1.LatestTotalStakes.Int64() {
				t.Errorf("Register is not effected, new stakes = %vs", stakeData2.LatestTotalStakes.Int64())
			}
		})
	}
}

func TestContractClient_Resign(t *testing.T) {
	var (
		candidate = common.HexToAddress(Candidate)
	)
	client, err := evrclient.Dial(TestNodeEndpoint)
	if err != nil {
		panic(err)
	}
	senderPk, err := crypto.HexToECDSA(CandidatePk)
	if err != nil {
		t.Error("private key invalid", "private key", senderPk)
	}

	type fields struct {
		Client    *evrclient.Client
		StakingSc common.Address
		SenderPk  *ecdsa.PrivateKey
		Candidate common.Address
		Owner     common.Address
		GasLimit  uint64
		Amount    *big.Int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []common.Address
		wantErr bool
	}{
		{
			name: "test resign",
			fields: fields{
				Client:    client,
				StakingSc: common.HexToAddress(StakingScAddress),
				Candidate: candidate,
				Owner:     candidate,
				GasLimit:  uint64(8000000),
				Amount:    new(big.Int).SetUint64(0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ContractClient{
				Client:    tt.fields.Client,
				StakingSc: tt.fields.StakingSc,
				AdminPk:   senderPk,
				Candidate: tt.fields.Candidate,
				Owner:     tt.fields.Owner,
				GasLimit:  tt.fields.GasLimit,
				Amount:    tt.fields.Amount,
			}
			candidates1, err := c.GetAllCandidates()
			if err != nil {
				t.Errorf("GetAllCandidates() error = %v", err)
				return
			}
			fmt.Println("Current candidates:")
			printCandidates(candidates1)
			tx, err := c.Resign()
			if err != nil && tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}

			txSucceed := false
			for i := 0; i < 10; i++ {
				var receipt *types.Receipt
				receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
				if err == nil {
					txSucceed = receipt.Status == uint64(1)
					break
				}
				time.Sleep(1 * time.Second)
			}

			if !txSucceed {
				t.Error("Register is not effected")
				return
			}

			time.Sleep(EpochTime * time.Second)
			candidates2, err := c.GetAllCandidates()
			if err != nil {
				t.Errorf("GetAllCandidates() error = %v", err)
				return
			}
			if len(candidates2) == len(candidates1) {
				t.Errorf("Register is not effected, new candidates = %vs", len(candidates2))
				return
			}

			fmt.Println("=========================================")
			fmt.Println("new candidates:")
			printCandidates(candidates2)
		})
	}
}

func printCandidates(candidates []common.Address) {
	for i := 0; i < len(candidates); i++ {
		fmt.Println(candidates[i].Hex())
	}

}
