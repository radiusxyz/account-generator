package main

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	ParseFlag()
	genesisFilePath := GetFlag("genesis")
	accountCount := GetFlag("account-count")
	balance := GetFlag("balance")
	count, err := strconv.Atoi(accountCount)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	accountFilePath := "./accounts"
	privateKeyFilePath := "./private-keys"
	Generate(count, accountFilePath, privateKeyFilePath)

	// genesis.json 파일 로드
	genesisData := loadGenesis(genesisFilePath)
	// accounts 파일에서 계정 로드
	accounts := loadAccounts(accountFilePath)
	// alloc에 계정 추가
	addAccountsToGenesis(genesisData, accounts, balance)
	// 수정된 genesis.json 파일 저장
	saveGenesis(genesisFilePath, genesisData)

	fmt.Println("Updated genesis.json with accounts.")

}

func ParseFlag() {
	flag.String("genesis", "", "genesis file")
	flag.String("account-count", "", "account count")
	flag.String("balance", "", "balance")
	flag.Parse()
}

func GetFlag(paramName string) string {
	return flag.Lookup(paramName).Value.(flag.Getter).Get().(string)
}

func Generate(accountCount int, accountFilePath string, privateKeyFilePath string) {
	accountFile, err := os.Create(accountFilePath)
	if err != nil {
		log.Fatal(err)
	}
	privateKeyFile, err := os.Create(privateKeyFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer accountFile.Close()
	defer privateKeyFile.Close()

	for i := 0; i < accountCount; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hexutil.Encode(privateKeyBytes)
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("Error casting public key to ECDSA")
		}

		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		if i <= accountCount-1 {
			_, err = accountFile.WriteString(fmt.Sprintf("%s\n", address))
		}
		if err != nil {
			log.Fatal(err)
		}

		if i <= accountCount-1 {
			_, err = privateKeyFile.WriteString(fmt.Sprintf("%s\n", privateKeyHex))
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(accountCount, "addresses and private key generated and saved")
}

// genesis.json 파일 로드
func loadGenesis(filePath string) map[string]interface{} {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open genesis file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read genesis file: %v", err)
	}

	var genesis map[string]interface{}
	if err := json.Unmarshal(data, &genesis); err != nil {
		log.Fatalf("Failed to parse genesis file: %v", err)
	}

	return genesis
}

// accounts 파일 로드
func loadAccounts(filePath string) []string {
	var accounts []string
	// 읽을 파일 열기
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// 파일을 한 줄씩 읽기
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text() // 현재 줄의 텍스트 가져오기
		accounts = append(accounts, line)
	}

	// 읽는 중 오류가 발생했는지 확인
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	return accounts
}

// genesis.json의 alloc에 계정 추가
func addAccountsToGenesis(genesis map[string]interface{}, accounts []string, balance string) {
	// alloc을 가져옴
	alloc, ok := genesis["alloc"].(map[string]interface{})
	if !ok {
		log.Fatalf("Invalid alloc format in genesis.json")
	}

	// 각 계정을 alloc에 추가
	for _, account := range accounts {
		alloc[account] = map[string]interface{}{
			"balance": balance, // 기본 잔액
		}
	}
}

// 수정된 genesis.json 파일 저장
func saveGenesis(filePath string, genesis map[string]interface{}) {
	data, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal genesis data: %v", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		log.Fatalf("Failed to save genesis file: %v", err)
	}
}
