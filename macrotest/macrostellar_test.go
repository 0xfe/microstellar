package macrotest

// This file implements an end-to-end integration test for the
// microstellar library.

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/0xfe/microstellar"
)

func init() {
	// logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
}

func debugf(msg string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"test": "macrotest"}).Infof(msg, args...)
}

// Send unused friendbot funds here
const homeAddress string = "GB7RTQME2RAOPRBDFBICCP3UDLCIJOSP7ZWCW5IL7Z6L4FNVLZMEWX2G"

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
	debugf("Created key pair: %v", keyPair)

	if useFriendBot {
		// Try to fund it with friendbot
		debugf("Funding new key with friendbot...")
		resp, _ := microstellar.FundWithFriendBot(keyPair.Address)
		debugf("Friendbot says: %v", resp)
	}

	// Load the account to see if there are funds
	debugf("Checking balance on new key...")
	account, err := ms.LoadAccount(keyPair.Address)
	var floatBalance float64 = 0

	if err == nil {
		balance := account.GetNativeBalance()
		debugf("Got balance: %v", balance)
		floatBalance = failOnError(strconv.ParseFloat(balance, 64)).(float64)
	}

	if floatBalance == 0 {
		debugf("Looks like it's empty. Funding via source account...")
		err := ms.FundAccount(fundSourceSeed, keyPair.Address, "200", microstellar.Opts().WithMemoText("initial fund"))
		if err != nil {
			log.Fatalf("Funding failed: %v", microstellar.ErrorString(err))
		}
		debugf("Payment sent: 100 lumens")
	} else {
		debugf("Yay! Friendbot paid us. Sending some lumens back to fundSource...")
		err := ms.PayNative(keyPair.Seed, homeAddress, "5000", microstellar.Opts().WithMemoText("friendbot payback"))

		if err != nil {
			log.Fatalf(microstellar.ErrorString(err))
		}
	}

	return *keyPair
}

// Helper function to show the asset balance of a specific account
func showBalance(ms *microstellar.MicroStellar, asset *microstellar.Asset, name, address string) {
	debugf("Balances for %s: %s", name, address)
	account, err := ms.LoadAccount(address)

	if err != nil {
		log.Fatalf("Canl't load balances for %s: %s", name, address)
	}

	log.Print("  Master weight: ", account.GetMasterWeight())
	log.Print("  XLM: ", account.GetNativeBalance())
	log.Print("  USD: ", account.GetBalance(asset))

	for i, s := range account.Signers {
		debugf("  signer %d (type: %v, weight: %v): %v", i, s.Type, s.Weight, s.PublicKey)
	}
}

// TestMicroStellarPayments implement end-to-end payment tests
func TestMicroStellarPayments(t *testing.T) {
	const fundSourceSeed = "SBW2N5EK5MZTKPQJZ6UYXEMCA63AO3AVUR6U5CUOIDFYCAR2X2IJIZAX"

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
	log.Print("Pair5 (signer2): ", keyPair5)

	log.Print("Watching for payments on distributor's ledger...")
	watcher, err := ms.WatchPayments(keyPair2.Address, microstellar.Opts().WithContext(context.Background()))

	if err != nil {
		log.Fatalf("Can't watch ledger: %+v", err)
	}

	paymentsReceived := 0
	go func() {
		for p := range watcher.Ch {
			debugf("  ## WatchPayments ## (distributor) %v: %v%v%v from %v%v", p.Type, p.Amount, p.StartingBalance, p.AssetCode, p.From, p.Account)
			paymentsReceived++
		}

		debugf("  ## WatchPayments ## (distributor) Done -- Error: %v", *watcher.Err)
	}()

	debugf("Creating new USD asset issued by %s (issuer)...", keyPair1.Address)
	USD := microstellar.NewAsset("USD", keyPair1.Address, microstellar.Credit4Type)

	debugf("Creating USD trustline for %s (distributor)...", keyPair2.Address)
	err = ms.CreateTrustLine(keyPair2.Seed, USD, "1000000")

	if err != nil {
		log.Fatalf("CreateTrustLine: %+v", err)
	}

	log.Print("Issuing USD from issuer to distributor...")
	err = ms.Pay(keyPair1.Seed, keyPair2.Address, "500000", USD)

	if err != nil {
		log.Fatalf("Pay: %+v", microstellar.ErrorString(err))
	}

	debugf("Creating USD trustline for %s (customer)...", keyPair3.Address)
	err = ms.CreateTrustLine(keyPair3.Seed, USD, "100000")

	if err != nil {
		log.Fatalf("CreateTrustLine: %+v", err)
	}

	debugf("Adding new signers to %s (distributor)...", keyPair2.Address)
	ms.AddSigner(keyPair2.Seed, keyPair4.Address, 1)
	ms.AddSigner(keyPair2.Seed, keyPair5.Address, 1)

	debugf("Killing master weight for %s (distributor)...", keyPair2.Address)
	err = ms.SetMasterWeight(keyPair2.Seed, 0)

	// See signers for key...
	showBalance(ms, USD, "distributor", keyPair2.Address)

	log.Print("Paying USD from distributor to customer (with dead master signer)...")
	err = ms.Pay(keyPair2.Seed, keyPair3.Address, "5000", USD, microstellar.Opts().WithMemoText("failed payment"))

	if err != nil {
		log.Print("Payment correctly failed.")
	} else {
		log.Fatalf("Payment suceeded. This should not have happened.")
	}

	log.Print("Paying USD from distributor to customer (with too many signers)...")
	err = ms.Pay(keyPair2.Address, keyPair3.Address, "5000", USD,
		microstellar.Opts().
			WithMemoText("real payment").
			WithSigner(keyPair4.Seed).
			WithSigner(keyPair5.Seed))

	if err != nil {
		log.Print("Payment correctly failed (too many signers).")
	} else {
		log.Fatalf("Payment succeeded. This should not have happened.")
	}

	log.Print("Paying USD from distributor to customer (with correct signers)...")
	err = ms.Pay(keyPair2.Address, keyPair3.Address, "5000", USD,
		microstellar.Opts().
			WithMemoText("real payment").
			WithSigner(keyPair5.Seed))

	if err != nil {
		log.Fatalf("Payment failed: %v", microstellar.ErrorString(err))
	}

	debugf("Require a total signing weight of 2 on distributor...")
	err = ms.SetThresholds(keyPair2.Address, 2, 2, 2, microstellar.Opts().WithSigner(keyPair4.Seed))

	if err != nil {
		log.Fatalf("SetThresholds failed: %v", microstellar.ErrorString(err))
	}

	log.Print("Paying USD from distributor to customer (with additional signer)...")
	err = ms.Pay(keyPair2.Address, keyPair3.Address, "5000", USD,
		microstellar.Opts().
			WithMemoText("real payment").
			WithSigner(keyPair4.Seed).
			WithSigner(keyPair5.Seed))

	if err != nil {
		log.Fatalf("Payment failed: %v", microstellar.ErrorString(err))
	}

	// Kill payment watcher
	log.Print("Killing payment watcher...")
	watcher.Done()

	log.Print("Sending back USD from customer to distributor before removing trustline...")
	err = ms.Pay(keyPair3.Seed, keyPair2.Address, "10000", USD,
		microstellar.Opts().WithMemoText("take it back"))

	if err != nil {
		log.Fatalf("Pay: %+v", err)
	}

	debugf("Removing USD trustline for %s (customer)...", keyPair3.Address)
	err = ms.RemoveTrustLine(keyPair3.Seed, USD)

	if err != nil {
		log.Fatalf("RemoveTrustLine: %v", microstellar.ErrorString(err))
	}

	showBalance(ms, USD, "issuer", keyPair1.Address)
	showBalance(ms, USD, "distributor", keyPair2.Address)
	showBalance(ms, USD, "customer", keyPair3.Address)
	showBalance(ms, USD, "signer1", keyPair4.Address)
	showBalance(ms, USD, "signer2", keyPair5.Address)

	debugf("Total payments on distributor's ledger: %d", paymentsReceived)
}

func showOffersAndBalances(ms *microstellar.MicroStellar, asset *microstellar.Asset, name, address string) {
	showBalance(ms, asset, name, address)

	offers, err := ms.LoadOffers(address)

	if err != nil {
		log.Fatalf("LoadOffers: %+v", microstellar.ErrorString(err))
	}

	debugf("Offers by %s: %d avaialable", name, len(offers))
	for i, o := range offers {
		offerJSON, err := json.MarshalIndent(o, "", "  ")
		if err != nil {
			log.Fatalf("MarshalIndent: %v", err)
		}
		debugf("Offer %d:\n%s", i, string(offerJSON))
	}
}

// TestMicroStellarOffers implement end-to-end DEX tests
func TestMicroStellarOffers(t *testing.T) {
	const fundSourceSeed = "SBW2N5EK5MZTKPQJZ6UYXEMCA63AO3AVUR6U5CUOIDFYCAR2X2IJIZAX"

	ms := microstellar.New("test")

	// Create a key pair
	issuer := createFundedAccount(ms, fundSourceSeed, true)
	bob := createFundedAccount(ms, issuer.Seed, false)
	mary := createFundedAccount(ms, issuer.Seed, false)

	log.Print("Pair1 (issuer): ", issuer)
	log.Print("Pair2 (bob): ", bob)
	log.Print("Pair3 (mary): ", mary)

	debugf("Creating new USD asset issued by %s (issuer)...", issuer.Address)
	USD := microstellar.NewAsset("USD", issuer.Address, microstellar.Credit4Type)

	debugf("Creating USD trustline for %s (bob)...", bob.Address)
	err := ms.CreateTrustLine(bob.Seed, USD, "1000000")

	if err != nil {
		log.Fatalf("CreateTrustLine: %+v", err)
	}

	debugf("Creating USD trustline for %s (mary)...", mary.Address)
	err = ms.CreateTrustLine(mary.Seed, USD, "1000000")

	if err != nil {
		log.Fatalf("CreateTrustLine: %+v", err)
	}

	log.Print("Issuing USD to bob...")
	err = ms.Pay(issuer.Seed, bob.Address, "500000", USD)

	if err != nil {
		log.Fatalf("Pay: %+v", microstellar.ErrorString(err))
	}

	log.Print("Issuing USD to mary...")
	err = ms.Pay(issuer.Seed, mary.Address, "500000", USD)

	if err != nil {
		log.Fatalf("Pay: %+v", microstellar.ErrorString(err))
	}

	// Check balances
	showBalance(ms, USD, "bob", bob.Address)
	showBalance(ms, USD, "mary", mary.Address)

	log.Print("Creating offer for Bob to sell 100 USD...")
	err = ms.CreateOffer(bob.Seed, USD, microstellar.NativeAsset, "2", "50")

	if err != nil {
		log.Fatalf("CreateOffer: %+v", microstellar.ErrorString(err))
	}

	// Check balances
	showOffersAndBalances(ms, USD, "bob", bob.Address)

	debugf("Creating offer for Mary to buy Bob's assets...")
	err = ms.CreateOffer(mary.Seed, microstellar.NativeAsset, USD, "0.5", "100")

	if err != nil {
		log.Fatalf("CreateOffer: %+v", microstellar.ErrorString(err))
	}

	debugf("Expecting 0 offers for Mary and Bob...")
	showOffersAndBalances(ms, USD, "bob", bob.Address)
	showOffersAndBalances(ms, USD, "mary", mary.Address)

	debugf("Creating another offer from Mary to sell XLM...")
	err = ms.CreateOffer(mary.Seed, microstellar.NativeAsset, USD, "0.5", "20")

	if err != nil {
		log.Fatalf("CreateOffer: %+v", microstellar.ErrorString(err))
	}

	// Load her offers to get qty and offer ID
	offers, err := ms.LoadOffers(mary.Address)

	if err != nil {
		log.Fatalf("LoadOffers: %+v", microstellar.ErrorString(err))
	}

	if len(offers) != 1 {
		log.Fatalf("wrong number of offers, want %v, got %v", 1, len(offers))
	}

	// Update her offer
	debugf("Updating mary's offer...")
	offerID := fmt.Sprintf("%d", offers[0].ID)
	err = ms.UpdateOffer(mary.Seed, offerID, microstellar.NativeAsset, USD, "0.5", "5")

	if err != nil {
		log.Fatalf("UpdateOffer: %+v", microstellar.ErrorString(err))
	}

	// Load her offers to get amount
	offers, err = ms.LoadOffers(mary.Address)

	if err != nil {
		log.Fatalf("LoadOffers: %+v", microstellar.ErrorString(err))
	}

	if offers[0].Amount != "5.0000000" {
		showOffersAndBalances(ms, USD, "mary", mary.Address)
		log.Fatalf("wrong amount, want %v, got %v", "5.0000000", offers[0].Amount)
	}

	// Delete the offer
	debugf("Deleting mary's offer...")
	err = ms.DeleteOffer(mary.Seed, offerID, microstellar.NativeAsset, USD, "0.5")

	if err != nil {
		log.Fatalf("DeleteOffer: %+v", microstellar.ErrorString(err))
	}

	// Load her offers to get amount
	offers, err = ms.LoadOffers(mary.Address)

	if err != nil {
		log.Fatalf("LoadOffers: %+v", microstellar.ErrorString(err))
	}

	if len(offers) != 0 {
		showOffersAndBalances(ms, USD, "mary", mary.Address)
		log.Fatalf("wrong number of offers, want %v, got %v", 0, len(offers))
	}
}
