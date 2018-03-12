package microstellar

import (
	"fmt"
	"log"
)

// This example creates a passive offer on stellar's DEX.
func ExampleMicroStellar_CreateOffer() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Sell 200 USD on the DEX for lumens (at 0.5 USD/lumen). This is a passive
	// offer.
	err := ms.CreateOffer("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK",
		USD, NativeAsset, "0.5", "200",
		Opts().MakePassive())

	if err != nil {
		log.Fatalf("CreateOffer: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example updates an existing offer on the DEX.
func ExampleMicroStellar_UpdateOffer() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Update Offer ID 23456 to sell 200 USD on the DEX for lumens (at 0.4 USD/lumen).
	err := ms.UpdateOffer("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK",
		"23456", USD, NativeAsset, "0.4", "200")

	if err != nil {
		log.Fatalf("UpdateOffer: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example deletes an existing offer on the DEX.
func ExampleMicroStellar_DeleteOffer() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Delete Offer ID 23456 on the DEX.
	err := ms.DeleteOffer("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK",
		"23456", USD, NativeAsset, "0.4")

	if err != nil {
		log.Fatalf("DeleteOffer: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example creates a new offer using the ManageOffer method.
func ExampleMicroStellar_ManageOffer() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Delete Offer ID 23456 on the DEX.
	err := ms.ManageOffer("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK",
		&OfferParams{
			OfferType:  OfferCreate,
			SellAsset:  USD,
			BuyAsset:   NativeAsset,
			BuyPrice:   "0.4",
			SellAmount: "200",
		})

	if err != nil {
		log.Fatalf("ManageOffer: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}
