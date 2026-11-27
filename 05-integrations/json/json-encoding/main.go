// Topic: json-encoding | Concept: json | Stage: 05-integrations
// Run:   go run ./05-integrations/json/json-encoding
//
// Explore:
//   - Only EXPORTED (capitalised) fields are encoded
//   - Struct tags rename fields: `json:"user_name"`; `omitempty` drops zero values;
//     `json:"-"` never encodes the field (passwords!); `json:",string"` quotes numbers
//   - time.Time encodes as RFC 3339 text
//   - Implement MarshalJSON / UnmarshalJSON to control a type's JSON form
//   - Unknown shape? Decode into map[string]any (numbers become float64)
//   - json.RawMessage delays decoding: read a "type" field first, then decode the rest
//   - Decoder.DisallowUnknownFields catches typos in input; Decoder also streams
//     many JSON values one after another (NDJSON logs)
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Money stores cents as an int but appears in JSON as a "12.34" string.
type Money int64

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%d.%02d", m/100, m%100))
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("money must be a string like \"12.34\": %w", err)
	}
	var whole, cents int64
	if _, err := fmt.Sscanf(s, "%d.%02d", &whole, &cents); err != nil {
		return fmt.Errorf("bad money %q: %w", s, err)
	}
	*m = Money(whole*100 + cents)
	return nil
}

type Order struct {
	ID        int       `json:"id"`
	Customer  string    `json:"customer_name"`
	Total     Money     `json:"total"`
	Items     []string  `json:"items,omitempty"`
	Note      string    `json:"note,omitempty"`
	Quantity  int       `json:"quantity,string"`
	CreatedAt time.Time `json:"created_at"`
	Secret    string    `json:"-"`
	internal  string    // unexported: never encoded
}

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"` // decoded later, once we know the type
}

type Login struct {
	User string `json:"user"`
}
type Purchase struct {
	Amount Money `json:"amount"`
}

func main() {
	fmt.Println("Topic: encoding/json")
	failed := false

	o := Order{
		ID: 7, Customer: "Asha", Total: 123456, Items: []string{"keyboard", "mouse"},
		Quantity: 2, CreatedAt: time.Date(2026, 11, 27, 10, 30, 0, 0, time.UTC),
		Secret: "do-not-leak", internal: "also hidden",
	}
	data, _ := json.Marshal(o)
	fmt.Println("Marshal:      ", string(data))
	pretty, _ := json.MarshalIndent(Order{ID: 8, Customer: "Ravi"}, "", "  ")
	fmt.Println("MarshalIndent (empty items/note omitted):\n" + string(pretty))
	if strings.Contains(string(data), "do-not-leak") || !strings.Contains(string(data), `"total":"1234.56"`) {
		failed = true
	}

	var back Order
	if err := json.Unmarshal(data, &back); err != nil {
		fmt.Fprintln(os.Stderr, "unmarshal:", err)
		os.Exit(1)
	}
	fmt.Printf("Unmarshal:     id=%d customer=%s total(cents)=%d qty=%d created=%s\n",
		back.ID, back.Customer, back.Total, back.Quantity, back.CreatedAt.Format(time.DateOnly))
	if back.Total != o.Total || back.Quantity != 2 {
		failed = true
	}

	// Unknown shape
	var anything map[string]any
	json.Unmarshal([]byte(`{"name":"sensor-1","reading":21.5,"tags":["a","b"],"ok":true}`), &anything)
	fmt.Printf("map[string]any: %v (reading is %T)\n", anything, anything["reading"])

	// Polymorphic events with RawMessage
	stream := `{"type":"login","data":{"user":"asha"}}
{"type":"purchase","data":{"amount":"49.99"}}
{"type":"logout","data":{}}`
	dec := json.NewDecoder(strings.NewReader(stream)) // streams one value at a time
	for {
		var ev Event
		err := dec.Decode(&ev)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "decode:", err)
			os.Exit(1)
		}
		switch ev.Type {
		case "login":
			var l Login
			json.Unmarshal(ev.Data, &l)
			fmt.Println("  event login    user =", l.User)
		case "purchase":
			var p Purchase
			json.Unmarshal(ev.Data, &p)
			fmt.Println("  event purchase cents =", int64(p.Amount))
			if p.Amount != 4999 {
				failed = true
			}
		default:
			fmt.Println("  event", ev.Type, "(ignored)")
		}
	}

	// Strict decoding catches a typo in the input
	strict := json.NewDecoder(strings.NewReader(`{"id":1,"customer_nmae":"typo"}`))
	strict.DisallowUnknownFields()
	var typo Order
	err := strict.Decode(&typo)
	fmt.Println("DisallowUnknownFields:", err)
	if err == nil {
		failed = true
	}

	// Custom UnmarshalJSON rejects bad input
	var bad Money
	err = json.Unmarshal([]byte(`12.34`), &bad)
	fmt.Println("Money from a bare number:", err)

	if failed {
		fmt.Println("MISMATCH")
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
