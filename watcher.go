package microstellar

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizon"
)

// Transaction represents a finalized transaction in the ledger. You can subscribe to transactions
// on the stellar network via the WatchTransactions call.
type Transaction horizon.Transaction

// TransactionWatcher is returned by WatchTransactions, which watches the ledger for transactions
// to and from an address.
type TransactionWatcher struct {
	// Ch gets a *Transaction everytime there's a new entry in the ledger.
	Ch chan *Transaction

	// TODO: another channel for structured Payment info

	// Call Done to stop watching the ledger. This closes Ch.
	Done func()

	// This is set if the stream terminates unexpectedly. Safe to check
	// after Ch is closed.
	Err *error
}

// WatchTransactions watches the ledger for transactions to and from address and streams them on a channel . Use
// Options.WithContext to set a context.Context, and Options.WithCursor to set a cursor.
func (ms *MicroStellar) WatchTransactions(address string, options ...*Options) (*TransactionWatcher, error) {
	var streamError error
	w := &TransactionWatcher{
		Ch:   make(chan *Transaction),
		Err:  &streamError,
		Done: func() {},
	}

	watcherFunc := func(params streamParams) {
		if params.tx.fake {
			w.Ch <- &Transaction{Account: "FAKE"}
			return
		}

		err := params.tx.GetClient().StreamTransactions(params.ctx, params.address, params.cursor, func(transaction horizon.Transaction) {
			debugf("WatchTransaction", "found transaction (%s) on %s", transaction.ID, transaction.Account)
			t := Transaction(transaction)
			w.Ch <- &t
		})

		if err != nil {
			debugf("WatchTransaction", "stream unexpectedly disconnected", err)
			*w.Err = errors.Wrapf(err, "stream disconnected")
			w.Done()
		}
	}

	cancelFunc, err := ms.watch("transaction", address, watcherFunc, options...)
	w.Done = cancelFunc

	return w, err
}

// Payment represents a finalized payment in the ledger. You can subscribe to payments
// on the stellar network via the WatchPayments call.
type Payment horizon.Payment

// PaymentWatcher is returned by WatchPayments, which watches the ledger for payments
// to and from an address.
type PaymentWatcher struct {
	// Ch gets a *Payment everytime there's a new entry in the ledger.
	Ch chan *Payment

	// TODO: another channel for structured Payment info

	// Call Done to stop watching the ledger. This closes Ch.
	Done func()

	// This is set if the stream terminates unexpectedly. Safe to check
	// after Ch is closed.
	Err *error
}

// WatchPayments watches the ledger for payments to and from address and streams them on a channel . Use
// Options.WithContext to set a context.Context, and Options.WithCursor to set a cursor.
func (ms *MicroStellar) WatchPayments(address string, options ...*Options) (*PaymentWatcher, error) {
	var streamError error
	w := &PaymentWatcher{
		Ch:   make(chan *Payment),
		Err:  &streamError,
		Done: func() {},
	}

	watcherFunc := func(params streamParams) {
		if params.tx.fake {
			w.Ch <- &Payment{Type: "fake"}
			return
		}

		err := params.tx.GetClient().StreamPayments(params.ctx, params.address, params.cursor, func(payment horizon.Payment) {
			debugf("WatchPayments", "found payment (%s) at %s, loading memo", payment.Type, address)
			params.tx.GetClient().LoadMemo(&payment)
			p := Payment(payment)
			w.Ch <- &p
		})

		if err != nil {
			debugf("WatchPayment", "stream unexpectedly disconnected", err)
			*w.Err = errors.Wrapf(err, "stream disconnected")
			w.Done()
		}

		close(w.Ch)
	}

	cancelFunc, err := ms.watch("payment", address, watcherFunc, options...)
	w.Done = cancelFunc

	return w, err
}

// streamParams is sent to streamFunc with the parameters for a horizon stream.
type streamParams struct {
	ctx        context.Context
	tx         *Tx
	cursor     *horizon.Cursor
	address    string
	cancelFunc func()
	err        *error
}

// streamFunc starts a horizon stream with the specified parameters.
type streamFunc func(streamParams)

// watch is a helper method to work with the Horizon Stream* methods. Returns a cancelFunc and error.
func (ms *MicroStellar) watch(entity string, address string, streamer streamFunc, options ...*Options) (func(), error) {
	logField := fmt.Sprintf("watch:%s", entity)
	debugf(logField, "watching address: %s", address)

	if err := ValidAddress(address); address != "" && err != nil {
		return nil, errors.Errorf("can't watch %s, invalid address: %s", entity, address)
	}

	tx := NewTx(ms.networkName, ms.params)

	var cursor *horizon.Cursor
	var ctx context.Context
	var cancelFunc func()

	if len(options) > 0 {
		tx.SetOptions(options[0])
		if options[0].hasCursor {
			// Ugh! Why do I have to do this?
			c := horizon.Cursor(options[0].cursor)
			cursor = &c
			debugf(logField, "starting stream for at cursor: %s", string(*cursor))
		}
		ctx = options[0].ctx
	}

	if ctx == nil {
		ctx, cancelFunc = context.WithCancel(context.Background())
	} else {
		ctx, cancelFunc = context.WithCancel(ctx)
	}

	go func() {
		if tx.fake {
		out:
			for {
				select {
				case <-ctx.Done():
					break out
				default:
					// continue
				}
				streamer(streamParams{ctx: ctx, tx: tx, cursor: cursor, address: address, cancelFunc: cancelFunc})
				time.Sleep(200 * time.Millisecond)
			}
		} else {
			streamer(streamParams{
				ctx:        ctx,
				tx:         tx,
				cursor:     cursor,
				address:    address,
				cancelFunc: cancelFunc,
			})
		}
	}()

	return cancelFunc, nil
}
