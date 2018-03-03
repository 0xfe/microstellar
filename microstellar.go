// Package microstellar is an easy-to-use Go client for the Stellar network.
//
//   go get github.com/0xfe/microstellar
//
// Author: Mohit Muthanna Cheppudira <mohit@muthanna.com>
//
// Usage notes
//
// In Stellar lingo, a private key is called a seed, and a public key is called an address. Seed
// strings start with "S", and address strings start with "G".
//
// In the methods below, "sourceSeed" is typically the private key of the account that needs
// to sign the transaction.
//
// Most methods allow you to add a *TxOptions struct at the end, which set extra parameters on the
// submitted transaction. If you add new signers via TxOptions, then sourceSeed will not be used to sign
// the transaction -- and it's okay to use a public address instead of a seed for sourceSeed.
// See examples for how to use TxOptions.
//
// You can use ErrorString(...) to extract the Horizon error from a returned error.
package microstellar

import (
	"fmt"

	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
)

// MicroStellar is the user handle to the Stellar network. Use the New function
// to create a new instance.
type MicroStellar struct {
	networkName string
	params      Params
	fake        bool
}

// Params lets you add optional parameters to New and NewTx.
type Params map[string]interface{}

// New returns a new MicroStellar client connected that operates on the network
// specified by networkName. The supported networks are:
//
//    public: the public horizon network
//    test: the public horizon testnet
//    fake: a fake network used for tests
//    custom: a custom network specified by the parameters
//
// If you're using "custom", provide the URL and Passphrase to your
// horizon network server in the parameters.
//
//    NewTx("custom", Params{
//        "url": "https://my-horizon-server.com",
//        "passphrase": "foobar"})
func New(networkName string, params ...Params) *MicroStellar {
	var p Params

	if len(params) > 0 {
		p = params[0]
	}

	return &MicroStellar{
		networkName: networkName,
		params:      p,
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
func (ms *MicroStellar) FundAccount(sourceSeed string, address string, amount string, options ...*TxOptions) error {
	if !ValidAddressOrSeed(sourceSeed) {
		return fmt.Errorf("FundAccount: invalid source address or seed: %s", sourceSeed)
	}

	if err := ValidAddress(address); err != nil {
		return fmt.Errorf("FundAccount: invalid address: %s: %v", address, err)
	}

	payment := build.CreateAccount(
		build.Destination{AddressOrSeed: address},
		build.NativeAmount{Amount: amount})

	tx := NewTx(ms.networkName, ms.params)

	if len(options) > 0 {
		tx.SetOptions(options[0])
	}

	tx.Build(sourceAccount(sourceSeed), payment)
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// LoadAccount loads the account information for the given address.
func (ms *MicroStellar) LoadAccount(address string) (*Account, error) {
	if !ValidAddressOrSeed(address) {
		return nil, fmt.Errorf("LoadAccount: invalid address: %v", address)
	}

	if ms.fake {
		return newAccount(), nil
	}

	tx := NewTx(ms.networkName, ms.params)
	account, err := tx.GetClient().LoadAccount(address)

	if err != nil {
		return nil, err
	}

	return newAccountFromHorizon(account), nil
}

// PayNative makes a native asset payment of amount from source to target.
func (ms *MicroStellar) PayNative(sourceSeed string, targetAddress string, amount string, options ...*TxOptions) error {
	return ms.Pay(sourceSeed, targetAddress, amount, NativeAsset, options...)
}

// Pay lets you make payments with credit assets.
//
//   ms.Pay("source_seed", "target_address", "3", NativeAsset, microstellar.Opts().WithMemoText("for shelter"))
func (ms *MicroStellar) Pay(sourceSeed string, targetAddress string, amount string, asset *Asset, options ...*TxOptions) error {
	if err := asset.Validate(); err != nil {
		return fmt.Errorf("Pay: invalid asset: %v", err)
	}

	if !ValidAddressOrSeed(sourceSeed) {
		return fmt.Errorf("Pay: invalid source address or seed: %s", sourceSeed)
	}

	if err := ValidAddress(targetAddress); err != nil {
		return fmt.Errorf("Pay: invalid address: %s: %v", targetAddress, err)
	}

	paymentMuts := []interface{}{
		build.Destination{AddressOrSeed: targetAddress},
	}

	if asset.IsNative() {
		paymentMuts = append(paymentMuts, build.NativeAmount{Amount: amount})
	} else {
		paymentMuts = append(paymentMuts,
			build.CreditAmount{Code: asset.Code, Issuer: asset.Issuer, Amount: amount})
	}

	tx := NewTx(ms.networkName, ms.params)

	if len(options) > 0 {
		tx.SetOptions(options[0])
	}

	tx.Build(sourceAccount(sourceSeed), build.Payment(paymentMuts...))
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// CreateTrustLine creates a trustline from sourceSeed to asset, with the specified trust limit. An empty
// limit string indicates no limit.
func (ms *MicroStellar) CreateTrustLine(sourceSeed string, asset *Asset, limit string, options ...*TxOptions) error {
	if !ValidAddressOrSeed(sourceSeed) {
		return fmt.Errorf("CreateTrustLine: invalid source address or seed: %s", sourceSeed)
	}

	if err := asset.Validate(); err != nil {
		return fmt.Errorf("CreateTrustLine: invalid asset: %v", err)
	}

	tx := NewTx(ms.networkName, ms.params)

	if len(options) > 0 {
		tx.SetOptions(options[0])
	}

	if limit == "" {
		tx.Build(sourceAccount(sourceSeed), build.Trust(asset.Code, asset.Issuer))
	} else {
		tx.Build(sourceAccount(sourceSeed), build.Trust(asset.Code, asset.Issuer, build.Limit(limit)))
	}

	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// RemoveTrustLine removes an trustline from sourceSeed to an asset.
func (ms *MicroStellar) RemoveTrustLine(sourceSeed string, asset *Asset, options ...*TxOptions) error {
	if !ValidAddressOrSeed(sourceSeed) {
		return fmt.Errorf("RemoveTrustLine: invalid source address or seed: %s", sourceSeed)
	}

	if err := asset.Validate(); err != nil {
		return fmt.Errorf("RemoveTrustLine: invalid asset: %v", err)
	}

	tx := NewTx(ms.networkName, ms.params)

	if len(options) > 0 {
		tx.SetOptions(options[0])
	}

	tx.Build(sourceAccount(sourceSeed), build.RemoveTrust(asset.Code, asset.Issuer))
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// SetMasterWeight changes the master weight of sourceSeed.
func (ms *MicroStellar) SetMasterWeight(sourceSeed string, weight uint32, options ...*TxOptions) error {
	if !ValidAddressOrSeed(sourceSeed) {
		return fmt.Errorf("SetMasterWeight: invalid source address or seed: %s", sourceSeed)
	}

	tx := NewTx(ms.networkName, ms.params)

	if len(options) > 0 {
		tx.SetOptions(options[0])
	}

	tx.Build(sourceAccount(sourceSeed), build.MasterWeight(weight))
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// AddSigner adds signerAddress as a signer to sourceSeed's account with weight signerWeight.
func (ms *MicroStellar) AddSigner(sourceSeed string, signerAddress string, signerWeight uint32, options ...*TxOptions) error {
	if !ValidAddressOrSeed(sourceSeed) {
		return fmt.Errorf("AddSigner: invalid source address or seed: %s", sourceSeed)
	}

	if !ValidAddressOrSeed(signerAddress) {
		return fmt.Errorf("AddSigner: invalid signer address or seed: %s", signerAddress)
	}

	tx := NewTx(ms.networkName, ms.params)

	if len(options) > 0 {
		tx.SetOptions(options[0])
	}

	tx.Build(sourceAccount(sourceSeed), build.AddSigner(signerAddress, signerWeight))
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// RemoveSigner removes signerAddress as a signer from sourceSeed's account.
func (ms *MicroStellar) RemoveSigner(sourceSeed string, signerAddress string, options ...*TxOptions) error {
	if !ValidAddressOrSeed(sourceSeed) {
		return fmt.Errorf("RemoveSigner: invalid source address or seed: %s", sourceSeed)
	}

	if !ValidAddressOrSeed(signerAddress) {
		return fmt.Errorf("AddSigner: invalid signer address or seed: %s", signerAddress)
	}

	tx := NewTx(ms.networkName, ms.params)

	if len(options) > 0 {
		tx.SetOptions(options[0])
	}

	tx.Build(sourceAccount(sourceSeed), build.RemoveSigner(signerAddress))
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}

// SetThresholds sets the signing thresholds for the account.
func (ms *MicroStellar) SetThresholds(sourceSeed string, low, medium, high uint32, options ...*TxOptions) error {
	if !ValidAddressOrSeed(sourceSeed) {
		return fmt.Errorf("SetThresholds: invalid source address or seed: %s", sourceSeed)
	}

	tx := NewTx(ms.networkName, ms.params)

	if len(options) > 0 {
		tx.SetOptions(options[0])
	}

	tx.Build(sourceAccount(sourceSeed), build.SetThresholds(low, medium, high))
	tx.Sign(sourceSeed)
	tx.Submit()
	return tx.Err()
}
