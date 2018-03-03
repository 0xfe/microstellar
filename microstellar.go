// Package microstellar is an easy-to-use Go client for the Stellar network.
//
//   go get github.com/0xfe/microstellar
//
// Author: Mohit Muthanna Cheppudira <mohit@muthanna.com>
package microstellar

import (
	"log"

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
func (ms *MicroStellar) FundAccount(sourceSeed string, address string, amount string) error {
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
		return newAccount(), nil
	}

	tx := NewTx(ms.networkName)
	account, err := tx.GetClient().LoadAccount(address)

	if err != nil {
		return nil, err
	}

	return newAccountFromHorizon(account), nil
}

// PayNative makes a native asset payment of amount from source to target.
func (ms *MicroStellar) PayNative(sourceSeed string, targetAddress string, amount string) error {
	return ms.Pay(NewPayment(sourceSeed, targetAddress, amount))
}

// Pay lets you create more advanced payment transactions (e.g., pay with credit assets, set memo, etc.)
func (ms *MicroStellar) Pay(payment *Payment) error {
	txMuts := []build.TransactionMutator{}

	paymentMuts := []interface{}{
		build.Destination{AddressOrSeed: payment.targetAddress},
	}

	if payment.asset.IsNative() {
		paymentMuts = append(paymentMuts, build.NativeAmount{Amount: payment.amount})
	} else {
		paymentMuts = append(paymentMuts,
			build.CreditAmount{Code: payment.asset.Code, Issuer: payment.asset.Issuer, Amount: payment.amount})
	}

	switch payment.memoType {
	case MemoText:
		txMuts = append(txMuts, build.MemoText{Value: payment.memoText})
	case MemoID:
		txMuts = append(txMuts, build.MemoID{Value: payment.memoID})
	}

	txMuts = append(txMuts, build.Payment(paymentMuts...))
	tx := NewTx(ms.networkName)
	tx.Build(sourceAccount(payment.sourceSeed), txMuts...)

	if len(payment.signerSeeds) > 0 {
		log.Printf("Signing multi...")
		tx.Sign(payment.signerSeeds...)
	} else {
		tx.Sign(payment.sourceSeed)
	}

	tx.Submit()
	return tx.Err()
}

// CreateTrustLine creates a trustline from sourceSeed to asset, with the specified trust limit. An empty
// limit string indicates no limit.
func (ms *MicroStellar) CreateTrustLine(sourceSeed string, asset *Asset, limit string) error {
	tx := NewTx(ms.networkName)

	if limit == "" {
		tx.Build(sourceAccount(sourceSeed), build.Trust(asset.Code, asset.Issuer))
	} else {
		tx.Build(sourceAccount(sourceSeed), build.Trust(asset.Code, asset.Issuer, build.Limit(limit)))
	}

	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// RemoveTrustLine removes an trustline from sourceSeed to an asset
func (ms *MicroStellar) RemoveTrustLine(sourceSeed string, asset *Asset) error {
	tx := NewTx(ms.networkName)
	tx.Build(sourceAccount(sourceSeed), build.RemoveTrust(asset.Code, asset.Issuer))
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// SetMasterWeight changes the master weight of sourceSeed
func (ms *MicroStellar) SetMasterWeight(sourceSeed string, weight uint32) error {
	tx := NewTx(ms.networkName)
	tx.Build(sourceAccount(sourceSeed), build.MasterWeight(weight))
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// AddSigner adds signerAddress as a signer to sourceSeed's account with weight signerWeight
func (ms *MicroStellar) AddSigner(sourceSeed string, signerAddress string, signerWeight uint32) error {
	tx := NewTx(ms.networkName)
	tx.Build(sourceAccount(sourceSeed), build.AddSigner(signerAddress, signerWeight))
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// RemoveSigner removes signerAddress as a signer from sourceSeed's account
func (ms *MicroStellar) RemoveSigner(sourceSeed string, signerAddress string) error {
	tx := NewTx(ms.networkName)
	tx.Build(sourceAccount(sourceSeed), build.RemoveSigner(signerAddress))
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// TODO:
// AddSigners
// ChangeSigners
// Set Thresholds
// Masterweight
// Op
