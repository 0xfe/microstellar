package macrotest

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/0xfe/microstellar"
)

func TestMicroStellarOps(t *testing.T) {
	const fundSourceSeed = "SBW2N5EK5MZTKPQJZ6UYXEMCA63AO3AVUR6U5CUOIDFYCAR2X2IJIZAX"
	ms := microstellar.New("test")

	bob := createFundedAccount(ms, fundSourceSeed, true)

	log.Print("Pair1 (bob): ", bob)

	log.Print("Setting a bunch of data on bob...")
	ms.Start(bob.Seed)
	ms.SetData(bob.Address, "foo", []byte("bar"))
	ms.SetData(bob.Address, "fuzz", []byte{0xff, 0xff})
	ms.SetHomeDomain(bob.Address, "qubit.sh")
	ms.SetFlags(bob.Address, microstellar.FlagAuthRequired)
	ms.SetMasterWeight(bob.Address, 3)
	ms.SetThresholds(bob.Address, 1, 1, 1)
	err := ms.Submit()

	if err != nil {
		t.Errorf("Error submitting transaction: %v", microstellar.ErrorString(err))
	}

	log.Print("Reading bob's data...")
	account, err := ms.LoadAccount(bob.Address)

	if err != nil {
		t.Errorf("Error submitting transaction: %v", microstellar.ErrorString(err))
	}

	accountJSON, _ := json.MarshalIndent(*account, "", "  ")
	log.Print(string(accountJSON))

	foo, _ := account.GetData("foo")
	fuzz, _ := account.GetData("fuzz")
	log.Print(fmt.Sprintf("Got data: foo: %v, fuzz: %v", string(foo), fuzz))

	if string(foo) != "bar" {
		t.Errorf("wrong data value for foo: want %v, got %v", "bar", foo)
	}

	if account.HomeDomain != "qubit.sh" {
		t.Errorf("wrong home domain: want %v, got %v", "qubit.sh", account.HomeDomain)
	}

	if account.Flags.AuthRequired != true {
		t.Errorf("expecting account.Flags.AuthRequired to be true")
	}

	log.Print("Changing a bunch of bob's data...")
	ms.Start(bob.Seed)
	ms.ClearData(bob.Address, "foo")
	ms.ClearFlags(bob.Address, microstellar.FlagAuthRequired)
	err = ms.Submit()

	if err != nil {
		t.Errorf("Error submitting transaction: %v", microstellar.ErrorString(err))
	}

	log.Print("Reading bob's data...")
	account, err = ms.LoadAccount(bob.Address)

	if err != nil {
		t.Errorf("Error submitting transaction: %v", microstellar.ErrorString(err))
	}

	accountJSON, _ = json.MarshalIndent(*account, "", "  ")
	log.Print(string(accountJSON))

	_, ok := account.Data["foo"]

	if ok {
		t.Errorf("data key foo should not exist")
	}

	if account.Flags.AuthRequired == true {
		t.Errorf("expecting account.Flags.AuthRequired to be false")
	}
}
