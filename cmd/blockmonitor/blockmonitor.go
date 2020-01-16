package main

import (
	"fmt"
	"log"
	"time"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/blockmonitor"
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

	ticker := time.NewTicker(TimeTicker * time.Second)
	for range ticker.C {
		client.checkNodeAlive()
	}
}

func (client *Client) checkNodeAlive() {
	lastBlock, err := client.BlcClient.GetLastBlock()
	if err != nil {
		sendAlert(client, err.Error(), "SOS", false)
		return
	}

	if lastBlock == nil {
		// send alert
		sendAlert(client,
			fmt.Sprintf("[%s] Block is stuck, latest block: %d", time.Now().Format(time.RFC3339), client.BlcClient.LatestBlock),
			"SOS", false)
		client.SendCount++
		return
	}
	if lastBlock.Cmp(client.BlcClient.LatestBlock) <= 0 {
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
