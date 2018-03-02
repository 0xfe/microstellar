package main

import (
	"log"

	"github.com/0xfe/microstellar"
)

var testSeed = "SBZH3TR7QLXPCQTAXVTWA3VSNDNPUUIK64KZOKBD2HODK7A3AFU5H63J"

func CreateAndFundAccount() {
	ms := microstellar.New("test")

	pair, err := ms.CreateKeyPair()

	if err != nil {
		log.Fatalf("CreateKeyPair: %v", err)
	}

	log.Printf("Pair: %v", pair)

	err = ms.FundAccount(pair.Address, "SBZH3TR7QLXPCQTAXVTWA3VSNDNPUUIK64KZOKBD2HODK7A3AFU5H63J", "1")

	if err != nil {
		log.Fatalf("FundAccount: %v", err)
	}
}

func main() {
	CreateAndFundAccount()
}
