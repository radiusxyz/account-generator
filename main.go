package main

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func Generate(account_count int, accountFilePath string, privateKeyFilePath string) {
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

	for i := 0; i < account_count; i++ {
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

		if i <= account_count-1 {
			_, err = accountFile.WriteString(fmt.Sprintf("%s\n", address))
		}
		if err != nil {
			log.Fatal(err)
		}

		if i <= account_count-1 {
			_, err = privateKeyFile.WriteString(fmt.Sprintf("%s\n", privateKeyHex))
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(account_count, "addresses and private key generated and saved")
}

func main() {
	accountFilePath := "./accounts"
	privateKeyFilePath := "./private-keys"
	genesisFilePath := "./genesis.json"
	Generate(3, accountFilePath, privateKeyFilePath)

	// genesis.json 파일 로드
	genesisData := loadGenesis(genesisFilePath)
	fmt.Println(genesisData)
	fmt.Println("111111-===-----------==-=-===-=-")
	// accounts 파일에서 계정 로드
	accounts := loadAccounts(accountFilePath)
	fmt.Println(accounts)
	fmt.Println("222222-===-----------==-=-===-=-")
	// alloc에 계정 추가
	addAccountsToGenesis(genesisData, accounts)
	fmt.Println(genesisData)
	// 수정된 genesis.json 파일 저장
	saveGenesis(genesisFilePath, genesisData)

	fmt.Println("Updated genesis.json with accounts.")

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
	//file, err := os.Open(filePath)
	//if err != nil {
	//	log.Fatalf("Failed to open accounts file: %v", err)
	//}
	//defer file.Close()
	//
	//data, err := ioutil.ReadAll(file)
	//if err != nil {
	//	log.Fatalf("Failed to read accounts file: %v", err)
	//}
	//
	//// accounts 파일의 각 라인을 슬라이스로 반환
	//var accounts []string
	//for _, line := range string(data) {
	//	if line > 0 {
	//		accounts = append(accounts, string(line))
	//	}
	//}
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
		fmt.Println(line)
	}

	// 읽는 중 오류가 발생했는지 확인
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	return accounts
}

// genesis.json의 alloc에 계정 추가
func addAccountsToGenesis(genesis map[string]interface{}, accounts []string) {
	// alloc을 가져옴
	alloc, ok := genesis["alloc"].(map[string]interface{})
	if !ok {
		log.Fatalf("Invalid alloc format in genesis.json")
	}

	// 각 계정을 alloc에 추가
	for _, account := range accounts {
		alloc[account] = map[string]interface{}{
			"balance": "20000000000000000000000000000", // 기본 잔액
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
