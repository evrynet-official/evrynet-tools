package tx_metric

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/evrclient"
	"github.com/evrynet-official/evrynet-tools/lib/timeutil"
)

type TxMetric struct {
	StartBlockNumber     uint64
	StartTime            uint64
	StopTime             uint64
	endTime              uint64
	NumBlock             uint64
	Duration             time.Duration
	EvrClient            *evrclient.Client
	mu                   *sync.Mutex
	minuteStats          []int64
	totalTx              int64
	totalBlock           int64
	numberOfBlockHasNoTx int64
}

func (tm *TxMetric) MetricByTime() error {
	// Update StartBlockNumber & StartTime if reaching a block exists transactions
	if err := tm.FirstBlockHasTx(); err != nil {
		return err
	}
	var (
		totalTx                 int64
		totalBlock              int64
		numberOfBlockHasNoTx    int64
		calculateEachMinuteFlag = tm.Duration.Minutes() > 1
		minuteStats             = make([]int64, int(tm.StopTime-tm.StartTime)/60+1)
		endTime                 uint64
	)

	// Scan every block to sum Tx
	fmt.Println("--- Starting calculate TPS ...")
	for i := tm.StartBlockNumber; ; i++ {
		bl, err := tm.EvrClient.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		if bl == nil || err != nil {
			return errors.Errorf("Can not get blocknumber %d. Error: %s", i, err)
		}
		if bl.Time() > tm.StopTime {
			break
		}
		fmt.Printf("Found blocknumber %d at time %s | Txs: %d\n", i, timeutil.TimestampSToTime(bl.Time()).UTC().String(), bl.Transactions().Len())
		numberOfTx := int64(bl.Transactions().Len())
		if calculateEachMinuteFlag {
			fmt.Printf("bl.Time is %d start time is %d index is %d\n", bl.Time(), tm.StartTime, int(bl.Time()-tm.StartTime)/60)
			minuteStats[int(bl.Time()-tm.StartTime)/60] += numberOfTx
		}

		totalTx += numberOfTx
		totalBlock += 1
		endTime = bl.Time()
		if bl.Transactions().Len() == 0 {
			numberOfBlockHasNoTx++
		}
	}

	tm.Report(endTime, totalBlock, totalTx, numberOfBlockHasNoTx, minuteStats)

	return nil
}

func (tm *TxMetric) UpdateMetric(bl *types.Block) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	numberOfTx := int64(bl.Transactions().Len())

	index := int(bl.Time()-tm.StartTime) / 60
	for index >= len(tm.minuteStats) {
		tm.minuteStats = append(tm.minuteStats, 0)
	}
	fmt.Printf("index is %d, len minutesStats is %d", index, len(tm.minuteStats))
	tm.minuteStats[index] += numberOfTx

	tm.totalTx += numberOfTx
	tm.totalBlock += 1
	if tm.endTime < bl.Time() {
		tm.endTime = bl.Time()
	}
	if bl.Transactions().Len() == 0 {
		tm.numberOfBlockHasNoTx++
	}
}

func (tm *TxMetric) GetBlock(i int64) error {
	var (
		bl  *types.Block
		err error
	)
	fmt.Printf("Getting block %d\n", i)
	for attempt := 0; attempt <= 10; attempt++ {
		bl, err = tm.EvrClient.BlockByNumber(context.Background(), big.NewInt(i))
		if bl == nil || err != nil {
			fmt.Printf("Can not get block number %d. Error: %s, attempt: %d\n", i, err, attempt)
		} else {
			fmt.Printf("Found blocknumber %d at time %s | Txs: %d\n", i, timeutil.TimestampSToTime(bl.Time()).UTC().String(), bl.Transactions().Len())
			tm.UpdateMetric(bl)
			return nil
		}
	}
	return err
}
func (tm *TxMetric) MetricByBlock() error {
	if err := tm.FirstBlockHasTx(); err != nil {
		return err
	}
	var (
		wg    = &sync.WaitGroup{}
		batch = 10
		blN   = tm.StartBlockNumber
	)

	// Scan every block to sum Tx
	fmt.Printf("--- Starting calculate TPS from Block %d to Block %d ---  ...\n", tm.StartBlockNumber, tm.StartBlockNumber+tm.NumBlock)
	for {
		for i := 0; i < batch; i++ {
			wg.Add(1)
			go func(blockNum int64) {
				defer wg.Done()
				if err := tm.GetBlock(blockNum); err != nil {
					panic(err)
				}
			}(int64(blN))
			blN++
			if blN > tm.StartBlockNumber+tm.NumBlock {
				break
			}
		}
		wg.Wait()
		fmt.Printf("Done 1 batch from block %d to block %d\n", blN-uint64(batch), blN)
		if blN > tm.StartBlockNumber+tm.NumBlock {
			break
		}
	}
	tm.Report(tm.endTime, tm.totalBlock, tm.totalTx, tm.numberOfBlockHasNoTx, tm.minuteStats)

	return nil
}

func (tm *TxMetric) Report(endTime uint64, totalBlock, totalTx, numberOfBlockHasNoTx int64, minuteStats []int64) {
	fmt.Println("-----------General Stats----------------")
	fmt.Println("Duration:", endTime-tm.StartTime, "s")
	fmt.Println("Total Tx:", totalTx)
	fmt.Println("Total Blocks:", totalBlock)
	fmt.Println("Total Blocks have 0 Tx:", numberOfBlockHasNoTx)
	fmt.Println("=> AVG Txs/Block:", float64(totalTx)/float64(totalBlock))
	fmt.Println("=> TPS:", float64(totalTx)/float64(endTime-tm.StartTime))
	fmt.Println("=> AVG BlockTime:", float64(endTime-tm.StartTime)/float64(totalBlock))
	// Calculate stats for each minute
	if len(minuteStats) > 1 {
		fmt.Println("-----------Detail stats for each minute----------------")
		for index, txs := range minuteStats {
			fmt.Printf("Txs at minute %d: %d, avg TPS in this minute %.4f\n", index+1, txs, float64(txs)/60.0)
		}
	}
}

func (tm *TxMetric) FirstBlockHasTx() error {
	// Update StartBlockNumber & StartTime if reaching a block exists transactions
	fmt.Println("--- Finding block has Tx ...")
	for ; ; tm.StartBlockNumber++ {
		bl, err := tm.EvrClient.BlockByNumber(context.Background(), big.NewInt(int64(tm.StartBlockNumber)))
		if err != nil {
			return err
		}
		fmt.Printf("Found Block %d | Txs: %d\n", tm.StartBlockNumber, bl.Transactions().Len())

		if bl.Transactions().Len() > 0 {
			tm.StartTime = bl.Time()
			tm.StopTime = tm.StartTime + uint64(tm.Duration.Seconds())
			return nil
		}
	}
	return nil
}
