package microstellar

// AssetType represents an asset type on the stellar network.
type AssetType string

// NativeType represents native assets (like lumens.)
const NativeType AssetType = "native"

// Credit4Type represents credit assets with 4-digit codes.
const Credit4Type AssetType = "credit_alphanum4"

// Credit12Type represents credit assets with 12-digit codes.
const Credit12Type AssetType = "credit_alphanum12"

// Asset represents a specific asset class on the stellar network. For native
// assets "Code" and "Issuer" are ignored.
type Asset struct {
	Code   string
	Issuer string
	Type   AssetType
}

// NativeAsset is a convenience const representing a native asset.
var NativeAsset = &Asset{"XLM", "", NativeType}

// NewAsset creates a new asset with the given code, issuer, and assetType
func NewAsset(code string, issuer string, assetType AssetType) *Asset {
	return &Asset{code, issuer, assetType}
}

// Equals returns true if "this" and "that" represent the same asset class.
func (this Asset) Equals(that Asset) bool {
	// For native assets, don't compare code or issuer
	if this.Type == NativeType || that.Type == NativeType {
		return this.Type == that.Type
	}

	return (this.Code == that.Code && this.Issuer == that.Issuer && this.Type == that.Type)
}

// IsNative returns true if the asset is a native asset (e.g., lumens.)
func (this Asset) IsNative() bool {
	return this.Type == NativeType
}
