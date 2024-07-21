# Wallet Balance Checker

This service periodically checks the balance of specified Ethereum wallets and sends alerts to a Slack channel if the balance falls below a defined threshold. The service is built with Go and can be run inside a Docker container.

## Features

- Periodically checks Ethereum wallet balances.
- Sends alerts to a Slack channel if the balance is below the threshold.
- Configurable via a YAML file.
- Runs as a Docker container.

## Prerequisites

- Docker
- Go 1.22.4 (if running locally)

## Configuration

Create a `config.yaml` file with the following structure:

```yaml
slackWebhook: "https://hooks.slack.com/services/your/slack/webhook"
rpcUrl: "https://eth.llamarpc.com"
explorerUrl: "https://etherscan.io/"
intervalInSecs: 600
wallets:
  - name: "My-wallet-1"
    address: "0xC2Ae2e17fE7A18CcA845bf59899C4eAAAAAAAAAA"
    thresholdInEth: 1
  - name: "My-wallet-2"
    address: "0xC2Ae2e17fE7A18CcA845bf59899C4eBBBBBBBBBB"
    thresholdInEth: 0.001
```

## Building and Running
### Using Docker

1. Build the Docker image:

    ```sh
    docker build -t wallet-checker .
    ```
2. Create and update the `config.yaml` file:
    ```sh
    cp sample-config.yaml config.yaml
    nano config.yaml
    ```

3. Run the Docker container:
    ```sh
    docker run -d --name wallet-checker-container -v ./config.yaml:/app/config.yaml wallet-checker
    ```

### Running Locally

1. Install dependencies:
    ```sh
    go mod download
    ```

2. Build the application:
    ```sh
    go build -o wallet-checker
    ```
3. Create and update the `config.yaml` file:
    ```sh
    cp sample-config.yaml config.yaml
    nano config.yaml
    ```
4. Run the application:
    ```sh
    ./wallet-checker
    ```

## Logging
The service logs important events, including balance checks and alerts, with timestamps. Logs are printed to the standard output.

## Slack Alerts
If a wallet’s balance falls below the defined threshold, an alert is sent to the configured Slack channel. The alert includes:

- Wallet name
- Current balance
- Threshold
- Link to the wallet on the explorer

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.

## Contact
For questions or support, please contact [devops@chainsafe.io].