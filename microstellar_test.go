package microstellar

import (
	"fmt"
	"log"
)

func Example() {
	// Create a new MicroStellar client connected to a mock network.
	ms := New("fake")

	// Generate a new random keypair.
	pair, _ := ms.CreateKeyPair()

	// In stellar, you can create all kinds of asset types -- dollars, houses, kittens. These
	// customized assets are called credit assets.
	//
	// However, the native asset is always lumens (XLM). Lumens are used to pay for transactions
	// on the stellar network, and are used to fund the operations of Stellar.
	//
	// When you first create a key pair, you need to fund it with atleast 0.5 lumens. This
	// is called the "base reserve", and makes the account valid. You can only transact to
	// and from accounts that maintain the base reserve.
	ms.FundAccount("S6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", pair.Address, "1")

	// On the test network, you can ask FriendBot to fund your account. You don't need to buy
	// lumens. (If you do want to buy lumens for the test network, call me!)
	FundWithFriendBot(pair.Address)

	// Now load the account from the ledger and check its balance.
	account, _ := ms.LoadAccount(pair.Address)
	log.Printf("Native Balance: %v XLM", account.GetNativeBalance())

	// Note that we used GetNativeBalance() above, which implies lumens as the asset. You
	// could also do the following.
	log.Printf("Native Balance: %v XLM", account.GetBalance(NativeAsset))

	// Pay your buddy 3 lumens.
	ms.PayNative("S6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA",
		"GAUYTZ24ATLEBIV63MXMPOPQO2T6NHI6TQYEXRTFYXWYZ3JOCVO6UYUM", "3")

	// Alternatively, be explicit about lumens.
	ms.Pay("S6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA",
		"GAUYTZ24ATLEBIV63MXMPOPQO2T6NHI6TQYEXRTFYXWYZ3JOCVO6UYUM", "3", NativeAsset)

	// Create a credit asset called USD issued by anchor GAT5GKDILNY2G6NOBEIX7XMGSPPZD5MCHZ47MGTW4UL6CX55TKIUNN53
	USD := NewAsset("USD", "GAT5GKDILNY2G6NOBEIX7XMGSPPZD5MCHZ47MGTW4UL6CX55TKIUNN53", Credit4Type)

	// Pay your buddy 3 USD and add a memo
	ms.Pay("S6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA",
		"GAUYTZ24ATLEBIV63MXMPOPQO2T6NHI6TQYEXRTFYXWYZ3JOCVO6UYUM",
		"3", USD,
		Opts().WithMemoText("for beer"))

	// Create a trust line to the USD credit asset with a limit of 1000.
	ms.CreateTrustLine("S4H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", USD, "10000")

	// Check your balances.
	account, _ = ms.LoadAccount("GAUYTZ24ATLEBIV63MXMPOPQO2T6NHI6TQYEXRTFYXWYZ3JOCVO6UYUM")
	log.Printf("USD Balance: %v USD", account.GetBalance(USD))
	log.Printf("Native Balance: %v XLM", account.GetNativeBalance())

	// Find your home domain.
	log.Printf("Home domain: %s", account.HomeDomain)

	// Who are the signers on the account?
	for i, s := range account.Signers {
		log.Printf("Signer %d (weight: %v): %v", i, s.PublicKey, s.Weight)
	}

	log.Printf("ok")
}

func Example_multisig() {
	// Create a new MicroStellar client connected to a mock network.
	ms := New("fake")

	// Add two signers to the source account with weight 1 each
	ms.AddSigner(
		"S8H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", // source account
		"G6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", // signer address
		1) // weight

	ms.AddSigner(
		"S8H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", // source account
		"G9H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGB", // signer address
		1) // weight

	// Set the low, medium, and high thresholds of the account. (Here we require a minimum
	// total signing weight of 2 for all operations.)
	ms.SetThresholds("S8H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", 2, 2, 2)

	// Kill the master weight of account, so only the new signers can sign transactions
	ms.SetMasterWeight("S8H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", 0,
		Opts().WithSigner("S2H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA"))

	// Make a payment (and sign with new signers). Note that the first parameter (source) here
	// can be an address instead of a seed (since the seed can't sign anymore.)
	ms.PayNative(
		"G6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", // from
		"GAUYTZ24ATLEBIV63MXMPOPQO2T6NHI6TQYEXRTFYXWYZ3JOCVO6UYUM", // to
		"3", // amount
		Opts().
			WithSigner("S1H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA").
			WithSigner("S2H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA"))

	log.Printf("ok")
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
	err = ms.FundAccount("SD3M3RG4G54JSFIG4RJYPPKTB4G77IPSXKZPTMN5CKAFWNRQP6V24ZDQ", pair.Address, "1")

	if err != nil {
		log.Fatalf("FundAccount: %v", ErrorString(err))
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example loads and displays the native and a non-native balance on an account.
func ExampleMicroStellar_LoadAccount_balance() {
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
	err := ms.PayNative("SAED4QHN3USETFHECASIM2LRI3H4QTVKZK44D2RC27IICZPZQEGXGXFC", "GDS2DXCCTW5VO5A2KCEBHAP3W4XOCJSI2QVHNN63TXVGBUIIW4DI3BCW", "1")

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
	err := ms.Pay("SAED4QHN3USETFHECASIM2LRI3H4QTVKZK44D2RC27IICZPZQEGXGXFC", "GAGTJGMT55IDNTFTF2F553VQBWRBLGTWLU4YOOIFYBR2F6H6S4AEC45E", "1", USD)

	if err != nil {
		log.Fatalf("Pay: %v", ErrorString(err))
	}

	fmt.Printf("ok")
	// Output: ok
}

// Payments with memotext and memoid
func ExampleMicroStellar_Pay_memotext() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GALC5V4UUUICHENN3ZZLQY6UWAC67CMKVXYT4MT7YGQRD6RMXXCAMHP6", Credit4Type)

	// Pay 1 USD to targetAddress and set the memotext field
	err := ms.Pay("SAED4QHN3USETFHECASIM2LRI3H4QTVKZK44D2RC27IICZPZQEGXGXFC", "GAGTJGMT55IDNTFTF2F553VQBWRBLGTWLU4YOOIFYBR2F6H6S4AEC45E", "1", USD, Opts().WithMemoText("boo"))

	if err != nil {
		log.Fatalf("Pay (memotext): %v", ErrorString(err))
	}

	// Pay 1 USD to targetAddress and set the memotext field
	err = ms.Pay("SAED4QHN3USETFHECASIM2LRI3H4QTVKZK44D2RC27IICZPZQEGXGXFC", "GAGTJGMT55IDNTFTF2F553VQBWRBLGTWLU4YOOIFYBR2F6H6S4AEC45E", "1", USD, Opts().WithMemoID(42))

	if err != nil {
		log.Fatalf("Pay (memoid): %v", ErrorString(err))
	}

	fmt.Printf("ok")
	// Output: ok
}

// Makes a multisignature payment
func ExampleMicroStellar_Pay_multisig() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Pay 1 USD to targetAddress and set the memotext field
	err := ms.Pay("SDKORMIXFL2QW2UC3HWJ4GKL4PYFUMDOPEJMGWVQBW4GWJ5W2ZBOGRSZ", "GAGTJGMT55IDNTFTF2F553VQBWRBLGTWLU4YOOIFYBR2F6H6S4AEC45E", "1", USD,
		Opts().WithMemoText("multisig").
			WithSigner("SAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ").
			WithSigner("SBIUIQNMSXTGR4TGZETSQCGBTIF32G2L5D4AML4LFTMTHKM44UABFDMS"))

	if err != nil {
		log.Fatalf("Pay (memotext): %v", ErrorString(err))
	}

	// Pay 1 USD to targetAddress and set the memotext field
	err = ms.Pay("SAED4QHN3USETFHECASIM2LRI3H4QTVKZK44D2RC27IICZPZQEGXGXFC", "GAGTJGMT55IDNTFTF2F553VQBWRBLGTWLU4YOOIFYBR2F6H6S4AEC45E", "1", USD, Opts().WithMemoID(73223))

	if err != nil {
		log.Fatalf("Pay (memoid): %v", ErrorString(err))
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example makes a path payment from XLM to INR via USD and EUR.
func ExampleMicroStellar_Pay_path() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	XLM := NativeAsset

	// Custom USD, EUR, and INR assets issued by Bank of America
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)
	EUR := NewAsset("EUR", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)
	INR := NewAsset("INR", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Pay 5000 INR with XLM, going through USD and EUR. Spend no more than 40 lumens on this
	// transaction.
	err := ms.Pay(
		"SAED4QHN3USETFHECASIM2LRI3H4QTVKZK44D2RC27IICZPZQEGXGXFC", // from
		"GAGTJGMT55IDNTFTF2F553VQBWRBLGTWLU4YOOIFYBR2F6H6S4AEC45E", // to
		"5000", INR, // they receive 5000 INR
		Opts().
			WithAsset(XLM, "40"). // we spend no more than 40 XLM
			Through(USD, EUR))    // go through USD and EUR

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

// This example adds a signer to an account.
func ExampleMicroStellar_AddSigner() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Add signer to account
	err := ms.AddSigner("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK", "GCCRUJJGPYWKQWM5NLAXUCSBCJKO37VVJ74LIZ5AQUKT6KPVCPNAGC4A", 10)

	if err != nil {
		log.Fatalf("AddSigner: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example removes a signer from an account.
func ExampleMicroStellar_RemoveSigner() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Remove signer from account
	err := ms.RemoveSigner("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK", "GCCRUJJGPYWKQWM5NLAXUCSBCJKO37VVJ74LIZ5AQUKT6KPVCPNAGC4A")

	if err != nil {
		log.Fatalf("RemoveSigner: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example sets the signing thresholds for an account
func ExampleMicroStellar_SetThresholds() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Set the low, medium, and high thresholds for an account
	err := ms.SetThresholds("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK", 2, 2, 2)

	if err != nil {
		log.Fatalf("SetThresholds: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example sets the home domain for an account
func ExampleMicroStellar_SetHomeDomain() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Set the home domain to qubit.sh
	err := ms.SetHomeDomain("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK", "qubit.sh")

	if err != nil {
		log.Fatalf("SetHomeDomain: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example sets flags on an issuer's account
func ExampleMicroStellar_SetFlags() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Set the AUTH_REQUIRED and AUTH_REVOCABLE flags on the account.
	err := ms.SetFlags("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK", FlagAuthRequired|FlagAuthRevocable)

	if err != nil {
		log.Fatalf("SetFlags: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example demonstrates multi-op transactions.
func ExampleMicroStellar_Start() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	feeSource := "GAKMTB3D6AOE5HZ3QK726TZG6A22NGN7B46B2UALVYCLLHLOBMUBXZBJ"
	signer := "SDPLQEABOETMI7PPKJZYBHHW2BSA3424CI3V5ZRNN3NP2H7KYQOKY5ST"

	// Start a new multi-op transaction and bill the fee to feeSource. Also provide the
	// seed of the signer with authority to sign all operations.
	ms.Start(feeSource, Opts().WithMemoText("multi-op").WithSigner(signer))

	// Set the home domain to qubit.sh
	err := ms.SetHomeDomain("GAD3LPHSTZHNZOJOPRS7OZ2P74VXFCP5J4QNYIGGHZ246XINHGKPJIQR", "qubit.sh")

	if err != nil {
		log.Fatalf("SetHomeDomain: %v", err)
	}

	// Set the AUTH_REQUIRED and AUTH_REVOCABLE flags on the account.
	err = ms.SetFlags("GAD3LPHSTZHNZOJOPRS7OZ2P74VXFCP5J4QNYIGGHZ246XINHGKPJIQR", FlagAuthRequired|FlagAuthRevocable)

	if err != nil {
		log.Fatalf("SetFlags: %v", err)
	}

	// Sign and submit the transaction to the network
	err = ms.Submit()
	if err != nil {
		log.Fatalf("Submit: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}
