package sc

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
	"github.com/Evrynetlabs/evrynet-node/common"
	stakingContracts "github.com/Evrynetlabs/evrynet-node/consensus/staking_contracts"
	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/evrclient"
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

// ContractClient returns a struct
type ContractClient struct {
	Contract  *stakingContracts.StakingContracts
	Client    *evrclient.Client
	StakingSc common.Address
	SenderPk  *ecdsa.PrivateKey
	Candidate common.Address
	GasLimit  uint64
	Amount    *big.Int
	TranOps   *bind.TransactOpts
	Logger    *zap.SugaredLogger
}

// NewNewStakingFromFlags returns new instance of ContractClient
func NewNewStakingFromFlags(ctx *cli.Context, logger *zap.SugaredLogger) (*ContractClient, error) {
	var (
		stakingSc      = ctx.String(stakingScFlag.Name)
		senderPkString = ctx.String(senderPkFlag.Name)
		candidate      = ctx.String(candidateFlag.Name)
		amount         = new(big.Int).SetInt64(ctx.Int64(amountFlag.Name))
		gasLimit       = ctx.Uint64(gasLimitFlag.Name)
	)

	if !common.IsHexAddress(stakingSc) {
		return nil, errors.New("the address of staking sc is invalid")
	}
	if !common.IsHexAddress(candidate) {
		return nil, errors.New("the address of candidate is invalid")
	}
	senderPk, err := crypto.HexToECDSA(senderPkString)
	if err != nil {
		return nil, err
	}

	client, err := node.NewEvrynetClientFromFlags(ctx)
	if err != nil {
		return nil, err
	}
	stakeSCAddr := common.HexToAddress(stakingSc)
	contract, err := stakingContracts.NewStakingContracts(stakeSCAddr, client)
	if err != nil {
		return nil, err
	}

	contractClient := &ContractClient{
		Contract:  contract,
		Client:    client,
		StakingSc: stakeSCAddr,
		SenderPk:  senderPk,
		Candidate: common.HexToAddress(candidate),
		GasLimit:  gasLimit,
		Amount:    amount,
		TranOps:   bind.NewKeyedTransactor(senderPk),
		Logger:    logger,
	}
	return contractClient, nil
}

// Vote sends a transaction to vote for a candidate
func (c ContractClient) Vote(optTrans *bind.TransactOpts) (*types.Transaction, error) {
	contract, err := stakingContracts.NewStakingContracts(c.StakingSc, c.Client)
	if err != nil {
		return nil, err
	}
	if optTrans == nil {
		optTrans = c.TranOps
	}
	tx, err := contract.Vote(optTrans, c.Candidate)
	if err != nil {
		return nil, err
	}
	c.Logger.Infow("transaction is sent", "candidate", c.Candidate.Hex())
	return tx, nil
}

// UnVote sends a transaction to un-vote for a candidate
func (c ContractClient) UnVote(optTrans *bind.TransactOpts) (*types.Transaction, error) {
	if optTrans == nil {
		optTrans = c.TranOps
	}
	tx, err := c.Contract.Unvote(optTrans, c.Candidate, c.Amount)
	if err != nil {
		return nil, err
	}
	c.Logger.Infow("transaction is sent", "candidate", c.Candidate.Hex())
	return tx, nil
}

// Resign sends a transaction to re-sign for a candidate
func (c ContractClient) Resign(optTrans *bind.TransactOpts) (*types.Transaction, error) {
	if optTrans == nil {
		optTrans = c.TranOps
	}
	tx, err := c.Contract.Resign(optTrans, c.Candidate)
	if err != nil {
		return nil, err
	}
	c.Logger.Infow("transaction is sent", "candidate", c.Candidate.Hex())
	return tx, nil
}

// Register sends a transaction to register for a candidate
func (c ContractClient) Register(optTrans *bind.TransactOpts) (*types.Transaction, error) {
	if optTrans == nil {
		optTrans = c.TranOps
	}
	tx, err := c.Contract.Register(optTrans, c.Candidate, optTrans.From)
	if err != nil {
		return nil, err
	}

	c.Logger.Infow("transaction is sent", "candidate", c.Candidate.Hex())
	return tx, nil
}

// GetAllCandidates returns list candidate from SC
func (c ContractClient) GetAllCandidates(opts *bind.CallOpts) ([]common.Address, error) {
	response, err := c.Contract.GetListCandidates(opts)
	if err != nil {
		return nil, err
	}
	return response.Candidates, nil
}

// GetVoters returns list voters for a candidate from SC
func (c ContractClient) GetVoters(opts *bind.CallOpts) ([]common.Address, error) {
	response, err := c.Contract.GetVoters(opts, c.Candidate)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetCandidateData returns the data of a candidate form SC
func (c ContractClient) getCandidateData(opts *bind.CallOpts) (bool, common.Address, *big.Int, error) {
	response, err := c.Contract.GetCandidateData(opts, c.Candidate)
	if err != nil {
		return false, common.Address{}, nil, err
	}
	return response.IsActiveCandidate, response.Owner, response.TotalStake, nil
}
