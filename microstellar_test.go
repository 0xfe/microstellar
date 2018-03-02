package microstellar

import (
	"fmt"
	"log"
)

func Example() {
	// Create a new MicroStellar client connected to the testnet.
	ms := New("test")

	// Generate a new random keypair.
	pair, _ := ms.CreateKeyPair()

	// Display address and key
	log.Printf("Private seed: %s, Public address: %s", pair.Seed, pair.Address)

	// Fund the account with 1 lumen from an existing account.
	ms.FundAccount(pair.Address, "S6 ... private key ... 3J", "1")

	// Fund an account on the test network with Friendbot
	FundWithFriendBot(pair.Address)

	// Now load account details from ledger.
	account, _ := ms.LoadAccount(pair.Address)
	log.Printf("Native Balance: %v XLM", account.GetNativeBalance())

	// Pay someone 3 lumens
	ms.PayNative("S-sourceSeed", "G-targetAccount", "3")

	// Pay someone 1 USD issued by an anchor
	USD := NewAsset("USD", "ISSUERACCOUNT", Credit4Type)
	ms.Pay("S-sourceSeed", "G-targetAccount", USD, "3")

	// Check their balance
	account, _ = ms.LoadAccount("G-targetaccount")
	log.Printf("USD Balance: %v USD", account.GetBalance(USD))
}

// This example creates a key pair and displays the private and
// public keys. In stellar-terminology, the private key is typically
// called a "seed", and the publick key. an "address."
func ExampleMicroStellar_CreateKeyPair() {
	ms := New("test")

	// Generate a new random keypair.
	pair, err := ms.CreateKeyPair()

	if err != nil {
		log.Fatalf("CreateKeyPair: %v", err)
	}

	// Display address and key
	log.Printf("Private seed: %s, Public address: %s", pair.Seed, pair.Address)

	fmt.Printf("ok")
	// Output: ok
}

// This example creates a key pair and funds the account with lumens. FundAccount is
// used for the initial funding of the account -- it is what turns a public address
// into an account.
func ExampleMicroStellar_FundAccount() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Generate a new random keypair.
	pair, err := ms.CreateKeyPair()

	if err != nil {
		log.Fatalf("CreateKeyPair: %v", err)
	}

	// Fund the account with 1 lumen from an existing account.
	err = ms.FundAccount(pair.Address, "SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK", "1")

	if err != nil {
		log.Fatalf("FundAccount: %v", ErrorString(err))
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example loads and displays the native and a non-native balance on an account.
func ExampleMicroStellar_LoadAccount_GetBalance() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Load account from ledger.
	account, err := ms.LoadAccount("GCCRUJJGPYWKQWM5NLAXUCSBCJKO37VVJ74LIZ5AQUKT6KPVCPNAGC4A")

	if err != nil {
		log.Fatalf("LoadAccount: %v", err)
	}

	// See balances
	log.Printf("Native Balance: %v", account.GetNativeBalance())
	log.Printf("USD Balance: %v", account.GetBalance(USD))
	fmt.Printf("ok")
	// Output: ok
}
