.PHONY: accounts tx_flood tx_metric blockmonitor stakingcontract stresssc

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

stakingcontract:
	go build -v -o ./build/sc ./cmd/stakingcontract
	@echo "Done building."
	@echo 'Run "./build/sc" to interact with staking contract.'

stresssc:
	go build -v -o ./build/stresssc ./cmd/stress_sc
	@echo "Done building."
	@echo 'Run "./build/stresssc" to stress test for staking contract.'