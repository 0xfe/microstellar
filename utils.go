package microstellar

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/strkey"
)

type AddressOrSeed string

// ValidAddress returns error if address is an invalid stellar address
func ValidAddress(address string) error {
	_, err := strkey.Decode(strkey.VersionByteAccountID, address)
	return err
}

// ValidSeed returns error if the seed is invalid
func ValidSeed(seed string) error {
	_, err := strkey.Decode(strkey.VersionByteSeed, seed)
	return err
}

// ValidAddressOrSeed returns true if the string is a valid address or seed
func ValidAddressOrSeed(addressOrSeed string) bool {
	err := ValidAddress(addressOrSeed)
	if err == nil {
		return true
	}

	err = ValidSeed(addressOrSeed)
	return err == nil
}

// HorizonErrorString parses the horizon error out of err.
func ErrorString(err error) string {
	var errorString string
	herr, isHorizonError := err.(*horizon.Error)

	if isHorizonError {
		resultCodes, err := herr.ResultCodes()
		if err != nil {
			errorString = fmt.Sprintf("%v", err)
		}
		errorString = fmt.Sprintf("Codes: %v, Problem: %v", resultCodes, herr.Problem)
	} else {
		errorString = fmt.Sprintf("%v", err)
	}

	return errorString
}

// FundWithFriendBot funds address on the test network with some initial funds.
func FundWithFriendBot(address string) (string, error) {
	resp, err := http.Get("https://horizon-testnet.stellar.org/friendbot?addr=" + address)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
