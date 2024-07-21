package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/yaml.v2"
)

type Config struct {
	SlackWebhook   string `yaml:"slackWebhook"`
	RpcUrl         string `yaml:"rpcUrl"`
	ExplorerUrl    string `yaml:"explorerUrl"`
	IntervalInSecs int    `yaml:"intervalInSecs"`
	Wallets        []struct {
		Name           string  `yaml:"name"`
		Address        string  `yaml:"address"`
		ThresholdInEth float64 `yaml:"thresholdInEth"`
	} `yaml:"wallets"`
}

func main() {
	// Create a new logger that writes to stdout with date and time
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Read the YAML file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	// Parse the YAML file
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	// Assign the Slack webhook URL to a variable
	slackWebhookURL := config.SlackWebhook
	logger.Println("Slack Webhook URL:", slackWebhookURL)

	// Log the interval in seconds
	interval := time.Duration(config.IntervalInSecs) * time.Second
	logger.Printf("Interval in Seconds: %d\n", config.IntervalInSecs)

	// Connect to the Ethereum RPC endpoint
	client, err := rpc.Dial(config.RpcUrl)
	if err != nil {
		logger.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Run the wallet balance check in a loop
	for {
		checkWalletBalances(client, config, slackWebhookURL, logger)
		time.Sleep(interval)
	}
}

func checkWalletBalances(client *rpc.Client, config Config, slackWebhookURL string, logger *log.Logger) {
	for _, wallet := range config.Wallets {
		balance, err := getBalance(client, wallet.Address)
		if err != nil {
			logger.Printf("Failed to get balance for wallet %s: %v\n", wallet.Name, err)
			continue
		}
		balanceInEth := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
		logger.Printf("Wallet Name: %s, Address: %s, Balance: %f ETH, Threshold in ETH: %f\n", wallet.Name, wallet.Address, balanceInEth, wallet.ThresholdInEth)

		// Check if the balance is less than the threshold
		if balanceInEth.Cmp(big.NewFloat(wallet.ThresholdInEth)) < 0 {
			explorerLink := fmt.Sprintf("%s/address/%s", config.ExplorerUrl, wallet.Address)
			warningMessage := fmt.Sprintf("*Warning: Low Wallet Balance*\n*Wallet Name:* %s\n*Balance:* %f ETH\n*Threshold:* %f ETH\n*Details:* %s", wallet.Name, balanceInEth, wallet.ThresholdInEth, explorerLink)
			logger.Println(warningMessage)
			sendSlackNotification(slackWebhookURL, warningMessage)
		}
	}
}

func getBalance(client *rpc.Client, address string) (*big.Int, error) {
	var result string
	err := client.CallContext(context.Background(), &result, "eth_getBalance", address, "latest")
	if err != nil {
		return nil, err
	}
	balance := new(big.Int)
	balance.SetString(result[2:], 16) // Convert hex string to big.Int
	return balance, nil
}

func sendSlackNotification(webhookURL, message string) {
	payload := map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"color": "danger", // This sets the color to red
				"text":  message,
			},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal Slack payload: %v\n", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Failed to send Slack notification: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Slack notification failed with status: %s\n", resp.Status)
	}
}
