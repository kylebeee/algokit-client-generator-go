# algokit-client-generator-go

CLI tool that generates type-safe Go client packages from [ARC-56](https://arc.algorand.foundation/ARCs/arc-0056)/[ARC-32](https://arc.algorand.foundation/ARCs/arc-0032) Algorand application specifications.

## Installation

```bash
go install github.com/kylebeee/algokit-client-generator-go@latest
```

## Usage

```bash
algokit-client-generator-go generate \
  --application path/to/app.arc56.json \
  --output ./mycontract \
  --package mycontract
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--application` | `-a` | Path to ARC-56 or ARC-32 app spec JSON file (required) |
| `--output` | `-o` | Output directory for the generated Go package (required) |
| `--package` | `-p` | Go package name (default: derived from contract name) |
| `--mode` | `-m` | Generation mode: `full` or `minimal` (default: `full`) |
| `--preserve-names` | | Preserve original method names without sanitization |

## Generated Output

The generator produces 5 files per contract:

| File | Contents |
|------|----------|
| `appspec.go` | Embedded ARC-56 JSON spec with `GetAppSpec()` helper |
| `types.go` | Argument structs, result structs, and ABI struct types |
| `client.go` | `Client` with `Send{Method}()` methods for each ABI call |
| `composer.go` | `Composer` for building atomic transaction groups |
| `factory.go` | `Factory` for deploying new contract instances |

## Example

### Deploy and call a contract

```go
package main

import (
    "context"
    "fmt"

    algokit "github.com/kylebeee/algokit-utils-go"
    "example.com/myapp/gate"
)

func main() {
    ctx := context.Background()

    // Connect to LocalNet
    client, _ := algokit.LocalNet()
    account, _ := client.Account().Random()
    client.Account().SetSignerFromAccount(account)
    signer := algokit.BasicAccountTransactionSigner(account)

    // Create a typed factory
    factory, _ := gate.NewFactory(algokit.AppFactoryParams{
        Algorand:      client,
        DefaultSender: account.Address,
        DefaultSigner: signer,
    })

    // Deploy a new instance with typed create args
    gateClient, result, _ := factory.Create(ctx, gate.FactoryCreateParams{
        Args: gate.CreateArgs{
            Version:  "v1.0",
            AkitaDao: 0,
        },
    })
    fmt.Printf("Deployed app %d (tx: %s)\n", gateClient.AppID(), result.TxID)

    // Call a method with typed args and typed result
    checkResult, _ := gateClient.SendCheck(ctx, algokit.CallParams[gate.CheckArgs]{
        Args: gate.CheckArgs{
            Caller: account.Address,
            GateID: 1,
            Args:   [][]byte{},
        },
        Sender: account.Address,
        Signer: signer,
    })
    fmt.Printf("Check passed: %v\n", checkResult.Return) // bool
}
```

### Compose atomic transaction groups

```go
composer := gateClient.NewComposer()
composer.Register(ctx, algokit.CallParams[gate.RegisterArgs]{
    Args: gate.RegisterArgs{Payment: paymentTxn, Filters: filters, Args: args},
    Sender: account.Address,
    Signer: signer,
})
composer.Check(ctx, algokit.CallParams[gate.CheckArgs]{
    Args: gate.CheckArgs{Caller: account.Address, GateID: 1, Args: [][]byte{}},
    Sender: account.Address,
    Signer: signer,
})
result, _ := composer.Send(ctx)
fmt.Printf("Group confirmed in round %d\n", result.ConfirmedRound)
```

## Requirements

Generated code depends on [algokit-utils-go](https://github.com/kylebeee/algokit-utils-go) at runtime:

```bash
go get github.com/kylebeee/algokit-utils-go
```
