package macrotest

// This file implements an end-to-end integration test for the
// microstellar library.

import (
	"context"
	"log"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/0xfe/microstellar"
)

func init() {
	// logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
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

	log.Print("Watching for transactions on distributor's ledger...")
	watcher, err := ms.WatchPayments(keyPair2.Address, microstellar.Opts().WithContext(context.Background()))
	tWatcher, err := ms.WatchTransactions(keyPair2.Address, microstellar.Opts().WithContext(context.Background()))

	if err != nil {
		log.Fatalf("Can't watch ledger: %+v", err)
	}

	paymentsReceived := 0
	transactionsSeen := 0

	go func() {
		for p := range watcher.Ch {
			debugf("  ## WatchPayments ## (distributor) %v: %v%v%v from %v%v", p.Type, p.Amount, p.StartingBalance, p.AssetCode, p.From, p.Account)
			paymentsReceived++
		}

		debugf("  ## WatchPayments ## (distributor) Done -- Error: %v", *watcher.Err)
	}()

	go func() {
		for t := range tWatcher.Ch {
			debugf("  ## WatchTransactions ## (distributor) %v: (ops: %v) (memo: %v) (fee: %d)", t.ID, t.OperationCount, t.Memo, t.FeePaid)
			transactionsSeen++
		}

		debugf("  ## WatchTransactions ## (distributor) Done -- Error: %v", *tWatcher.Err)
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

	log.Print("Killing watchers...")
	watcher.Done()
	tWatcher.Done()

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
	debugf("Total transactions on distributor's ledger: %d", transactionsSeen)
}

func TestMicroStellarMultiOp(t *testing.T) {
	const fundSourceSeed = "SBW2N5EK5MZTKPQJZ6UYXEMCA63AO3AVUR6U5CUOIDFYCAR2X2IJIZAX"

	ms := microstellar.New("test")

	bob := createFundedAccount(ms, fundSourceSeed, true)
	mary := createFundedAccount(ms, bob.Seed, false)

	log.Print("Pair1 (bob): ", bob)
	log.Print("Pair2 (mary): ", mary)

	debugf("Adding mary as bob's signer...")
	ms.AddSigner(bob.Seed, mary.Address, 1)

	debugf("Starting self-signed multi-op transaction...")
	ms.Start(bob.Seed, microstellar.Opts().WithMemoText("multi-op"))
	ms.SetHomeDomain(bob.Address, "qubit.sh")
	ms.SetFlags(bob.Address, microstellar.FlagAuthRequired)
	ms.PayNative(bob.Address, mary.Address, "1")

	err := ms.Submit()
	if err != nil {
		log.Fatalf("Submit: %v", microstellar.ErrorString(err))
	}

	debugf("Starting multi-op transaction with alternate signer...")
	ms.Start(bob.Address, microstellar.Opts().WithMemoText("multi-op").WithSigner(mary.Seed))
	ms.SetHomeDomain(bob.Address, "qubit.sh")
	ms.SetFlags(bob.Address, microstellar.FlagAuthRequired)
	ms.PayNative(bob.Address, mary.Address, "1")

	err = ms.Submit()
	if err != nil {
		log.Fatalf("Submit: %v", microstellar.ErrorString(err))
	}
}
