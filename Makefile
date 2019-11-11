tools:
	go build -v -o account ./cmd/accounts
	@echo "Done building."
	@echo "Run \"./account\" to launch account manager."