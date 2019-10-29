package nrpush

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// NRPush represents an NewRelic Push data
type NRPush struct {
	Endpoint  string
	InsertKey string
	EventType string
}

// Endpoint defines the New Relic HTTP URL where the payload will be sent
const Endpoint = "https://insights-collector.newrelic.com/v1/accounts/{:accountID}/events"

// New creates sets a new NRPush configuration.
// `accountID` numeric. Can be found as part of the URL endpoint in Insigts dashboard
func New(insertKey, accountID, eventType string) NRPush {
	return NRPush{
		Endpoint:  strings.Replace(Endpoint, "{:accountID}", accountID, 1),
		InsertKey: insertKey,
		EventType: eventType,
	}
}

func interfaceToMap(data interface{}) (map[string]interface{}, error) {
	var (
		body   map[string]interface{}
		b, err = json.Marshal(data)
	)

	if err != nil {
		return map[string]interface{}{}, fmt.Errorf("JSON transformation Error(1) - %s", err.Error())
	}

	if err = json.Unmarshal(b, &body); err != nil {
		return map[string]interface{}{}, fmt.Errorf("JSON transformation Error(2)- %s", err.Error())
	}

	return body, nil
}

func mustJSONMarshal(d interface{}) []byte {
	b, err := json.Marshal(d)

	if err != nil {
		panic(fmt.Errorf("mustJSONMarshal - %s", err.Error()))
	}

	return b
}

func mapSafeReplace(m map[string]interface{}, key string, replace interface{}) map[string]interface{} {
	if _, ok := m[key]; ok {
		m["_"+key] = m[key]
	}
	m[key] = replace
	return m
}

// PushBatch inserts a batch of new custom data and return the UUID when successful
func (n NRPush) PushBatch(ctx context.Context, data []interface{}) (string, error) {
	for i, v := range data {
		body, err := interfaceToMap(v)
		if err != nil {
			return "", err
		}

		data[i] = mapSafeReplace(body, "eventType", n.EventType)
	}

	return n.sendWithContext(ctx, mustJSONMarshal(data))
}

// Push inserts a new custom data and return the UUID when successful
func (n NRPush) Push(ctx context.Context, data interface{}) (string, error) {
	var (
		body map[string]interface{}
		err  error
	)

	if body, err = interfaceToMap(data); err != nil {
		return "", err
	}

	// Put in eventType property
	body = mapSafeReplace(body, "eventType", n.EventType)

	return n.sendWithContext(ctx, mustJSONMarshal(body))
}

func (n NRPush) sendWithContext(ctx context.Context, data []byte) (string, error) {
	var (
		client       = http.Client{}
		request      *http.Request
		response     *http.Response
		body         map[string]interface{}
		responseBody []byte
		err          error
	)

	// HTTP
	request, err = http.NewRequest(http.MethodPost, n.Endpoint, bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Insert-Key", n.InsertKey)
	request.WithContext(ctx)

	if err != nil {
		return "", err
	}

	if response, err = client.Do(request); err != nil {
		return "", err
	}

	defer response.Body.Close()

	if responseBody, err = ioutil.ReadAll(response.Body); err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &body)
	return body["uuid"].(string), err
}
