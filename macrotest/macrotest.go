package main

// This file implements an end-to-end integration test for the
// microstellar library.

import (
	"log"
	"strconv"

	"github.com/0xfe/microstellar"
)

// Helper function to remove stupid "if err != nil" checks
func failOnError(i interface{}, err error) interface{} {
	if err != nil {
		log.Fatal(err)
	}

	return i
}

// Helper function to create a new funded stellar account on the testnet
func createFundedAccount(ms *microstellar.MicroStellar, fundSourceSeed string, useFriendBot bool) microstellar.KeyPair {
	// Create a key pair
	keyPair := failOnError(ms.CreateKeyPair()).(*microstellar.KeyPair)
	log.Printf("Created key pair: %v", keyPair)

	if useFriendBot {
		// Try to fund it with friendbot
		log.Printf("Funding new key with friendbot...")
		resp := failOnError(microstellar.FundWithFriendBot(keyPair.Address))
		log.Printf("Friendbot says: %v", resp)
	}

	// Load the account to see if there are funds
	log.Printf("Checking balance on new key...")
	account, err := ms.LoadAccount(keyPair.Address)
	var floatBalance float64 = 0

	if err == nil {
		balance := account.GetNativeBalance()
		log.Printf("Got balance: %v", balance)
		floatBalance = failOnError(strconv.ParseFloat(balance, 64)).(float64)
	}

	if floatBalance == 0 {
		log.Printf("Looks like it's empty. Funding via source account...")
		err := ms.FundAccount(fundSourceSeed, keyPair.Address, "100")
		if err != nil {
			log.Fatalf("Funding failed: %v", microstellar.ErrorString(err))
		}
		log.Printf("Payment sent: 100 lumens")
	}

	return *keyPair
}

// Helper function to show the asset balance of a specific account
func showBalance(ms *microstellar.MicroStellar, asset *microstellar.Asset, name, address string) {
	log.Printf("Balances for %s: %s", name, address)
	account, err := ms.LoadAccount(address)

	if err != nil {
		log.Fatalf("Canl't load balances for %s: %s", name, address)
	}

	log.Print("  Master weight: ", account.GetMasterWeight())
	log.Print("  XLM: ", account.GetNativeBalance())
	log.Print("  USD: ", account.GetBalance(asset))

	for i, s := range account.Signers {
		log.Printf("  signer %d (type: %v, weight: %v): %v", i, s.Type, s.Weight, s.PublicKey)
	}
}

// This method implements the full end-to-end test
func runTest(fundSourceSeed string) {
	ms := microstellar.New("test")

	// Create a key pair
	keyPair1 := createFundedAccount(ms, fundSourceSeed, true)
	keyPair2 := createFundedAccount(ms, keyPair1.Seed, false)
	keyPair3 := createFundedAccount(ms, keyPair1.Seed, false)
	keyPair4 := createFundedAccount(ms, keyPair1.Seed, false)
	keyPair5 := createFundedAccount(ms, keyPair1.Seed, false)

	log.Print("Pair1 (issuer): ", keyPair1)
	log.Print("Pair2 (distributor): ", keyPair2)
	log.Print("Pair3 (customer): ", keyPair3)
	log.Print("Pair4 (signer1): ", keyPair4)
	log.Print("Pair4 (signer2): ", keyPair5)

	log.Printf("Creating new USD asset issued by %s (issuer)...", keyPair1.Address)
	USD := microstellar.NewAsset("USD", keyPair1.Address, microstellar.Credit4Type)

	log.Printf("Creating USD trustline for %s (distributor)...", keyPair2.Address)
	err := ms.CreateTrustLine(keyPair2.Seed, USD, "10000")

	if err != nil {
		log.Fatalf("CreateTrustLine: %v", err)
	}

	log.Print("Issuing USD from issuer to distributor...")
	err = ms.Pay(microstellar.NewPayment(keyPair1.Seed, keyPair2.Address, "5000").WithAsset(USD))

	if err != nil {
		log.Fatalf("Pay: %v", err)
	}

	log.Printf("Creating USD trustline for %s (customer)...", keyPair3.Address)
	err = ms.CreateTrustLine(keyPair3.Seed, USD, "10000")

	if err != nil {
		log.Fatalf("CreateTrustLine: %v", err)
	}

	log.Printf("Adding new signers to %s (distributor)...", keyPair2.Address)
	ms.AddSigner(keyPair2.Seed, keyPair4.Address, 1)
	ms.AddSigner(keyPair2.Seed, keyPair5.Address, 1)

	log.Printf("Killing master weight for %s (distributor)...", keyPair2.Address)
	err = ms.SetMasterWeight(keyPair2.Seed, 0)

	// See signers for key...
	showBalance(ms, USD, "distributor", keyPair2.Address)

	log.Print("Paying USD from distributor to customer (with dead master signer)...")
	err = ms.Pay(
		microstellar.NewPayment(keyPair2.Seed, keyPair3.Address, "5000").
			WithAsset(USD).
			WithMemoText("failed payment"))

	if err != nil {
		log.Printf("Payment correctly failed: %v", microstellar.ErrorString(err))
	} else {
		log.Fatalf("Payment suceeded. This should not have happened.")
	}

	log.Print("Paying USD from distributor to customer (with too many signers)...")
	err = ms.Pay(
		microstellar.NewPayment(keyPair2.Address, keyPair3.Address, "5000").
			WithAsset(USD).
			WithMemoText("real payment").
			WithSigner(keyPair4.Seed).
			WithSigner(keyPair5.Seed))

	if err != nil {
		log.Printf("Payment correctly failed (too many signers): %v", microstellar.ErrorString(err))
	} else {
		log.Fatalf("Payment succeeded. This should not have happened.")
	}

	log.Print("Paying USD from distributor to customer (with correct signers)...")
	err = ms.Pay(
		microstellar.NewPayment(keyPair2.Address, keyPair3.Address, "5000").
			WithAsset(USD).
			WithMemoText("real payment").
			WithSigner(keyPair5.Seed))

	if err != nil {
		log.Fatalf("Payment failed: %v", microstellar.ErrorString(err))
	}

	log.Print("Sending back USD from customer to distributor before removing trustline...")
	err = ms.Pay(microstellar.NewPayment(keyPair3.Seed, keyPair2.Address, "5000").WithAsset(USD).WithMemoText("take it back"))

	if err != nil {
		log.Fatalf("Pay: %v", err)
	}

	log.Printf("Removing USD trustline for %s (customer)...", keyPair3.Address)
	err = ms.RemoveTrustLine(keyPair3.Seed, USD)

	if err != nil {
		log.Fatalf("RemoveTrustLine: %v", err)
	}

	showBalance(ms, USD, "issuer", keyPair1.Address)
	showBalance(ms, USD, "distributor", keyPair2.Address)
	showBalance(ms, USD, "customer", keyPair3.Address)
	showBalance(ms, USD, "signer", keyPair4.Address)
}

func main() {
	runTest("SBW2N5EK5MZTKPQJZ6UYXEMCA63AO3AVUR6U5CUOIDFYCAR2X2IJIZAX")
}
