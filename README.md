# MicroStellar

MicroStellar is an easy-to-use and fully functional (WIP) Go client for the Stellar network. See [API docs](https://godoc.org/github.com/0xfe/microstellar) for more.

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
pair, _ := ms.CreateKeyPair()
log.Printf("Private seed: %s, Public address: %s", pair.Seed, pair.Address)

// Fund the account with 1 lumen from an existing account.
ms.FundAccount(pair.Address, "S6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", "1")

// Fund an account on the test network with Friendbot.
microstellar.FundWithFriendBot(pair.Address)

// Now load account details from ledger.
account, _ := ms.LoadAccount(pair.Address)
log.Printf("Native Balance: %v XLM", account.GetNativeBalance())

// Pay someone 3 lumens.
ms.PayNative("S6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", "GAUYTZ24ATLEBIV63MXMPOPQO2T6NHI6TQYEXRTFYXWYZ3JOCVO6UYUM", "3")

// Pay someone 1 USD issued by an anchor.
USD := microstellar.NewAsset("USD", "S6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", Credit4Type)
ms.Pay("S6H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", "GAUYTZ24ATLEBIV63MXMPOPQO2T6NHI6TQYEXRTFYXWYZ3JOCVO6UYUM", USD, "3")

// Create a trust line to a credit asset with a limit of 1000.
ms.CreateTrustLine("S4H4HQPE6BRZKLK3QNV6LTD5BGS7S6SZPU3PUGMJDJ26V7YRG3FRNPGA", USD, "10000")

// Check balance.
account, _ = ms.LoadAccount("GAUYTZ24ATLEBIV63MXMPOPQO2T6NHI6TQYEXRTFYXWYZ3JOCVO6UYUM")
log.Printf("USD Balance: %v USD", account.GetBalance(USD))
log.Printf("Native Balance: %v XLM", account.GetNativeBalance())

// What's their home domain?
log.Printf("Home domain: %s", account.HomeDomain)

// Who are the signers on the account?
for i, s := range account.Signers {
    log.Printf("Signer %d (weight: %v): %v", i, s.PublicKey, s.Weight)
}
```

## Documentation

* API Docs - https://godoc.org/github.com/0xfe/microstellar
* End-to-end test - https://github.com/0xfe/microstellar/blob/master/macrotest/macrotest.go

### Supported Features

* Account creation and funding
* Lookup balances, home domain, and account signers
* Payment of native and custom assets
* Add and remove trust lines
* Change key weights

### Coming Soon

* Offer management
* Streaming ledger data (transactions, offers, etc.)
* Operations and options

## Hacking on MicroStellar

### Contribution Guidelines

* We're managing dependencies with [dep](https://github.com/golang/dep).
  * Add a new dependency with `dep ensure -add ...`
* If you're adding a new feature:
  * Add unit tests
  * Add godoc comments
  * If necessary, update the integration test in `macrotest/`
  * If necessary, add examples and verify that they show up in godoc

### Environment Setup

This package uses [dep](https://github.com/golang/dep) to manage dependencies. Before
hacking on this package, install all dependencies:

```
dep ensure
```

### Run tests

Run all unit tests:

```
go test -v
```

Run end-to-end integration test:

```
go run -v macrotest/macrotest.go
```

Test documentation with:

```
godoc -v -http=:6060
```

Then: http://localhost:6060/pkg/github.com/0xfe/microstellar/

## MIT License

Copyright Mohit Muthanna Cheppudira 2018 <mohit@muthanna.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.