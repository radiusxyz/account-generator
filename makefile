build:
	go build -o ./bin/account-generator

run: build
	./bin/account-generator -genesis=./genesis.json -account-count=10 -balance=10000000000000000000000000000