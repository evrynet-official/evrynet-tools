package main

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/blockmonitor"
)

const (
	// MaxTimes is a max try to sending messages
	MaxTimes = 3 //the times try sends message
)

func blcMonitor(ctx *cli.Context) {
	client := &Client{
		SendCount: 0,
	}

	teleClient, err := blockmonitor.NewTeleClientFromFlag(ctx)
	if err != nil {
		log.Printf("can not init telegram bot %s", err.Error())
		return
	}
	log.Print("Connected to telegram")
	client.TeleClient = teleClient

	blcClient, err := blockmonitor.NewBlcClientFromFlags(ctx)
	if err != nil {
		log.Printf("can not connect to evrynet node %s", err.Error())
		// send SOS
		sendAlert(client, err.Error(), "SOS", false)
		return
	}
	client.BlcClient = blcClient
	log.Print("evrynet client is created")

	ticker := time.NewTicker(client.BlcClient.Duration * time.Second)
	for range ticker.C {
		client.checkNodeAlive()
	}
}

func (client *Client) isBlockChanged(lastBlock *big.Int) bool {
	if lastBlock == nil || lastBlock.Cmp(client.BlcClient.LatestBlock) <= 0 {
		return false
	}
	return true
}

func (client *Client) checkNodeAlive() {
	var (
		blockChanged = false
		lastBlock    *big.Int
		err          error
	)

	// re-try max-times times if can get latest block, may be this reason is the system is loading
	for i := 0; i < MaxTimes; i++ {
		lastBlock, err = client.BlcClient.GetLastBlock()
		if err != nil {
			log.Printf("get latest block: %s at %d times", err.Error(), i+1)
		}
		blockChanged = client.isBlockChanged(lastBlock)
		if blockChanged {
			break
		}
		time.Sleep(3 * time.Second)
	}

	if !blockChanged {
		// send alert
		sendAlert(client,
			fmt.Sprintf("[%s] Block is stuck, latest block: %d", time.Now().Format(time.RFC3339), client.BlcClient.LatestBlock),
			"SOS", false)
		return
	}

	// check if node is ok from failed
	if client.SendCount >= MaxTimes {
		sendAlert(client, fmt.Sprintf("[%s] Node is ok", time.Now().Format(time.RFC3339)), "OK", true)
	}
	client.SendCount = 0
	client.BlcClient.LatestBlock = lastBlock
	log.Printf("Current block is: %d\n", lastBlock)
}

func sendAlert(client *Client, msg string, caption string, forceSend bool) {
	if forceSend {
		// send message not increase counter
		log.Printf("================send msg: %s", msg)
		client.TeleClient.SendMessage(msg, caption)
		return
	}
	if client.SendCount >= MaxTimes {
		return
	}
	log.Printf("================send msg: %s", msg)
	client.TeleClient.SendMessage(msg, caption)
	client.SendCount++
}
