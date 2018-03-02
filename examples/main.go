package main

import (
	"log"

	"github.com/0xfe/microstellar"
)

var testSeed = "SBZH3TR7QLXPCQTAXVTWA3VSNDNPUUIK64KZOKBD2HODK7A3AFU5H63J"

func CreateAndFundAccount() {
	// Create a new MicroStellar client connected to the testnet.
	ms := microstellar.New("test")

	// Generate a new random keypair.
	pair, err := ms.CreateKeyPair()

	if err != nil {
		log.Fatalf("CreateKeyPair: %v", err)
	}

	// Display address and key
	log.Printf("Private seed: %s, Public address: %s", pair.Seed, pair.Address)

	// Fund the account with 1 lumen from an existing account.
	err = ms.FundAccount(pair.Address, "SBZH3TR7QLXPCQTAXVTWA3VSNDNPUUIK64KZOKBD2HODK7A3AFU5H63J", "1")

	if err != nil {
		log.Fatalf("FundAccount: %v", err)
	}
}

func GetBalance() {
	// Create a new MicroStellar client connected to the testnet.
	ms := microstellar.New("test")
	address := "GCCRUJJGPYWKQWM5NLAXUCSBCJKO37VVJ74LIZ5AQUKT6KPVCPNAGC4A"

	// Custom USD asset issued by specified issuer
	usdAsset := microstellar.NewAsset("USD",
		"GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", microstellar.Credit4Type)

	// Load account from ledger.
	account, err := ms.LoadAccount(address)

	if err != nil {
		log.Fatalf("LoadAccount: %v", err)
	}

	// See balances
	log.Printf("Native Balance: %v", account.GetNativeBalance())
	log.Printf("USD Balance: %v", account.GetBalance(usdAsset))
}

func main() {
	// CreateAndFundAccount()
	GetBalance()
}
