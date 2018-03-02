# MicroStellar

MicroStellar is an easy-to-use Go client for the Stellar network.

## QuickStart

```go
// Create a new MicroStellar client connected to the testnet.
ms := microstellar.New("test")

// Generate a new random keypair.
pair, err := ms.CreateKeyPair()
if err != nil { log.Fatalf("CreateKeyPair: %v", err) }

// Display address and key
log.Printf("Private seed: %s, Public address: %s", pair.Seed, pair.Address)

// Fund the account with 1 lumen from an existing account.
err = ms.FundAccount(pair.Address, "S6 ... private key ... 3J", "1")
if err != nil { log.Fatalf("FundAccount: %v", err) }

// Now load account details from ledger.
account, err := ms.LoadAccount(pair.Address)
if err != nil { log.Fatalf("LoadAccount: %v", err) }

// See balance
log.Printf("Native Balance: %v XLM", account.GetNativeBalance())
```

## Status

### Working

* Connect to network and send operations
* Create and fund accounts
* Get balances
* Pay in native and credit assets

### In progress

* Trust lines

### TODO

* Issue Assets
* Add/change signers
* Offer management

## MIT License

... insert MIT blah blah here ...

Copyright Mohit Muthanna Cheppudira 2018