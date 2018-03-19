package microstellar

import (
	"log"
	"testing"
)

func TestValidAddress(t *testing.T) {
	if err := ValidAddress("GAB6FX3WVKZZRUE64H77BRWLDIOIOR4MU27L3ATNVUYKXPX5GF22TOZO"); err != nil {
		t.Error("address should be valid: ", err)
	}

	if err := ValidAddress("GA6FX3WVKZZRUE64H77BRWLDIOIOR4MU27L3ATNVUYKXPX5GF22TOZO"); err == nil {
		t.Error("address should not be valid: ", err)
	}

	if err := ValidAddress("SC56TNS6XKRCZE74OQDGHDE5KWSZQAW63CRJPH6ZQF34V3B62EQM6OP7"); err == nil {
		t.Error("seed is not a valid address: ", err)
	}
}

func TestValidSeed(t *testing.T) {
	if err := ValidSeed("GAB6FX3WVKZZRUE64H77BRWLDIOIOR4MU27L3ATNVUYKXPX5GF22TOZO"); err == nil {
		t.Error("address is not a valid seed: ", err)
	}

	if err := ValidSeed("SA6FX3WVKZZRUE64H77BRWLDIOIOR4MU27L3ATNVUYKXPX5GF22TOZO"); err == nil {
		t.Error("seed should not be valid: ", err)
	}

	if err := ValidSeed("SA6UC3LRJVNZ6DO3ZIBWUXHG6O7LKWWFTTAG2HK6QHSXZROMCVDU73RH"); err != nil {
		t.Error("seed is valid: ", err)
	}
}

func TestValidAddressOrSeed(t *testing.T) {
	if !ValidAddressOrSeed("GAB6FX3WVKZZRUE64H77BRWLDIOIOR4MU27L3ATNVUYKXPX5GF22TOZO") {
		t.Error("this is a valid address")
	}

	if ValidAddressOrSeed("SA6FX3WVKZZRUE64H77BRWLDIOIOR4MU27L3ATNVUYKXPX5GF22TOZO") {
		t.Error("this is not a valid address or a seed")
	}

	if !ValidAddressOrSeed("SA6UC3LRJVNZ6DO3ZIBWUXHG6O7LKWWFTTAG2HK6QHSXZROMCVDU73RH") {
		t.Error("this is a valid seed")
	}
}

func TestDecodeTx(t *testing.T) {
	tx := "AAAAAJb3jlBt5y04F3kXk47T9MO/Se7NcfhnIxXvWjOCzZ14AAAAZAB50HAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAuIMOnlpDFWhoO8o6VVzH4MZdIpgqr21GMRGG2riMxNoAAAAAAAAAAACYloAAAAAAAAAAAA"

	txe, err := DecodeTx(tx)
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	log.Printf("txe: %+v", txe)

	txe, err = DecodeTx(tx + "==")
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}
