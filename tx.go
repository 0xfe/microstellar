package microstellar

import (
	"fmt"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

// Tx represents a unique stellar operation. This is used by the MicroStellar
// library to abstract away the Horizon API and transport. To reuse Tx
// instances, you must call Tx.Reset() between operations.
//
// This struct is not thread-safe by design -- you must use separate instances
// in different goroutines.
//
// Unless you're hacking around in the guts, you should not need to use Tx.
type Tx struct {
	client      *horizon.Client
	networkName string
	network     build.Network
	fake        bool
	builder     *build.TransactionBuilder
	payload     string
	submitted   bool
	response    *horizon.TransactionSuccess
	err         error
}

// NewTx returns a new Tx that operates on networkName ("test", "public".)
func NewTx(networkName string) *Tx {
	var network build.Network
	var client *horizon.Client

	fake := false

	if networkName == "test" {
		network = build.TestNetwork
		client = horizon.DefaultTestNetClient
	} else if networkName == "fake" {
		network = build.TestNetwork
		client = horizon.DefaultTestNetClient
		fake = true
	} else {
		network = build.PublicNetwork
		client = horizon.DefaultPublicNetClient
	}

	return &Tx{
		networkName: networkName,
		client:      client,
		network:     network,
		fake:        fake,
		builder:     nil,
		payload:     "",
		submitted:   false,
		response:    nil,
		err:         nil,
	}
}

// GetClient returns the underlying horizon client handle.
func (tx *Tx) GetClient() *horizon.Client {
	return tx.client
}

// Err returns the error from the most recent failed operation.
func (tx *Tx) Err() error {
	return tx.err
}

// Response returns the horison response for the submitted operation.
func (tx *Tx) Response() string {
	return fmt.Sprintf("%v", tx.response)
}

// Reset clears all internal state, so you can run a new operation.
func (tx *Tx) Reset() {
	tx.err = nil
	tx.builder = nil
	tx.payload = ""
	tx.submitted = false
	tx.response = nil
	tx.err = nil
}

func sourceAccount(addressOrSeed string) build.SourceAccount {
	return build.SourceAccount{AddressOrSeed: addressOrSeed}
}

// Build creates a new operation out of the provided mutators.
func (tx *Tx) Build(sourceAccount build.TransactionMutator, muts ...build.TransactionMutator) error {
	if tx.err != nil {
		return tx.err
	}

	if tx.builder != nil {
		tx.err = fmt.Errorf("Tx.Build: transaction already built")
		return tx.err
	}

	muts = append([]build.TransactionMutator{
		sourceAccount,
		tx.network,
		build.AutoSequence{SequenceProvider: tx.client},
	}, muts...)

	builder, err := build.Transaction(muts...)
	tx.builder = builder
	tx.err = err
	return err
}

// IsSigned returns true of the transaction is signed.
func (tx *Tx) IsSigned() bool {
	return tx.payload != ""
}

// Sign signs the transaction with key.
func (tx *Tx) Sign(key string) error {
	if tx.err != nil {
		return tx.err
	}

	if tx.builder == nil {
		tx.err = fmt.Errorf("Tx.Sign: transaction not built")
		return tx.err
	}

	if tx.IsSigned() {
		tx.err = fmt.Errorf("Tx.Sign: transaction already signed")
		return tx.err
	}

	txe, err := tx.builder.Sign(key)

	if err != nil {
		tx.err = err
		return err
	}

	tx.payload, err = txe.Base64()

	if err != nil {
		tx.err = err
		return err
	}

	return nil
}

// Submit sends the transaction to the stellar network.
func (tx *Tx) Submit() error {
	if tx.fake {
		tx.response = &horizon.TransactionSuccess{Result: "fake_ok"}
		return nil
	}

	if tx.err != nil {
		return tx.err
	}

	if !tx.IsSigned() {
		tx.err = fmt.Errorf("Tx.Submit: transaction not signed")
		return tx.err
	}

	if tx.submitted {
		tx.err = fmt.Errorf("Tx.Submit: transaction already submitted")
		return tx.err
	}

	resp, err := tx.client.SubmitTransaction(tx.payload)

	if err != nil {
		tx.err = err
		return err
	}

	tx.response = &resp
	tx.submitted = true
	return nil
}
