// Topic: json-encoding | Concept: io | Stage: 03-advanced
// Run:   go run ./03-advanced/io/json-encoding
//
// Explore:
// * json.Marshal / json.Unmarshal between structs and JSON bytes
// * Struct tags: `json:"name"`, `json:",omitempty"`, `json:"-"`
// * Only exported (Capitalised) fields are encoded
// * Unknown JSON fields are ignored; missing ones keep zero values
// * Decoding into map[string]any when the shape is unknown (numbers become float64)
// * json.NewEncoder / NewDecoder stream to and from any io.Writer / io.Reader
// * Custom formats with MarshalJSON
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Alert struct {
	ID       int       `json:"id"`
	Customer string    `json:"customer"`
	Channels []string  `json:"channels,omitempty"`
	Sent     bool      `json:"sent"`
	Secret   string    `json:"-"`
	internal string    // unexported: never encoded
	Created  Timestamp `json:"created"`
}

// Timestamp shows a custom JSON format
type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format("2006-01-02"))
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := time.Parse("2006-01-02", s)
	*t = Timestamp(parsed)
	return err
}

func main() {
	fmt.Println("Topic: encoding/json")

	a := Alert{
		ID:       7,
		Customer: "Asha",
		Channels: []string{"email", "sms"},
		Secret:   "token-123",
		internal: "hidden",
		Created:  Timestamp(time.Date(2026, 11, 27, 0, 0, 0, 0, time.UTC)),
	}

	b, _ := json.Marshal(a)
	fmt.Println("Marshal:", string(b))

	pretty, _ := json.MarshalIndent(Alert{ID: 8, Customer: "Ravi"}, "", "  ")
	fmt.Println("MarshalIndent (omitempty dropped channels):")
	fmt.Println(string(pretty))

	input := `{"id": 9, "customer": "Meera", "sent": true, "created": "2026-01-15", "extra": "ignored"}`
	var decoded Alert
	if err := json.Unmarshal([]byte(input), &decoded); err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}
	fmt.Printf("Unmarshal: id=%d customer=%s sent=%v created=%s\n",
		decoded.ID, decoded.Customer, decoded.Sent, time.Time(decoded.Created).Format("Jan 2 2006"))

	var generic map[string]any
	_ = json.Unmarshal([]byte(`{"count": 3, "tags": ["a", "b"]}`), &generic)
	fmt.Printf("map[string]any: count is %T (%v), tags is %T\n", generic["count"], generic["count"], generic["tags"])

	var bad Alert
	err := json.Unmarshal([]byte(`{"id": "not-a-number"}`), &bad)
	fmt.Println("type mismatch error:", err)

	// Streaming: decode many objects from a reader, encode to stdout
	dec := json.NewDecoder(strings.NewReader(`{"id":1,"customer":"x"} {"id":2,"customer":"y"}`))
	enc := json.NewEncoder(os.Stdout)
	for dec.More() {
		var al Alert
		if err := dec.Decode(&al); err != nil {
			break
		}
		al.Sent = true
		fmt.Print("  stream out: ")
		_ = enc.Encode(map[string]any{"id": al.ID, "sent": al.Sent})
	}
}
