package sc

import (
	"context"
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
	Candidate        = "0x71562b71999873DB5b286dF957af199Ec94617F7"
	EpochTime        = 2 * 40 //seconds
)

func TestContractClient(t *testing.T) {
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

	fmt.Printf("*****************register for new candidate = %s\n", Candidate)
	contractClient := ContractClient{
		Client:    client,
		StakingSc: common.HexToAddress(StakingScAddress),
		SenderPk:  senderPk,
		Candidate: candidate,
		GasLimit:  uint64(8000000),
		Amount:    new(big.Int).SetUint64(0),
	}
	candidates1, err := contractClient.GetAllCandidates()
	if err != nil {
		t.Errorf("GetAllCandidates() error = %v", err)
	}
	fmt.Println("Current candidates:")
	printCandidates(candidates1)
	tx, err := contractClient.Register()
	if err != nil {
		t.Errorf("Register() error = %v", err)
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
		t.Error("can not register new candidate")
	} else {
		time.Sleep(EpochTime * time.Second)
		candidates2, err := contractClient.GetAllCandidates()
		if err != nil {
			t.Errorf("GetAllCandidates() error = %v", err)
		}
		if len(candidates2) == len(candidates1) {
			t.Errorf("Register is failed, new candidates = %vs", len(candidates2))
		} else {
			fmt.Println("successful register a candidate, new candidates:")
			printCandidates(candidates2)
		}
	}
	fmt.Println("***************************************************")

	fmt.Printf("*****************vote for the candidate = %s\n", Candidate)
	contractClient.Amount = new(big.Int).SetInt64(10)
	stakeData1, err := contractClient.GetCandidateData()
	if err != nil {
		t.Errorf("GetCandidateData() error = %v", err)
	}
	fmt.Printf("current staking before vote is %v\n", stakeData1.LatestTotalStakes.Int64())
	tx, err = contractClient.Vote()
	if err != nil {
		t.Errorf("Vote() error = %v", err)
	}

	txSucceed = false
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
		t.Error("can not vote")
	} else {
		stakeData2, err := contractClient.GetCandidateData()
		if err != nil {
			t.Errorf("GetCandidateData() error = %v", err)
		}
		if stakeData2.LatestTotalStakes.Int64() != 10 {
			t.Errorf("Vote is failed, new stakes = %v", stakeData2.LatestTotalStakes.Int64())
		} else {
			fmt.Printf("successful vote, last staking is %v\n", stakeData2.LatestTotalStakes.Int64())
		}
	}
	fmt.Println("***************************************************")

	fmt.Printf("*****************Un-vote for the candidate = %s\n", Candidate)
	contractClient.Amount = new(big.Int).SetInt64(1)
	stakeData1, err = contractClient.GetCandidateData()
	if err != nil {
		t.Errorf("GetCandidateData() error = %v", err)
	}
	fmt.Printf("current staking before un-vote is %v\n", stakeData1.LatestTotalStakes.Int64())
	tx, err = contractClient.UnVote()
	if err != nil {
		t.Errorf("UnVote() error = %v", err)
	}

	txSucceed = false
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
		t.Error("can not un-vote")
	} else {
		stakeData2, err := contractClient.GetCandidateData()
		if err != nil {
			t.Errorf("GetCandidateData() error = %v", err)
		}
		if stakeData2.LatestTotalStakes.Int64() != 9 {
			t.Errorf("Unvote is failed, new stakes = %v", stakeData2.LatestTotalStakes.Int64())
		} else {
			fmt.Printf("successful un-vote, last staking is %v\n", stakeData2.LatestTotalStakes.Int64())
		}
	}
	fmt.Println("***************************************************")

	fmt.Printf("*****************Resign for the candidate = %s\n", Candidate)
	candidates1, err = contractClient.GetAllCandidates()
	if err != nil {
		t.Errorf("GetAllCandidates() error = %v", err)
	}
	fmt.Println("Current candidates:")
	printCandidates(candidates1)
	tx, err = contractClient.Resign()
	if err != nil {
		t.Errorf("Resign() error = %v", err)
	}

	txSucceed = false
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
		t.Error("can not resign")
	} else {
		time.Sleep(EpochTime * time.Second)
		candidates2, err := contractClient.GetAllCandidates()
		if err != nil {
			t.Errorf("GetAllCandidates() error = %v", err)
		}
		if len(candidates2) == len(candidates1) {
			t.Errorf("Resign is failed, new candidates = %vs", len(candidates2))
		} else {
			fmt.Println("successful resign a candidate, new candidates:")
			printCandidates(candidates2)
		}
	}
}

func printCandidates(candidates []common.Address) {
	for i := 0; i < len(candidates); i++ {
		fmt.Println(candidates[i].Hex())
	}

}
