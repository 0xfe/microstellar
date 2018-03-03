package microstellar

import (
	"fmt"
	"log"
	"testing"
)

// Payments with memotext and memoid
func ExampleTxOptions_memotext() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Pay 1 USD to targetAddress and set the memotext field
	err := ms.Pay("sourceSeed", "targetAddress", "1", USD, Opts().WithMemoText("boo"))

	if err != nil {
		log.Fatalf("Pay (memotext): %v", ErrorString(err))
	}

	// Pay 1 USD to targetAddress and set the memotext field
	err = ms.Pay("sourceSeed", "targetAddress", "1", USD, Opts().WithMemoID(42))

	if err != nil {
		log.Fatalf("Pay (memoid): %v", ErrorString(err))
	}

	fmt.Printf("ok")
	// Output: ok
}

// Makes a multisignature payment
func ExampleTxOptions_multisig() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Pay 1 USD to targetAddress and set the memotext field
	err := ms.Pay("sourceSeed", "targetAddress", "1", USD,
		Opts().WithMemoText("multisig").
			WithSigner("SAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ").
			WithSigner("SBIUIQNMSXTGR4TGZETSQCGBTIF32G2L5D4AML4LFTMTHKM44UABFDMS"))

	if err != nil {
		log.Fatalf("Pay (memotext): %v", ErrorString(err))
	}

	// Pay 1 USD to targetAddress and set the memotext field
	err = ms.Pay("sourceSeed", "targetAddress", "1", USD, Opts().WithMemoID(73223))

	if err != nil {
		log.Fatalf("Pay (memoid): %v", ErrorString(err))
	}

	fmt.Printf("ok")
	// Output: ok
}

func TestTx(t *testing.T) {
	ms := New("fake")
	tx := NewTx("fake")

	keyPair, _ := ms.CreateKeyPair()

	err := tx.Sign(keyPair.Seed)

	if err == nil {
		t.Errorf("signing should not succeed: want %v, got nil", err)
	}

	err = tx.Build(sourceAccount(keyPair.Seed))
	if err == nil {
		t.Errorf("build failed: want nil, got %v", err)
	}

	err = tx.Build(sourceAccount(keyPair.Seed))
	if err == nil {
		t.Errorf("duplicate build should fail: want %v, got nil", err)
	}

	tx.Reset()
	err = tx.Build(sourceAccount(keyPair.Seed))
	if err != nil {
		t.Errorf("build failed: want nil, got %v", err)
	}

	err = tx.Sign(keyPair.Seed)
	if err != nil {
		t.Errorf("sign failed: want nil, got %v", err)
	}

	err = tx.Submit()
	if err != nil {
		t.Errorf("submit failed: want nil, got %v", err)
	}

	if tx.Err() != nil {
		t.Errorf("tx.Err() should be nil: got %v", tx.Err())

	}
}
