package sc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Evrynetlabs/evrynet-node"
	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/evrclient"
)

// CheckTransStatus returns stratus of a transaction
func CheckTransStatus(client *evrclient.Client, tx *types.Transaction) error {
	var err error
	if tx == nil {
		return errors.New("transaction is nil")
	}
	receipt, err := WaitForTx(client, tx.Hash())
	if err != nil {
		return err
	}
	if receipt.Status == uint64(0) {
		return errors.New("transaction got status is failed")
	}
	return nil
}

// WaitForTx wait for a transaction is finished
func WaitForTx(client *evrclient.Client, hash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), hash)
		switch err {
		case evrynet.NotFound:
		case nil:
			return receipt, nil
		default:
			return receipt, err
		}
		time.Sleep(1 * time.Second)
	}
}

// PrintCandidates prints result on console view
func PrintCandidates(candidates []common.Address) {
	for i := 0; i < len(candidates); i++ {
		fmt.Println(candidates[i].Hex())
	}
}
