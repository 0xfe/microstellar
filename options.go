package microstellar

import (
	"context"
	"time"
)

// SortOrder is used with WithSortOrder
type SortOrder int

// For use with WithSortOrder
const (
	SortAscending  = SortOrder(0)
	SortDescending = SortOrder(1)
)

// MemoType sets the memotype field on the payment request.
type MemoType int

// Supported memo types.
const (
	MemoNone   = MemoType(0) // No memo
	MemoID     = MemoType(1) // ID memo
	MemoText   = MemoType(2) // Text memo (max 28 chars)
	MemoHash   = MemoType(3) // Hash memo
	MemoReturn = MemoType(4) // Return hash memo
)

// Options are additional parameters for a transaction. Use Opts() or NewOptions()
// to create a new instance.
type Options struct {
	// Defaults to context.Background if unset.
	ctx context.Context

	// Use With* methods to set these options
	hasFee        bool
	fee           uint32
	hasTimeBounds bool
	timeBounds    time.Duration

	// Used by all transactions.
	memoType MemoType // defaults to no memo
	memoText string   // additional memo text
	memoID   uint64   // additional memo ID

	signerSeeds []string

	// Options for query methods (Watch*, Load*)
	hasCursor      bool
	cursor         string
	hasLimit       bool
	limit          uint
	sortDescending bool

	// For offer management.
	passiveOffer bool
}

// NewOptions creates a new options structure for Tx.
func NewOptions() *Options {
	return &Options{
		ctx:            nil,
		hasFee:         false,
		hasTimeBounds:  false,
		memoType:       MemoNone,
		hasCursor:      false,
		hasLimit:       false,
		sortDescending: false,
		passiveOffer:   false,
	}
}

// Opts is just an alias for NewOptions
func Opts() *Options {
	return NewOptions()
}

// mergeOptions takes a slice of Options and merges them.
func mergeOptions(opts []*Options) *Options {
	// for now, just return the first option
	if len(opts) > 0 {
		return opts[0]
	}

	return NewOptions()
}

// WithMemoText sets the memoType and memoText fields on Payment p
func (o *Options) WithMemoText(text string) *Options {
	o.memoType = MemoText
	o.memoText = text
	return o
}

// WithMemoID sets the memoType and memoID fields on Payment p
func (o *Options) WithMemoID(id uint64) *Options {
	o.memoType = MemoID
	o.memoID = id
	return o
}

// WithSigner adds a signer to Payment p
func (o *Options) WithSigner(signerSeed string) *Options {
	o.signerSeeds = append(o.signerSeeds, signerSeed)
	return o
}

// WithContext sets the context.Context for the connection
func (o *Options) WithContext(context context.Context) *Options {
	o.ctx = context
	return o
}

// WithCursor sets the cursor for watchers and queries
func (o *Options) WithCursor(cursor string) *Options {
	o.hasCursor = true
	o.cursor = cursor
	return o
}

// WithLimit sets the limit for queries
func (o *Options) WithLimit(limit uint) *Options {
	o.hasLimit = true
	o.limit = limit
	return o
}

// WithSortOrder sets the sort order of the results
func (o *Options) WithSortOrder(order SortOrder) *Options {
	if order == SortDescending {
		o.sortDescending = true
	}
	return o
}

// MakePassive turns this into a passive offer
func (o *Options) MakePassive() *Options {
	o.passiveOffer = true
	return o
}

// TxOptions is a deprecated alias for TxOptoins
type TxOptions Options
