.PHONY: accounts tx_flood tx_metric

accounts:
	go build -v -o ./build/accounts ./cmd/accounts
	@echo "Done building."
	@echo "Run \"./build/accounts\" to launch accounts manager."

tx_flood:
	go build -v -o ./build/tx_flood ./cmd/tx_flood
	@echo "Done building."
	@echo "Run \"./build/tx_flood\" to launch transactions manager."

tx_metric:
	go build -v -o ./build/tx_metric ./cmd/tx_metric
	@echo "Done building."
	@echo "Run \"./build/tx_metric\" to launch transactions metric."

blockmonitor:
	go build -v -o ./build/blockmonitor ./cmd/blockmonitor
	@echo "Done building."
	@echo "Run \"./build/blockmonitor\" to nonitor node."

smartcontract:
	go build -v -o ./build/sc ./cmd/smartcontract
	@echo "Done building."
	@echo 'Run "./build/sc" to interact with staking contract.'