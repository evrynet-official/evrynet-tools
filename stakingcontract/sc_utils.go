package sc

import (
	"fmt"

	"github.com/Evrynetlabs/evrynet-node/common"
)

// PrintCandidates prints result on console view
func PrintCandidates(candidates []common.Address) {
	for i := 0; i < len(candidates); i++ {
		fmt.Println(candidates[i].Hex())
	}
}
