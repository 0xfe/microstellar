package macrotest

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/0xfe/microstellar"
)

func TestMicroStellarWatchers(t *testing.T) {
	const fundSourceSeed = "SBW2N5EK5MZTKPQJZ6UYXEMCA63AO3AVUR6U5CUOIDFYCAR2X2IJIZAX"
	ms := microstellar.New("test")

	// Starting ledger watcher immediately...
	lCount := 0
	lW, err := ms.WatchLedgers(microstellar.Opts().WithCursor("now"))

	go func() {
		for l := range lW.Ch {
			lCount++
			debugf("  ## (ledger) %v: coins: %v, op_count: %v, tx_count: %v", l.ID, l.TotalCoins, l.OperationCount, l.TransactionCount)
		}

		debugf("  ## (ledger) Done -- Error: %v", *lW.Err)
	}()

	issuer := createFundedAccount(ms, fundSourceSeed, true)
	bob := createFundedAccount(ms, issuer.Seed, false)
	mary := createFundedAccount(ms, issuer.Seed, false)

	log.Print("Pair1 (issuer): ", issuer)
	log.Print("Pair2 (bob): ", bob)
	log.Print("Pair3 (mary): ", mary)

	log.Print("Starting watchers...")

	pCount := 0
	pW, err := ms.WatchPayments(issuer.Address, microstellar.Opts().WithContext(context.Background()))
	tCount := 0
	tW, err := ms.WatchTransactions(bob.Address, microstellar.Opts().WithContext(context.Background()))

	if err != nil {
		log.Fatalf("Can't watch ledger: %+v", err)
	}

	go func() {
		for p := range pW.Ch {
			pCount++
			debugf("  ## (payment:issuer) %v: %v%v%v from %v%v", p.Type, p.Amount, p.StartingBalance, p.AssetCode, p.From, p.Account)
		}

		debugf("  ## (payment:issuer) Done -- Error: %v", *pW.Err)
	}()

	go func() {
		for tx := range tW.Ch {
			tCount++
			debugf("  ## (tx:bob) %v: op_count: %v, fees: %v", tx.ID, tx.OperationCount, tx.FeePaid)
		}

		debugf("  ## (tx:bob) Done -- Error: %v", *tW.Err)
	}()

	debugf("Creating new USD asset issued by %s (issuer)...", issuer.Address)
	USD := microstellar.NewAsset("USD", issuer.Address, microstellar.Credit4Type)

	debugf("Creating USD trustline for %s (bob)...", bob.Address)
	err = ms.CreateTrustLine(bob.Seed, USD, "1000000")

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

	log.Print("Creating offer for Bob to sell 100 USD...")
	err = ms.CreateOffer(bob.Seed, USD, microstellar.NativeAsset, "2", "50")

	if err != nil {
		log.Fatalf("CreateOffer: %+v", microstellar.ErrorString(err))
	}

	debugf("Creating offer for Mary to buy Bob's assets...")
	err = ms.CreateOffer(mary.Seed, microstellar.NativeAsset, USD, "0.5", "100")

	if err != nil {
		log.Fatalf("CreateOffer: %+v", microstellar.ErrorString(err))
	}

	pW.Done()
	tW.Done()
	lW.Done()

	debugf("Sleeping for 5 seconds...")
	time.Sleep(5 * time.Second)

	debugf("Totals: tCount: %d, lCount: %d, pCount %d", tCount, lCount, pCount)

	if tCount < 2 {
		t.Errorf("didn't see enough transactions, tCount: %d", tCount)
	}
	if pCount < 2 {
		t.Errorf("didn't see enough payments, pCount: %d", tCount)
	}
	if lCount < 2 {
		t.Errorf("didn't see enough ledger entries, lCount: %d", lCount)
	}
}
