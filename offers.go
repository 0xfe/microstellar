package microstellar

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/stellar/go/build"
)

// OfferType tells ManagedOffer what operation to perform
type OfferType int

// The available offer types.
const (
	OfferCreate        = OfferType(0)
	OfferCreatePassive = OfferType(1)
	OfferUpdate        = OfferType(2)
	OfferDelete        = OfferType(3)
)

// OfferParams specify the parameters
type OfferParams struct {
	// Create, update, or delete.
	OfferType OfferType

	// The asset that's being sold on the DEX.
	SellAsset *Asset

	// The asset that you want to buy on the DEX.
	BuyAsset *Asset

	// How much you're willing to pay (in SellAsset units) per unit of BuyAsset.
	BuyPrice string

	// How many units of SellAsset are you selling?
	SellAmount string

	// Existing offer ID (for Update and Delete)
	OfferID string
}

// ManageOffer lets you trade on the DEX. See the Create/Update/DeleteOffer methods below
// to see how this is used.
func (ms *MicroStellar) ManageOffer(sourceSeed string, params *OfferParams, options ...*TxOptions) error {
	if !ValidAddressOrSeed(sourceSeed) {
		return errors.Errorf("invalid source address or seed: %s", sourceSeed)
	}

	if err := params.BuyAsset.Validate(); err != nil {
		return errors.Wrap(err, "ManageOffer")
	}

	if err := params.SellAsset.Validate(); err != nil {
		return errors.Wrap(err, "ManageOffer")
	}

	rate := build.Rate{
		Selling: genBuildAsset(params.SellAsset),
		Buying:  genBuildAsset(params.BuyAsset),
		Price:   build.Price(params.BuyPrice),
	}

	amount := build.Amount(params.SellAmount)

	var offerID uint64
	if params.OfferID != "" {
		var err error
		if offerID, err = strconv.ParseUint(params.OfferID, 10, 64); err != nil {
			return errors.Wrapf(err, "ManageOffer: bad OfferID: %v", params.OfferID)
		}
	}

	var builder build.ManageOfferBuilder
	switch params.OfferType {
	case OfferCreate:
		builder = build.CreateOffer(rate, amount)
	case OfferCreatePassive:
		builder = build.CreatePassiveOffer(rate, amount)
	case OfferUpdate:
		builder = build.UpdateOffer(rate, amount, build.OfferID(offerID))
	case OfferDelete:
		builder = build.DeleteOffer(rate, build.OfferID(offerID))
	default:
		return errors.Errorf("ManageOffer: bad OfferType: %v", params.OfferType)
	}

	tx := NewTx(ms.networkName, ms.params)

	if len(options) > 0 {
		tx.SetOptions(options[0])
	}

	tx.Build(sourceAccount(sourceSeed), builder)
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// CreateOffer creates an offer to trade sellAmount of sellAsset held by sourceSeed for buyAsset at
// buyPrice (per unit of buyAsset.) The offer is made on Stellar's decentralized exchange (DEX.)
//
// You can use add Opts().MakePassive() to make this a passive offer.
func (ms *MicroStellar) CreateOffer(sourceSeed string, sellAsset *Asset, buyAsset *Asset, sellAmount string, buyPrice string, options ...*TxOptions) error {
	offerType := OfferCreate

	if len(options) > 0 {
		if options[0].passiveOffer {
			offerType = OfferCreatePassive
		}
	}

	return ms.ManageOffer(sourceSeed, &OfferParams{
		OfferType:  offerType,
		SellAsset:  sellAsset,
		SellAmount: sellAmount,
		BuyAsset:   buyAsset,
		BuyPrice:   buyPrice,
	}, options...)
}

// UpdateOffer updates the existing offer with ID offerID on the DEX.
func (ms *MicroStellar) UpdateOffer(sourceSeed string, offerID string, sellAsset *Asset, buyAsset *Asset, sellAmount string, buyPrice string, options ...*TxOptions) error {
	return ms.ManageOffer(sourceSeed, &OfferParams{
		OfferType:  OfferUpdate,
		SellAsset:  sellAsset,
		SellAmount: sellAmount,
		BuyAsset:   buyAsset,
		BuyPrice:   buyPrice,
		OfferID:    offerID,
	}, options...)
}

// DeleteOffer deletes the specified parameters (assets, price, ID) on the DEX.
func (ms *MicroStellar) DeleteOffer(sourceSeed string, offerID string, sellAsset *Asset, buyAsset *Asset, buyPrice string, options ...*TxOptions) error {
	return ms.ManageOffer(sourceSeed, &OfferParams{
		OfferType: OfferUpdate,
		SellAsset: sellAsset,
		BuyAsset:  buyAsset,
		BuyPrice:  buyPrice,
		OfferID:   offerID,
	}, options...)
}
