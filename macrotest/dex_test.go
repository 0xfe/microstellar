package macrotest

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/0xfe/microstellar"
)

func showOffersAndBalances(ms *microstellar.MicroStellar, asset *microstellar.Asset, name, address string) {
	showBalance(ms, asset, name, address)

	offers, err := ms.LoadOffers(address)

	if err != nil {
		log.Fatalf("LoadOffers: %+v", microstellar.ErrorString(err))
	}

	debugf("Offers by %s: %d avaialable", name, len(offers))
	for i, o := range offers {
		offerJSON, err := json.MarshalIndent(o, "", "  ")
		if err != nil {
			log.Fatalf("MarshalIndent: %v", err)
		}
		debugf("Offer %d:\n%s", i, string(offerJSON))
	}
}

// TestMicroStellarOffers implement end-to-end DEX tests
func TestMicroStellarOffers(t *testing.T) {
	const fundSourceSeed = "SBW2N5EK5MZTKPQJZ6UYXEMCA63AO3AVUR6U5CUOIDFYCAR2X2IJIZAX"

	ms := microstellar.New("test")

	// Create a key pair
	issuer := createFundedAccount(ms, fundSourceSeed, true)
	bob := createFundedAccount(ms, issuer.Seed, false)
	mary := createFundedAccount(ms, issuer.Seed, false)

	log.Print("Pair1 (issuer): ", issuer)
	log.Print("Pair2 (bob): ", bob)
	log.Print("Pair3 (mary): ", mary)

	debugf("Creating new USD asset issued by %s (issuer)...", issuer.Address)
	USD := microstellar.NewAsset("USD", issuer.Address, microstellar.Credit4Type)

	debugf("Creating USD trustline for %s (bob)...", bob.Address)
	err := ms.CreateTrustLine(bob.Seed, USD, "1000000")

	if err != nil {
		log.Fatalf("CreateTrustLine: %+v", err)
	}

	debugf("Creating USD trustline for %s (mary)...", mary.Address)
	err = ms.CreateTrustLine(mary.Seed, USD, "1000000")

	if err != nil {
		log.Fatalf("CreateTrustLine: %+v", err)
	}

	log.Print("Issuing USD to bob...")
	err = ms.Pay(issuer.Seed, bob.Address, "500000", USD)

	if err != nil {
		log.Fatalf("Pay: %+v", microstellar.ErrorString(err))
	}

	log.Print("Issuing USD to mary...")
	err = ms.Pay(issuer.Seed, mary.Address, "500000", USD)

	if err != nil {
		log.Fatalf("Pay: %+v", microstellar.ErrorString(err))
	}

	// Check balances
	showBalance(ms, USD, "bob", bob.Address)
	showBalance(ms, USD, "mary", mary.Address)

	log.Print("Creating offer for Bob to sell 100 USD...")
	err = ms.CreateOffer(bob.Seed, USD, microstellar.NativeAsset, "2", "50")

	if err != nil {
		log.Fatalf("CreateOffer: %+v", microstellar.ErrorString(err))
	}

	// Check balances
	showOffersAndBalances(ms, USD, "bob", bob.Address)

	debugf("Creating offer for Mary to buy Bob's assets...")
	err = ms.CreateOffer(mary.Seed, microstellar.NativeAsset, USD, "0.5", "100")

	if err != nil {
		log.Fatalf("CreateOffer: %+v", microstellar.ErrorString(err))
	}

	debugf("Expecting 0 offers for Mary and Bob...")
	showOffersAndBalances(ms, USD, "bob", bob.Address)
	showOffersAndBalances(ms, USD, "mary", mary.Address)

	debugf("Creating another offer from Mary to sell XLM...")
	err = ms.CreateOffer(mary.Seed, microstellar.NativeAsset, USD, "0.5", "20")

	if err != nil {
		log.Fatalf("CreateOffer: %+v", microstellar.ErrorString(err))
	}

	// Load her offers to get qty and offer ID
	offers, err := ms.LoadOffers(mary.Address)

	if err != nil {
		log.Fatalf("LoadOffers: %+v", microstellar.ErrorString(err))
	}

	if len(offers) != 1 {
		log.Fatalf("wrong number of offers, want %v, got %v", 1, len(offers))
	}

	// Update her offer
	debugf("Updating mary's offer...")
	offerID := fmt.Sprintf("%d", offers[0].ID)
	err = ms.UpdateOffer(mary.Seed, offerID, microstellar.NativeAsset, USD, "0.5", "5")

	if err != nil {
		log.Fatalf("UpdateOffer: %+v", microstellar.ErrorString(err))
	}

	// Load her offers to get amount
	offers, err = ms.LoadOffers(mary.Address)

	if err != nil {
		log.Fatalf("LoadOffers: %+v", microstellar.ErrorString(err))
	}

	if offers[0].Amount != "5.0000000" {
		showOffersAndBalances(ms, USD, "mary", mary.Address)
		log.Fatalf("wrong amount, want %v, got %v", "5.0000000", offers[0].Amount)
	}

	// Delete the offer
	debugf("Deleting mary's offer...")
	err = ms.DeleteOffer(mary.Seed, offerID, microstellar.NativeAsset, USD, "0.5")

	if err != nil {
		log.Fatalf("DeleteOffer: %+v", microstellar.ErrorString(err))
	}

	// Load her offers to get amount
	offers, err = ms.LoadOffers(mary.Address)

	if err != nil {
		log.Fatalf("LoadOffers: %+v", microstellar.ErrorString(err))
	}

	if len(offers) != 0 {
		showOffersAndBalances(ms, USD, "mary", mary.Address)
		log.Fatalf("wrong number of offers, want %v, got %v", 0, len(offers))
	}
}

func TestMicroStellarPathPayments(t *testing.T) {
	const fundSourceSeed = "SBW2N5EK5MZTKPQJZ6UYXEMCA63AO3AVUR6U5CUOIDFYCAR2X2IJIZAX"

	ms := microstellar.New("test")

	issuer := createFundedAccount(ms, fundSourceSeed, true)
	bob := createFundedAccount(ms, issuer.Seed, false)
	mary := createFundedAccount(ms, issuer.Seed, false)
	usdHolder := createFundedAccount(ms, issuer.Seed, false)
	eurHolder := createFundedAccount(ms, issuer.Seed, false)
	inrHolder := createFundedAccount(ms, issuer.Seed, false)

	/*
		Uncomment this section, and comment out the trustline loop below (and the funding above) to speed up debugging.

			issuer := microstellar.KeyPair{"SAWFL2IHE3WVXYQ7DNU2ERZFMJ5ESN7G7Z4FKW5EATGDOB3M7SLVX7CG", "GCFEQ72ADTAK4NH5VQ2STBSSIBN5GSHSYPLIORH6ILW3LZEIJ7XJKVDE"}
			bob := microstellar.KeyPair{"SCXTE2J7YBRM7UY5RRD7AH3DLMCSOAK4EM7VXK7T7ZUBUYZJL5P2XXJM", "GCBPBDZNVOKXW5DZIMNCUQ4UQAQLZNPHIAVW2FHUDXY72DNSMVJFXUBK"}
			mary := microstellar.KeyPair{"SDSVM3SBYNIHGCGKJTC22OSXKYSB4ESFNUTTTVRYSIIAOSRKZSI2P2RD", "GDSHCBTJ6ARCSG3BCLBGLO7LD67LM4OCQ4V3VFUWTRVCJKWUY2SDGHVM"}
			usdHolder := microstellar.KeyPair{"SAIGAGTM2L2QWR437QPMM4BBAEFXTLTDGRYBZRD7QJSHRZDO2DPG266X", "GA57OMXZVWLBWBMM7Q3GHQ2EA7UUSFWY4O7RLPAFCAC3GHHPKOZNPWMW"}
			eurHolder := microstellar.KeyPair{"SANYOWPAKMHQYMCK4R4OUJBCVJM4SFG7AIGEIKZ3QFWYUQKINB3GRSFF", "GAOUDB7VAN52IMNA7PJ3FOQQXF3YN3R3TZJPAUIOU6ZZA4VAPWXIGHZC"}
			inrHolder := microstellar.KeyPair{"SBYC4QC57C46QR7OHPCFCEEBAZSKVOZDDAAC5ZMW6KUMYJR3KJTJJPSN", "GAYQ25MNFVJJP4MJT2VM2TDA5FJW6NFVIALAVOOT4SSVIIZY6XZXLX63"}
	*/

	log.Print("Pair1 (issuer): ", issuer)
	log.Print("Pair2 (bob): ", bob)
	log.Print("Pair3 (mary): ", mary)
	log.Print("Pair4 (usdHolder): ", usdHolder)
	log.Print("Pair5 (eurHolder): ", eurHolder)
	log.Print("Pair6 (inrHolder): ", inrHolder)

	debugf("Creating new USD, EUR, and INR assets issued by %s (issuer)...", issuer.Address)
	XLM := microstellar.NativeAsset
	USD := microstellar.NewAsset("USD", issuer.Address, microstellar.Credit4Type)
	EUR := microstellar.NewAsset("EUR", issuer.Address, microstellar.Credit4Type)
	INR := microstellar.NewAsset("INR", issuer.Address, microstellar.Credit4Type)

	debugf("Creating USD, EUR, and INR trustlines for everyone...")
	for i, k := range []microstellar.KeyPair{bob, mary, usdHolder, eurHolder, inrHolder} {
		debugf("Keypair %d USD", i+1)
		failOnError0("CreateTrustLine", ms.CreateTrustLine(k.Seed, USD, "1000000"))
		failOnError0("Pay (issue)", ms.Pay(issuer.Seed, k.Address, "500000", USD))

		debugf("Keypair %d EUR", i+1)
		failOnError0("CreateTrustLine", ms.CreateTrustLine(k.Seed, EUR, "1000000"))
		failOnError0("Pay (issue)", ms.Pay(issuer.Seed, k.Address, "500000", EUR))

		debugf("Keypair %d INR", i+1)
		failOnError0("CreateTrustLine", ms.CreateTrustLine(k.Seed, INR, "1000000"))
		failOnError0("Pay (issue)", ms.Pay(issuer.Seed, k.Address, "500000", INR))
	}

	debugf("Creating offers on the DEX...")
	failOnError0("CreateOffer (usdHolder)", ms.CreateOffer(usdHolder.Seed, USD, XLM, "1", "20"))
	failOnError0("CreateOffer (eurHolder)", ms.CreateOffer(eurHolder.Seed, EUR, USD, "1", "20"))
	failOnError0("CreateOffer (inrHolder)", ms.CreateOffer(inrHolder.Seed, INR, EUR, "1", "20"))

	// Show offers
	showOffersAndBalances(ms, USD, "usdHolder", usdHolder.Address)
	showOffersAndBalances(ms, EUR, "eurHolder", eurHolder.Address)
	showOffersAndBalances(ms, INR, "inrHolder", inrHolder.Address)

	// Check Mary's balance
	debugf("Before path payment...")
	showBalance(ms, INR, "bob", bob.Address)
	showBalance(ms, INR, "mary", mary.Address)

	// Call the path finder
	// logrus.SetLevel(logrus.DebugLevel)
	paths, err := ms.FindPaths(bob.Address, mary.Address, INR, "5")
	debugf("Got: %+v", paths)

	debugf("Executing path payment...")
	err = ms.Pay(bob.Seed, mary.Address, "5", INR,
		microstellar.Opts().WithAsset(XLM, "10").Through(USD, EUR))

	if err != nil {
		log.Fatalf("Pay: %v", microstellar.ErrorString(err))
	}

	debugf("After path payment...")
	showBalance(ms, INR, "bob", bob.Address)
	showBalance(ms, INR, "mary", mary.Address)
}
