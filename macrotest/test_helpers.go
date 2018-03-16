package macrotest

import (
	"log"
	"strconv"

	"github.com/0xfe/microstellar"
	"github.com/sirupsen/logrus"
)

func debugf(msg string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"test": "macrotest"}).Infof(msg, args...)
}

// Send unused friendbot funds here.
const homeAddress string = "GB7RTQME2RAOPRBDFBICCP3UDLCIJOSP7ZWCW5IL7Z6L4FNVLZMEWX2G"

// Helper function to remove stupid "if err != nil" checks.
func failOnError(i interface{}, err error) interface{} {
	if err != nil {
		log.Fatal(err)
	}

	return i
}

// Helper function for methods with only error return values.
func failOnError0(method string, err error) {
	if err != nil {
		log.Fatalf("%s: %+v", method, err)
	}
}

// Helper function to create a new funded stellar account on the testnet.
func createFundedAccount(ms *microstellar.MicroStellar, fundSourceSeed string, useFriendBot bool) microstellar.KeyPair {
	// Create a key pair
	keyPair := failOnError(ms.CreateKeyPair()).(*microstellar.KeyPair)
	debugf("Created key pair: %v", keyPair)

	if useFriendBot {
		// Try to fund it with friendbot
		debugf("Funding new key with friendbot...")
		resp, _ := microstellar.FundWithFriendBot(keyPair.Address)
		debugf("Friendbot says: %v", resp)
	}

	// Load the account to see if there are funds
	debugf("Checking balance on new key...")
	account, err := ms.LoadAccount(keyPair.Address)
	var floatBalance float64 = 0

	if err == nil {
		balance := account.GetNativeBalance()
		debugf("Got balance: %v", balance)
		floatBalance = failOnError(strconv.ParseFloat(balance, 64)).(float64)
	}

	if floatBalance == 0 {
		amount := "2000"
		if !useFriendBot {
			amount = "300"
		}

		debugf("Looks like it's empty. Funding via source account...")
		err := ms.FundAccount(fundSourceSeed, keyPair.Address, amount, microstellar.Opts().WithMemoText("initial fund"))
		if err != nil {
			log.Fatalf("Funding failed: %v", microstellar.ErrorString(err))
		}
		debugf("Payment sent: %s lumens", amount)
	} else {
		debugf("Yay! Friendbot paid us. Sending some lumens back to fundSource...")
		err := ms.PayNative(keyPair.Seed, homeAddress, "5000", microstellar.Opts().WithMemoText("friendbot payback"))

		if err != nil {
			log.Fatalf(microstellar.ErrorString(err))
		}
	}

	return *keyPair
}

// Helper function to show the asset balance of a specific account
func showBalance(ms *microstellar.MicroStellar, asset *microstellar.Asset, name, address string) {
	debugf("Balances for %s: %s", name, address)
	account, err := ms.LoadAccount(address)

	if err != nil {
		log.Fatalf("Canl't load balances for %s: %s", name, address)
	}

	log.Print("  Master weight: ", account.GetMasterWeight())
	log.Print("  XLM: ", account.GetNativeBalance())
	log.Printf("  %s: %v", asset.Code, account.GetBalance(asset))

	for i, s := range account.Signers {
		debugf("  signer %d (type: %v, weight: %v): %v", i, s.Type, s.Weight, s.PublicKey)
	}
}
