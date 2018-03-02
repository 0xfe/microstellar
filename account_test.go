package microstellar

import "testing"

func TestAccounts(t *testing.T) {
	account := &Account{}
	asset := NewAsset("USD", "foobar", Credit4Type)

	account.NativeBalance = Balance{
		Asset:  NativeAsset,
		Amount: "10",
		Limit:  "",
	}

	account.Balances = append(account.Balances, Balance{
		Asset:  asset,
		Amount: "1",
		Limit:  "1000",
	})

	if balance := account.GetNativeBalance(); balance != "10" {
		t.Errorf("wrong native balance: want %v, got %v", "10", balance)
	}

	if balance := account.GetBalance(asset); balance != "1" {
		t.Errorf("wrong native balance: want %v, got %v", "1", balance)
	}
}
