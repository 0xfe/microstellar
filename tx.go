package microstellar

import (
	"fmt"
	"time"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

// Tx represents a unique stellar transaction. This is used by the MicroStellar
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
	options     *TxOptions
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

// SetOptions sets the Tx options
func (tx *Tx) SetOptions(options *TxOptions) {
	tx.options = options
}

// WithOptions sets the Tx options and returns the Tx
func (tx *Tx) WithOptions(options *TxOptions) *Tx {
	tx.SetOptions(options)
	return tx
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

	if tx.fake {
		tx.builder = &build.TransactionBuilder{}
		return nil
	}

	muts = append([]build.TransactionMutator{
		sourceAccount,
		tx.network,
		build.AutoSequence{SequenceProvider: tx.client},
	}, muts...)

	if tx.options != nil {
		switch tx.options.memoType {
		case MemoText:
			muts = append(muts, build.MemoText{Value: tx.options.memoText})
		case MemoID:
			muts = append(muts, build.MemoID{Value: tx.options.memoID})
		}
	}

	builder, err := build.Transaction(muts...)
	tx.builder = builder
	tx.err = err
	return err
}

// IsSigned returns true of the transaction is signed.
func (tx *Tx) IsSigned() bool {
	return tx.payload != ""
}

// Sign signs the transaction with every key in keys.
func (tx *Tx) Sign(keys ...string) error {
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

	if tx.fake {
		tx.payload = "FAKE"
		return nil
	}

	var txe build.TransactionEnvelopeBuilder
	var err error

	if tx.options != nil && len(tx.options.signerSeeds) > 0 {
		txe, err = tx.builder.Sign(tx.options.signerSeeds...)
	} else {
		txe, err = tx.builder.Sign(keys...)
	}

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

	if tx.fake {
		tx.response = &horizon.TransactionSuccess{Result: "fake_ok"}
		return nil
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

// MemoType sets the memotype field on the payment request.
type MemoType int

const (
	MemoNone   = MemoType(0) // No memo
	MemoID     = MemoType(1) // ID memo
	MemoText   = MemoType(2) // Text memo (max 28 chars)
	MemoHash   = MemoType(3) // Hash memo
	MemoReturn = MemoType(4) // Return hash memo
)

// TxOptions are additional parameters for a transaction. Use Opts() or NewTxOptions()
// to create a new instance.
type TxOptions struct {
	// Use With* methods to set these options
	hasFee bool
	fee    uint32

	hasTimeBounds bool
	timeBounds    time.Duration

	memoType MemoType // defaults to no memo
	memoText string   // additional memo text
	memoID   uint64   // additional memo ID

	signerSeeds []string
}

// NewTxOptions creates a new options structure for Tx.
func NewTxOptions() *TxOptions {
	return &TxOptions{
		hasFee:        false,
		hasTimeBounds: false,
		memoType:      MemoNone,
	}
}

// Opts is just an alias for NewTxOptions
func Opts() *TxOptions {
	return NewTxOptions()
}

// WithMemoText sets the memoType and memoText fields on Payment p
func (o *TxOptions) WithMemoText(text string) *TxOptions {
	o.memoType = MemoText
	o.memoText = text
	return o
}

// WithMemoID sets the memoType and memoID fields on Payment p
func (o *TxOptions) WithMemoID(id uint64) *TxOptions {
	o.memoType = MemoID
	o.memoID = id
	return o
}

// WithSigner adds a signer to Payment p
func (o *TxOptions) WithSigner(signerSeed string) *TxOptions {
	o.signerSeeds = append(o.signerSeeds, signerSeed)
	return o
}
