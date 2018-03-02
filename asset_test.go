package microstellar

import "testing"

func TestAssetTypes(t *testing.T) {
	asset := NewAsset("QBIT", "ISSUER", Credit4Type)

	if isNative := asset.IsNative(); isNative {
		t.Errorf("wrong asset type: want %v, got %v", false, isNative)
	}

	if isNative := NativeAsset.IsNative(); !isNative {
		t.Errorf("NativeAsset is not native: want %v, got %v", true, isNative)
	}

	other1 := NewAsset("QBIT", "ISSUER", Credit4Type)

	if equals := asset.Equals(*other1); !equals {
		t.Errorf("asset.Equals: want %v, got %v", true, equals)
	}

	other2 := NewAsset("QBIT", "ISSUER2", Credit4Type)

	if equals := asset.Equals(*other2); equals {
		t.Errorf("asset.Equals: want %v, got %v", false, equals)
	}
}
