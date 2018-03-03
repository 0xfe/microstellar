package microstellar

// Payment is what you pass to PayWithOptions to tweak your payment transaction.
type Payment struct {
	// Use NewPayment to set these fields
	sourceSeed    string // seed for source address
	targetAddress string // address of payee
	amount        string // amount to pay

	// Use WithAsset method to set the asset type. (Defaults to NativeAsset)
	asset *Asset // asset to pay with

	// Use WithMemo* methods to set these options
	memoType MemoType // defaults to no memo
	memoText string   // additional memo text
	memoID   uint64   // additional memo ID

	// Use WithSigner to add signers to this transaction
	signerSeeds []string
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

// NewPayment creates a new payment tx from sourceSeed to targetAddress of
// amount tokens. Use the With* methods to customize the payment (e.g., change the
// asset type.)
func NewPayment(sourceSeed, targetAddress, amount string) *Payment {
	return &Payment{
		sourceSeed:    sourceSeed,
		targetAddress: targetAddress,
		amount:        amount,
		asset:         NativeAsset,
		memoType:      MemoNone,
	}
}

// WithAsset sets the Asset field on Payment p
func (p *Payment) WithAsset(asset *Asset) *Payment {
	p.asset = asset
	return p
}

// WithMemoText sets the memoType and memoText fields on Payment p
func (p *Payment) WithMemoText(text string) *Payment {
	p.memoType = MemoText
	p.memoText = text
	return p
}

// WithMemoID sets the memoType and memoID fields on Payment p
func (p *Payment) WithMemoID(id uint64) *Payment {
	p.memoType = MemoID
	p.memoID = id
	return p
}

// WithSigner adds a signer to Payment p
func (p *Payment) WithSigner(signerSeed string) *Payment {
	p.signerSeeds = append(p.signerSeeds, signerSeed)
	return p
}
