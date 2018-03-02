// Package microstellar is an easy-to-use Go client for the Stellar network.
//
// Usage:
//
//   // Create a new MicroStellar client connected to the testnet.
//   ms := microstellar.New("test")
//
//   // Generate a new random keypair.
//   pair, err := ms.CreateKeyPair()
//
//   // Display address and key
//   log.Printf("Private seed: %s, Public address: %s", pair.Seed, pair.Address)
//
//   // Fund the account with 1 lumen from an existing account.
//   err = ms.FundAccount(pair.Address, "S6 ... private key ... 3J", "1")
//
//   // Fund an account on the test network with Friendbot
//   resp, err := microstellar.FundWithFriendBot(pair.Address)
//
//   // Now load account details from ledger.
//   account, err := ms.LoadAccount(pair.Address)
//   log.Printf("Native Balance: %v XLM", account.GetNativeBalance())
//
//   // Pay someone 3 lumens
//   ms.PayNative("S-sourceSeed", "G-targetAccount", "3")
//
//   // Pay someone 1 USD issued by an anchor
//   USD := microstellar.NewAsset("USD", "ISSUERACCOUNT", microstellar.Credit4Type)
//   ms.Pay("S-sourceSeed", "G-targetAccount", USD, "3")
//
//   // Check their balance
//   account, err := ms.LoadAccount("G-targetaccount")
//   log.Printf("USD Balance: %v USD", account.GetBalance(USD))
//
// Author: Mohit Muthanna Cheppudira <mohit@muthanna.com>
package microstellar

import (
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
)

type MicroStellar struct {
	networkName string
}

// KeyPair represents a key pair for a signer on a stellar account. An account
// can have multiple signers.
type KeyPair struct {
	Seed    string // private key
	Address string // public key
}

// New returns a new MicroStellar client connected to networkName ("test", "public")
func New(networkName string) *MicroStellar {
	return &MicroStellar{
		networkName: networkName,
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
