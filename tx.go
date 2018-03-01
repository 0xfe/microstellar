package microstellar

import (
	"fmt"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

type Tx struct {
	client      *horizon.Client
	networkName string
	network     build.Network
	builder     *build.TransactionBuilder
	payload     string
	submitted   bool
	response    *horizon.TransactionSuccess
	err         error
}

func NewTx(networkName string) *Tx {
	var network build.Network
	var client *horizon.Client

	if networkName == "test" {
		network = build.TestNetwork
		client = horizon.DefaultTestNetClient
	} else {
		network = build.PublicNetwork
		client = horizon.DefaultPublicNetClient
	}

	return &Tx{
		networkName: networkName,
		client:      client,
		network:     network,
		builder:     nil,
		payload:     "",
		submitted:   false,
		response:    nil,
		err:         nil,
	}
}

func (tx *Tx) Err() error {
	return tx.err
}

func (tx *Tx) Response() string {
	return fmt.Sprintf("%v", tx.response)
}

func (tx *Tx) Reset() {
	tx.err = nil
	tx.builder = nil
	tx.payload = ""
	tx.submitted = false
	tx.response = nil
	tx.err = nil
}

func (tx *Tx) Build(muts ...build.TransactionMutator) error {
	if tx.err != nil {
		return tx.err
	}

	if tx.builder != nil {
		tx.err = fmt.Errorf("Tx.Build: transaction already built")
		return tx.err
	}

	muts = append([]build.TransactionMutator{
		tx.network,
		build.AutoSequence{SequenceProvider: tx.client},
	}, muts...)

	builder, err := build.Transaction(muts...)
	tx.builder = builder
	return err
}

func (tx *Tx) IsSigned() bool {
	return tx.payload != ""
}

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

	resp, err := tx.client.SubmitTransaction(tx.payload)

	if err != nil {
		tx.err = err
		return err
	}

	tx.response = &resp
	tx.submitted = true
	return nil
}
