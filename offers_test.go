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

	// Sell 200 USD on the DEX for lumens (at 2 lumens per USD). This is a passive
	// offer. (This is equivalent to an offer to buy 400 lumens for 200 USD.)
	err := ms.CreateOffer("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK",
		USD, NativeAsset, "2", "200",
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

	// Update Offer ID 23456 to sell 200 USD on the DEX for lumens (at 1 lumen / USD.)
	err := ms.UpdateOffer("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK",
		"23456", USD, NativeAsset, "1", "200")

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

	// Create an offer to buy 200 lumens at 2 lumens/dollar.
	err := ms.ManageOffer("SCSMBQYTXKZYY7CLVT6NPPYWVDQYDOQ6BB3QND4OIXC7762JYJYZ3RMK",
		&OfferParams{
			OfferType:  OfferCreate,
			SellAsset:  USD,
			BuyAsset:   NativeAsset,
			Price:      "2",
			SellAmount: "100",
		})

	if err != nil {
		log.Fatalf("ManageOffer: %v", err)
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example lists the offers currently out by an address.
func ExampleMicroStellar_LoadOffers() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Get at most 10 offers made by address in descending order
	offers, err := ms.LoadOffers("GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ",
		Opts().WithLimit(10).WithSortOrder(SortDescending))

	if err != nil {
		log.Fatalf("LoadOffers: %v", err)
	}

	for _, o := range offers {
		log.Printf("Offer ID: %v, Selling: %v, Price: %v, Amount: %v", o.ID, o.Selling.Code, o.Price, o.Amount)
	}

	fmt.Printf("ok")
	// Output: ok
}

// This example lists all asks on the DEX between USD <-> XLM
func ExampleMicroStellar_LoadOrderBook() {
	// Create a new MicroStellar client connected to a fake network. To
	// use a real network replace "fake" below with "test" or "public".
	ms := New("fake")

	// Custom USD asset issued by specified issuer
	USD := NewAsset("USD", "GAIUIQNMSXTTR4TGZETSQCGBTIF32G2L5P4AML4LFTMTHKM44UHIN6XQ", Credit4Type)

	// Get at most 10 orders made between USD and XLM
	orderbook, err := ms.LoadOrderBook(USD, NativeAsset,
		Opts().WithLimit(10).WithSortOrder(SortDescending))

	if err != nil {
		log.Fatalf("LoadOrderBook: %v", err)
	}

	// List all the returned asks.
	for _, ask := range orderbook.Asks {
		log.Printf("ask: %s %s for %s %s/%s", ask.Amount, orderbook.Base.Code, ask.Price, orderbook.Counter.Code, orderbook.Base.Code)
	}

	fmt.Printf("ok")
	// Output: ok
}
