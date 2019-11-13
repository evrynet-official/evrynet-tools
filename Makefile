.PHONY: accounts transactions

accounts:
	go build -v -o ./build/accounts ./cmd/accounts
	@echo "Done building."
	@echo "Run \"./build/accounts\" to launch accounts manager."

transactions:
	go build -v -o ./build/transactions ./cmd/transactions
	@echo "Done building."
	@echo "Run \"./build/transactions\" to launch transactions manager."