.PHONY: accounts

accounts:
	go build -v -o ./build/accounts ./cmd/accounts
	@echo "Done building."
	@echo "Run \"./accounts\" to launch account manager."