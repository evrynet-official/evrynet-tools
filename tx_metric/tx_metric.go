package tx_metric

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"

	"github.com/evrynet-official/evrynet-client/ethclient"
)

type TxMetric struct {
	StartBlockNumber int64
	StartTime        uint64
	StopTime         uint64
	Duration         time.Duration
	EvrClient        *ethclient.Client
}

func (tm *TxMetric) Start() error {
	// Update StartBlockNumber & StartTime if reaching a block exists transactions
	fmt.Println("--- Finding block has Tx ...")
	for ; ; tm.StartBlockNumber++ {
		bl, err := tm.EvrClient.BlockByNumber(context.Background(), big.NewInt(tm.StartBlockNumber))
		if err != nil {
			return err
		}
		fmt.Printf("Found blocknumber %d | Txs: %d\n", tm.StartBlockNumber, bl.Transactions().Len())

		if bl.Transactions().Len() > 0 {
			tm.StartTime = bl.Time()
			tm.StopTime = tm.StartTime + uint64(tm.Duration.Seconds())
			break
		}
	}

	var (
		totalTx              = 0
		totalBlock           int64
		numberOfBlockHasNoTx = 0
	)
	// Scan every block to sum Tx
	fmt.Println("--- Starting calculate TPS ...")
	for i := tm.StartBlockNumber; ; i++ {
		bl, err := tm.EvrClient.BlockByNumber(context.Background(), big.NewInt(i))
		if bl == nil || err != nil {
			return errors.Errorf("Can not get blocknumber %d. Error: %s", i, err)
		}

		if bl.Time() <= tm.StopTime {
			fmt.Printf("Found blocknumber %d | Txs: %d\n", i, bl.Transactions().Len())
			totalTx += bl.Transactions().Len()
			totalBlock = i - tm.StartBlockNumber + 1
			if bl.Transactions().Len() == 0 {
				numberOfBlockHasNoTx++
			}
		} else {
			break
		}
	}

	fmt.Println("------------------------------")
	fmt.Println("Duration:", tm.Duration.Seconds(), "s")
	fmt.Println("Total Tx:", totalTx)
	fmt.Println("Total Blocks:", totalBlock)
	fmt.Println("Total Blocks have 0 Tx:", numberOfBlockHasNoTx)
	fmt.Println("=> AVG Txs/Block:", float64(totalTx)/float64(totalBlock))
	fmt.Println("=> TPS:", float64(totalTx)/tm.Duration.Seconds())
	return nil
}
