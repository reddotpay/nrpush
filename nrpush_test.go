package nrpush_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/reddotpay/nrpush"
)

func ExampleNRPush_Push() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"success": true, "uuid": "73dda6fb-001f-b000-0000-016e157e6878"}`)
	}))

	defer ts.Close()

	n := nrpush.New("somepushkey123", "111111", "transaction")
	n.Endpoint = ts.URL

	// nrpush.Verbose = true
	uuid, err := n.Push(context.Background(), map[string]interface{}{
		"amount":    100.00,
		"product":   "test",
		"eventType": "someotherevent",
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("New Relic UUID: %s", uuid)

	// Output:
	// New Relic UUID: 73dda6fb-001f-b000-0000-016e157e6878
}

func ExampleNRPush_PushBatch() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"success": true, "uuid": "73dda6fb-001f-b000-0000-016e157e6878"}`)
	}))

	defer ts.Close()

	n := nrpush.New("somepushkey123", "111111", "transaction")
	n.Endpoint = ts.URL

	uuid, err := n.PushBatch(context.Background(), []interface{}{
		map[string]interface{}{
			"amount":    100.00,
			"product":   "test",
			"eventType": "someotherevent",
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("New Relic UUID: %s", uuid)

	// Output:
	// New Relic UUID: 73dda6fb-001f-b000-0000-016e157e6878
}
