package microstellar

import (
	"context"
	"fmt"
	"log"
	"time"
)

func ExampleMicroStellar_WatchPayments() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Watch for payments to address. (The fake network sends payments every 200ms.)
	watcher, err := ms.WatchPayments("GCCRUJJGPYWKQWM5NLAXUCSBCJKO37VVJ74LIZ5AQUKT6KPVCPNAGC4A",
		Opts().WithContext(context.Background()))

	if err != nil {
		log.Fatalf("Can't watch ledger: %+v", err)
	}

	// Count the number of payments received.
	paymentsReceived := 0

	go func() {
		for p := range watcher.Ch {
			paymentsReceived++
			log.Printf("WatchPayments %d: %v -- %v %v from %v to %v\n", paymentsReceived, p.Type, p.Amount, p.AssetCode, p.From, p.To)
		}

		log.Printf("## WatchPayments ## Done -- Error: %v\n", *watcher.Err)
	}()

	// Stream the ledger for about a second then stop the watcher.
	time.Sleep(1 * time.Second)
	watcher.Done()

	// Sleep a bit to wait for the done message from the goroutine.
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("%d payments received", paymentsReceived)
	// Output: 5 payments received
}

func ExampleMicroStellar_WatchTransactions() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Watch for transactions to address. (The fake network sends transactions every 200ms.)
	watcher, err := ms.WatchTransactions("GCCRUJJGPYWKQWM5NLAXUCSBCJKO37VVJ74LIZ5AQUKT6KPVCPNAGC4A",
		Opts().WithContext(context.Background()))

	if err != nil {
		log.Fatalf("Can't watch ledger: %+v", err)
	}

	// Count the number of transactions received.
	received := 0

	go func() {
		for t := range watcher.Ch {
			received++
			log.Printf("WatchTransactions %d: %v %v %v\n", received, t.ID, t.Account, t.Ledger)
		}

		log.Printf("## WatchTransactions ## Done -- Error: %v\n", *watcher.Err)
	}()

	// Stream the ledger for about a second then stop the watcher.
	time.Sleep(1 * time.Second)
	watcher.Done()

	// Sleep a bit to wait for the done message from the goroutine.
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("%d transactions received", received)
	// Output: 5 transactions received
}

func ExampleMicroStellar_WatchLedgers() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Get notified on new ledger entries in Stellar.
	watcher, err := ms.WatchLedgers(Opts().WithCursor("now"))

	if err != nil {
		log.Fatalf("Can't watch ledger: %+v", err)
	}

	// Count the number of entries seen.
	entries := 0

	go func() {
		for l := range watcher.Ch {
			entries++
			log.Printf("WatchLedgers %d: %v -- %v\n", entries, l.ID, l.TotalCoins)
		}

		log.Printf("## WatchLedgers ## Done -- Error: %v\n", *watcher.Err)
	}()

	// Stream the ledger for about a second then stop the watcher.
	time.Sleep(1 * time.Second)
	watcher.Done()

	// Sleep a bit to wait for the done message from the goroutine.
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("%d entries seen", entries)
	// Output: 5 entries seen
}
