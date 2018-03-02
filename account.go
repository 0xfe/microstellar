package microstellar

import (
	"github.com/stellar/go/clients/horizon"
)

// KeyPair represents a key pair for a signer on a stellar account. An account
// can have multiple signers.
type KeyPair struct {
	Seed    string // private key
	Address string // public key
}

// Balance is the balance amount of the asset in the account.
type Balance struct {
	Asset  *Asset
	Amount string
}

// Account represents an account on the stellar network.
type Account struct {
	Balances      []Balance
	NativeBalance Balance
}

// newAccountFromHorizon creates a new account from a Horizon JSON response.
func newAccountFromHorizon(ha horizon.Account) *Account {
	account := &Account{}
	for _, b := range ha.Balances {
		if b.Asset.Type == string(NativeType) {
			account.NativeBalance = Balance{NativeAsset, b.Balance}
			continue
		}

		balance := Balance{
			Asset:  NewAsset(b.Asset.Code, b.Asset.Issuer, AssetType(b.Asset.Type)),
			Amount: b.Balance,
		}

		account.Balances = append(account.Balances, balance)
	}

	return account
}

// GetBalance returns the balance for asset in account. If no balance is
// found for the asset, returns "".
func (account *Account) GetBalance(asset *Asset) string {
	for _, b := range account.Balances {
		if asset.Equals(*b.Asset) {
			return b.Amount
		}
	}

	return ""
}

// GetNativeBalance returns the balance of the native currency (typically lumens)
// in the account.
func (account *Account) GetNativeBalance() string {
	return account.NativeBalance.Amount
}
