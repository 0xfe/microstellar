package microstellar

import (
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
)

type MicroStellar struct {
	networkName string
}

type KeyPair struct {
	Seed    string
	Address string
}

func sourceAccount(addressOrSeed string) build.SourceAccount {
	return build.SourceAccount{AddressOrSeed: addressOrSeed}
}

func New(networkName string) *MicroStellar {
	return &MicroStellar{
		networkName: networkName,
	}
}

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

func (ms *MicroStellar) LoadAccount(address string) (*Account, error) {
	tx := NewTx(ms.networkName)
	account, err := tx.GetClient().LoadAccount(address)

	if err != nil {
		return nil, err
	}

	wrappedAccount := Account(account)
	return &wrappedAccount, nil
}

// GetBalances returns the balances in the account
