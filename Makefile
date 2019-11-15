.PHONY: accounts tx_flood

accounts:
	go build -v -o ./build/accounts ./cmd/accounts
	@echo "Done building."
	@echo "Run \"./build/accounts\" to launch accounts manager."

tx_flood:
	go build -v -o ./build/tx_flood ./cmd/tx_flood
	@echo "Done building."
	@echo "Run \"./build/tx_flood\" to launch transactions manager."