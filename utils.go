package microstellar

import (
	"fmt"

	"github.com/stellar/go/clients/horizon"
)

// ErrorString parses the horizon error out of err.
func ErrorString(err error) string {
	var errorString string
	herr, isHorizonError := err.(*horizon.Error)

	if isHorizonError {
		resultCodes, err := herr.ResultCodes()
		if err != nil {
			errorString = fmt.Sprintf("%v", err)
		}
		errorString = fmt.Sprintf("%v", resultCodes)
	}
	return errorString
}
