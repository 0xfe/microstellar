# MicroStellar

MicroStellar is an easy-to-use Go client for the Stellar network. See [API docs](https://godoc.org/github.com/0xfe/microstellar) for more.

## QuickStart

### Installation

```
go get github.com/0xfe/microstellar
```

### Usage

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

## Documentation

* API Docs - https://godoc.org/github.com/0xfe/microstellar

## Hacking

### Contribution Guidelines

* We're managing dependencies with [dep](https://github.com/golang/dep).
  * Add a new dependency with `dep ensure -add ...`
* If you're adding a new feature:
  * Add unit tests
  * Add godoc comments
  * If necessary, update the integration test in `macrotest/`
  * If necessary, add examples and verify that they show up in godoc

### Setup

This package uses [dep](https://github.com/golang/dep) to manage dependencies. To install dependencies:

```
dep ensure
```

### Run tests

```
go test -v
```

### GoDoc

Test your documentation with:

```
godoc -v -http=:6060
```

Then: http://localhost:6060/pkg/github.com/0xfe/microstellar/

## Status

### Working

* Connect to network and send operations
* Create and fund accounts
* Get balances
* Pay in native and credit assets
* Trust lines
* Issue Assets
* Integration test against horizon-test

### In progress

* Remove trust lines

### TODO

* Add/change signers
* Offer management

## MIT License

... insert MIT blah blah here ...

Copyright Mohit Muthanna Cheppudira 2018