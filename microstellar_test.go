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
	log.Printf("Private seed: %s, Public address: %s", pair.Seed, pair.Address)

	// Fund the account with 1 lumen from an existing account.
	ms.FundAccount(pair.Address, "S6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", "1")

	// Fund an account on the test network with Friendbot.
	FundWithFriendBot(pair.Address)

	// Now load account details from ledger.
	account, _ := ms.LoadAccount(pair.Address)
	log.Printf("Native Balance: %v XLM", account.GetNativeBalance())

	// Pay someone 3 lumens.
	ms.PayNative("S6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", "GAUYTZ24ATLEBIV63MXMPOPQO2T6NHI6TQYEXRTFYXWYZ3JOCVO6UYUM", "3")

	// Pay someone 1 USD issued by an anchor.
	USD := NewAsset("USD", "S6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", Credit4Type)
	ms.Pay("S6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", "GAUYTZ24ATLEBIV63MXMPOPQO2T6NHI6TQYEXRTFYXWYZ3JOCVO6UYUM", USD, "3")

	// Create a trust line to a credit asset with a limit of 1000.
	ms.CreateTrustLine("S4H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", USD, "10000")

	// Check balance.
	account, _ = ms.LoadAccount("GAUYTZ24ATLEBIV63MXMPOPQO2T6NHI6TQYEXRTFYXWYZ3JOCVO6UYUM")
	log.Printf("USD Balance: %v USD", account.GetBalance(USD))
	log.Printf("Native Balance: %v XLM", account.GetNativeBalance())

	// What's their home domain?
	log.Printf("Home domain: %s", account.HomeDomain)

	// Who are the signers on the account?
	for i, s := range account.Signers {
		log.Printf("Signer %d (weight: %v): %v", i, s.PublicKey, s.Weight)
	}
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

// This example makes a native asset (lumens) payment from one account to another.
func ExampleMicroStellar_PayNative() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Pay 1 XLM to targetAddress
	err := ms.PayNative("sourceSeed", "targetAddress", "1")

	if err != nil {
		log.Fatalf("PayNative: %v", ErrorString(err))
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example makes a payment of a credit asset from one account to another.
func ExampleMicroStellar_Pay() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Pay 1 USD to targetAddress
	err := ms.Pay("sourceSeed", "targetAddress", USD, "1")

	if err != nil {
		log.Fatalf("Pay: %v", ErrorString(err))
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example creates a trust line to a credit asset.
func ExampleMicroStellar_CreateTrustLine() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Create a trustline to the custom asset with no limit
	err := ms.CreateTrustLine("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK", USD, "")

	if err != nil {
		log.Fatalf("CreateTrustLine: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example removes a trust line to a credit asset.
func ExampleMicroStellar_RemoveTrustLine() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Remove the trustline (if exists)
	err := ms.RemoveTrustLine("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK", USD)

	if err != nil {
		log.Fatalf("RemoveTrustLine: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example sets the weight of the accounts primary signer (the master weight) to
// zero. This effectively kills the account.
func ExampleMicroStellar_SetMasterWeight() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Set master weight to zero.
	err := ms.SetMasterWeight("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK", 0)

	if err != nil {
		log.Fatalf("SetMasterWeight: %v", err)
	}

	// Load the account and check its master weight
	account, err := ms.LoadAccount("GCCRUJJGPYWKQWM5NLAXUCSBCJKO37VVJ74LIZ5AQUKT6KPVCPNAGC4A")

	if err != nil {
		log.Fatalf("LoadAccount: %v", err)
	}

	log.Printf("Master weight: %v", account.GetMasterWeight())
	fmt.Printf("ok")
	// Output: ok
}
