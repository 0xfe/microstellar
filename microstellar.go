// Package microstellar is an easy-to-use Go client for the Stellar network.
//
//   go get github.com/0xfe/microstellar
//
// Author: Mohit Muthanna Cheppudira <mohit@muthanna.com>
package microstellar

import (
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
)

// MicroStellar is the user handle to the Stellar network. Use the New function
// to create a new instance.
type MicroStellar struct {
	networkName string
	fake        bool
}

// New returns a new MicroStellar client connected to networkName ("test", "public")
func New(networkName string) *MicroStellar {
	return &MicroStellar{
		networkName: networkName,
		fake:        networkName == "fake",
	}
}

// CreateKeyPair generates a new random key pair.
func (ms *MicroStellar) CreateKeyPair() (*KeyPair, error) {
	pair, err := keypair.Random()
	if err != nil {
		return nil, err
	}

	return &KeyPair{pair.Seed(), pair.Address()}, nil
}

// FundAccount creates a new account out of address by funding it with lumens
// from sourceSeed. The minimum funding amount today is 0.5 XLM.
func (ms *MicroStellar) FundAccount(address string, sourceSeed string, amount string) error {
	payment := build.CreateAccount(
		build.Destination{AddressOrSeed: address},
		build.NativeAmount{Amount: amount})

	tx := NewTx(ms.networkName)
	tx.Build(sourceAccount(sourceSeed), payment)
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// LoadAccount loads the account information for the given address.
func (ms *MicroStellar) LoadAccount(address string) (*Account, error) {
	if ms.fake {
		return &Account{}, nil
	}

	tx := NewTx(ms.networkName)
	account, err := tx.GetClient().LoadAccount(address)

	if err != nil {
		return nil, err
	}

	return newAccountFromHorizon(account), nil
}

// Pay makes a payment of amount from source to target in the currency specified by asset.
func (ms *MicroStellar) Pay(sourceSeed string, targetAddress string, asset *Asset, amount string) error {
	var payment build.PaymentBuilder

	if asset.IsNative() {
		payment = build.Payment(
			build.Destination{AddressOrSeed: targetAddress},
			build.NativeAmount{Amount: amount})
	} else {
		payment = build.Payment(
			build.Destination{AddressOrSeed: targetAddress},
			build.CreditAmount{Code: asset.Code, Issuer: asset.Issuer, Amount: amount})
	}

	tx := NewTx(ms.networkName)
	tx.Build(sourceAccount(sourceSeed), payment)
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// PayNative makes a native asset payment of amount from source to target.
func (ms *MicroStellar) PayNative(sourceSeed string, targetAddress string, amount string) error {
	return ms.Pay(sourceSeed, targetAddress, NativeAsset, amount)
}

// CreateTrustLine creates a trustline to asset, with the specified trust limit. An empty
// limit string indicates no limit.
func (ms *MicroStellar) CreateTrustLine(sourceSeed string, asset *Asset, limit string) error {
	tx := NewTx(ms.networkName)

	if limit == "" {
		tx.Build(sourceAccount(sourceSeed), build.Trust(asset.Code, asset.Issuer))
	} else {
		tx.Build(sourceAccount(sourceSeed), build.Trust(asset.Code, asset.Issuer, limit))
	}

	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// GetBalances returns the balances in the account
// PayLumens
// Pay
// IssueAsset
// AddTrustLine
// ChangeTrustline
// AddSigners
// ChangeSigners
// Masterweight
// Op
