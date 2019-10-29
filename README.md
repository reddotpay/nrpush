# nrpush
Push custom New Relic Insight datasource in golang

## Installation

```
import "github.com/reddotpay/nrpush"
```

## Usage

```go
const Endpoint = "https://insights-collector.newrelic.com/v1/accounts/{:accountID}/events"
```
Endpoint defines the New Relic HTTP URL where the payload will be sent

#### type NRPush

```go
type NRPush struct {
	Endpoint  string
	InsertKey string
	EventType string
}
```

NRPush represents an NewRelic Push data

#### func  New

```go
func New(insertKey, accountID, eventType string) NRPush
```
New creates sets a new NRPush configuration. `accountID` numeric. Can be found
as part of the URL endpoint in Insigts dashboard

#### func (NRPush) Push

```go
func (n NRPush) Push(ctx context.Context, data interface{}) (string, error)
```
Push inserts a new custom data and return the UUID when successful

#### func (NRPush) PushBatch

```go
func (n NRPush) PushBatch(ctx context.Context, data []interface{}) (string, error)
```
PushBatch inserts a batch of new custom data and return the UUID when successful
