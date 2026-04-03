# go-edgeservpos

[![CI](https://github.com/jeffresc/go-edgeservpos/actions/workflows/ci.yml/badge.svg)](https://github.com/jeffresc/go-edgeservpos/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/jeffresc/go-edgeservpos/branch/main/graph/badge.svg)](https://codecov.io/gh/jeffresc/go-edgeservpos)
[![Go Reference](https://pkg.go.dev/badge/github.com/jeffresc/go-edgeservpos.svg)](https://pkg.go.dev/github.com/jeffresc/go-edgeservpos)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go client library for the [EdgeServ POS](https://edgeserv.com/) API. Provides programmatic access to restaurant point-of-sale data including customer management.

## Installation

```bash
go get github.com/jeffresc/go-edgeservpos
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	edgeservpos "github.com/jeffresc/go-edgeservpos"
)

func main() {
	client := edgeservpos.NewClient(
		"https://api.example.com", // API host
		"RESTAURANT_CODE",         // Restaurant code
		"CLIENT_ID",               // OAuth client ID
		"CLIENT_SECRET",           // OAuth client secret
		"USERNAME",                // Username
		"PASSWORD",                // Password
	)

	customers, err := client.ListCustomers()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range customers {
		fmt.Printf("%s %s (%s)\n", c.FirstName, c.LastName, c.EmailAddress)
	}
}
```

## API

### Client

`NewClient(host, restaurantCode, clientID, clientSecret, username, password string) *Client` — Creates a new API client.

### Methods

| Method | Description |
|---|---|
| `GetOAuthToken() (string, error)` | Retrieves an OAuth bearer token. |
| `ListCustomers() ([]Customer, error)` | Retrieves all customers for the restaurant. Handles token management automatically. |

### Types

- **Customer** — `ServerID`, `FirstName`, `LastName`, `EmailAddress`, `Point`, `PhoneNumbers`, `LastVisitDate`, `Addresses`
- **Address** — `Address`, `Address2`, `City`, `State`, `ZipCode`

## Testing

```bash
go test -race ./...
```

## License

[MIT](LICENSE)
