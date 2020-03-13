package sc

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
	"github.com/Evrynetlabs/evrynet-node/common"
	stakingContracts "github.com/Evrynetlabs/evrynet-node/consensus/staking_contracts"
	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/evrclient"
	"github.com/evrynet-official/evrynet-tools/lib/log"
)

const (
	TestNodeEndpoint = "http://0.0.0.0:22001"
	StakingScAddress = "0x0000000000000000000000000000000000000011"
	AdminPk          = "85af6fd1be0b4314fc00e8da30091541fff1a6a7159032ba9639fea4449e86cc"
	Candidate        = "0x45F8B547A7f16730c0C8961A21b56c31d84DdB49"
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

	stakingScAddr := common.HexToAddress(StakingScAddress)
	contract, err := stakingContracts.NewStakingContracts(stakingScAddr, client)
	if err != nil {
		t.Error("cannot create the instance of staking contract", "staking address", StakingScAddress)
	}

	optTrans := bind.NewKeyedTransactor(senderPk)

	zap, flush, err := log.NewSugaredLogger(nil)
	if err != nil {
		t.Error("cannot create the instance of zap logger", "error", err)
	}

	defer flush()
	fmt.Printf("*****************register for new candidate = %s\n", Candidate)
	contractClient := ContractClient{
		Contract:  contract,
		Client:    client,
		StakingSc: stakingScAddr,
		SenderPk:  senderPk,
		Candidate: candidate,
		Amount:    new(big.Int).SetUint64(0),
		TranOps:   optTrans,
		Logger:    zap,
	}

	candidates1, err := contractClient.GetAllCandidates(nil)
	if err != nil {
		t.Errorf("GetAllCandidates() error = %v", err)
	}
	fmt.Println("Current candidates:")
	printCandidates(candidates1)
	tx, err := contractClient.Register(nil)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	if err = checkTransStatus(client, tx); err != nil {
		t.Error("can not register new candidate", err)
	} else {
		time.Sleep(EpochTime * time.Second)
		candidates2, err := contractClient.GetAllCandidates(nil)
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
	contractClient.Amount = new(big.Int).SetInt64(100)
	stakeData1, err := contractClient.GetCandidateData(nil)
	if err != nil {
		t.Errorf("GetCandidateData() error = %v", err)
	}
	fmt.Printf("current staking before vote is %v\n", stakeData1.LatestTotalStakes.Int64())
	tx, err = contractClient.Vote(nil)
	if err != nil {
		t.Errorf("Vote() error = %v", err)
	}

	if err = checkTransStatus(client, tx); err != nil {
		t.Error("can not vote", err)
	} else {
		stakeData2, err := contractClient.GetCandidateData(nil)
		if err != nil {
			t.Errorf("GetCandidateData() error = %v", err)
		}
		if stakeData2.LatestTotalStakes.Int64() != 100 {
			t.Errorf("Vote is failed, new stakes = %v", stakeData2.LatestTotalStakes.Int64())
		} else {
			fmt.Printf("successful vote, last staking is %v\n", stakeData2.LatestTotalStakes.Int64())
		}
	}
	fmt.Println("***************************************************")

	fmt.Printf("*****************Un-vote for the candidate = %s\n", Candidate)
	contractClient.Amount = new(big.Int).SetInt64(10)
	stakeData1, err = contractClient.GetCandidateData(nil)
	if err != nil {
		t.Errorf("GetCandidateData() error = %v", err)
	}
	fmt.Printf("current staking before un-vote is %v\n", stakeData1.LatestTotalStakes.Int64())
	tx, err = contractClient.UnVote(nil)
	if err != nil {
		t.Errorf("UnVote() error = %v", err)
	}

	if err = checkTransStatus(client, tx); err != nil {
		t.Error("can not un-vote", err)
	} else {
		stakeData2, err := contractClient.GetCandidateData(nil)
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
	candidates1, err = contractClient.GetAllCandidates(nil)
	if err != nil {
		t.Errorf("GetAllCandidates() error = %v", err)
	}
	fmt.Println("Current candidates:")
	printCandidates(candidates1)
	tx, err = contractClient.Resign(nil)
	if err != nil {
		t.Errorf("Resign() error = %v", err)
	}

	if err = checkTransStatus(client, tx); err != nil {
		t.Error("can not resign", err)
	} else {
		time.Sleep(EpochTime * time.Second)
		candidates2, err := contractClient.GetAllCandidates(nil)
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

func checkTransStatus(client *evrclient.Client, tx *types.Transaction) error {
	var err error
	if tx == nil {
		return errors.New("transaction is nil")
	}
	for i := 0; i < 10; i++ {
		var receipt *types.Receipt
		receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
		if err == nil {
			if receipt.Status != uint64(1) {
				return errors.New("transaction's status is failed")
			}
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return err
}
