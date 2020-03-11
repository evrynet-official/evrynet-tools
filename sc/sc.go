package sc

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/urfave/cli"

	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
	"github.com/Evrynetlabs/evrynet-node/common"
	stakingContracts "github.com/Evrynetlabs/evrynet-node/consensus/staking_contracts"
	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/evrclient"
	"github.com/Evrynetlabs/evrynet-node/log"
	"github.com/evrynet-official/evrynet-tools/lib/node"
)

var (
	stakingScFlag = cli.StringFlag{
		Name:  "stakingsc",
		Usage: "address of staking smart-contract (Address in string format: 0x...)",
		Value: "0x2d5bd25efa0ab97aaca4e888c5fbcb4866904e46",
	}
	senderPkFlag = cli.StringFlag{
		Name:  "senderpk",
		Usage: "the private key of admin/ owner/ user",
		Value: "85af6fd1be0b4314fc00e8da30091541fff1a6a7159032ba9639fea4449e86cc",
	}
	candidateFlag = cli.StringFlag{
		Name:  "candidate",
		Usage: "the address of candidate (Address in string format: 0x...)",
		Value: "0x71562b71999873DB5b286dF957af199Ec94617F7",
	}
	gasLimitFlag = cli.Uint64Flag{
		Name:  "gaslimit",
		Usage: "the gaslimit to execute the call to the contract.",
		Value: 8000000,
	}
	amountFlag = cli.Int64Flag{
		Name:  "amount",
		Usage: "the amount.",
		Value: 0,
	}
)

// NewStakingFlag returns flags for Staking contract client
func NewStakingFlag() []cli.Flag {
	return []cli.Flag{stakingScFlag, senderPkFlag, candidateFlag, gasLimitFlag, amountFlag}
}

type ContractClient struct {
	Client    *evrclient.Client
	StakingSc common.Address
	SenderPk  *ecdsa.PrivateKey
	Candidate common.Address
	GasLimit  uint64
	Amount    *big.Int
}

func NewNewStakingFromFlags(ctx *cli.Context) (*ContractClient, error) {
	var (
		stakingSc      = ctx.String(stakingScFlag.Name)
		senderPkString = ctx.String(senderPkFlag.Name)
		candidate      = ctx.String(candidateFlag.Name)
		amount         = ctx.Int64(amountFlag.Name)
		gasLimit       = ctx.Uint64(gasLimitFlag.Name)
	)

	if !common.IsHexAddress(stakingSc) {
		return nil, errors.New("the address of staking sc is invalid")
	}
	if !common.IsHexAddress(candidate) {
		return nil, errors.New("the address of candidate is invalid")
	}
	if senderPkString == "" {
		return nil, errors.New("the private key of admin is invalid")
	}

	senderPk, err := crypto.HexToECDSA(senderPkString)
	if err != nil {
		return nil, err
	}

	client, err := node.NewEvrynetClientFromFlags(ctx)
	if err != nil {
		return nil, err
	}

	contractClient := &ContractClient{
		Client:    client,
		StakingSc: common.HexToAddress(stakingSc),
		SenderPk:  senderPk,
		Candidate: common.HexToAddress(candidate),
		Amount:    new(big.Int).SetInt64(amount),
		GasLimit:  gasLimit,
	}
	return contractClient, nil
}

func (c ContractClient) Vote() (*types.Transaction, error) {
	contract, err := stakingContracts.NewStakingContracts(c.StakingSc, c.Client)
	if err != nil {
		return nil, err
	}

	optTrans := bind.NewKeyedTransactor(c.SenderPk)
	optTrans.GasLimit = c.GasLimit
	optTrans.Value = c.Amount
	tx, err := contract.Vote(optTrans, c.Candidate)
	if err != nil {
		return nil, err
	}
	log.Info("vote for a candidate is succeed", "candidate", c.Candidate)
	return tx, nil
}

func (c ContractClient) GetCandidateData() (*struct {
	IsActiveCandidate bool
	Owner             common.Address
	LatestTotalStakes *big.Int
}, error) {
	contract, err := stakingContracts.NewStakingContracts(c.StakingSc, c.Client)
	if err != nil {
		return nil, err
	}
	opts := new(bind.CallOpts)
	response, err := contract.GetCandidateData(opts, c.Candidate)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c ContractClient) GetAllCandidates() ([]common.Address, error) {
	contract, err := stakingContracts.NewStakingContracts(c.StakingSc, c.Client)
	if err != nil {
		return nil, err
	}
	opts := new(bind.CallOpts)
	response, err := contract.GetAllCandidates(opts)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c ContractClient) UnVote() (*types.Transaction, error) {
	contract, err := stakingContracts.NewStakingContracts(c.StakingSc, c.Client)
	if err != nil {
		return nil, err
	}

	optTrans := bind.NewKeyedTransactor(c.SenderPk)
	optTrans.GasLimit = c.GasLimit
	tx, err := contract.Unvote(optTrans, c.Candidate, c.Amount)
	if err != nil {
		return nil, err
	}
	log.Info("un-vote for a candidate is succeed", "candidate", c.Candidate)
	return tx, nil
}

func (c ContractClient) Resign() (*types.Transaction, error) {
	contract, err := stakingContracts.NewStakingContracts(c.StakingSc, c.Client)
	if err != nil {
		return nil, err
	}

	optTrans := bind.NewKeyedTransactor(c.SenderPk)
	optTrans.GasLimit = c.GasLimit
	tx, err := contract.Resign(optTrans, c.Candidate)
	if err != nil {
		return nil, err
	}
	log.Info("re-sign for a candidate is succeed", "candidate", c.Candidate)
	return tx, nil
}

func (c ContractClient) Register() (*types.Transaction, error) {
	contract, err := stakingContracts.NewStakingContracts(c.StakingSc, c.Client)
	if err != nil {
		return nil, err
	}

	optTrans := bind.NewKeyedTransactor(c.SenderPk)
	optTrans.GasLimit = c.GasLimit
	tx, err := contract.Register(optTrans, c.Candidate, optTrans.From)
	if err != nil {
		return nil, err
	}

	log.Info("register for a candidate is succeed", "candidate", c.Candidate)
	return tx, nil
}