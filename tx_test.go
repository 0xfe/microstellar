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
func ExampleTxOptions_multisig() {
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

	// Test the EvBeforeSubmit handler.
	tx.Reset()
	handler := func(args ...interface{}) (bool, error) {
		log.Printf("Presubmit: %v", args[0])
		return true, nil
	}

	txHandler := TxHandler(handler)
	tx.SetOptions(Opts().On(EvBeforeSubmit, &txHandler))
	tx.Build(sourceAccount(keyPair.Seed))
	tx.Sign()
	tx.Submit()

	// Test SkipSignatures
	tx.Reset()
	tx.SetOptions(Opts().SkipSignatures())
	tx.Sign()
	tx.Submit()
}
