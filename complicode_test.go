package complicode

import (
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	authCode := "29040011007"
	key := "9rCB7Sv4X29d)5k7N%3ab89p-3(5[A"
	date, _ := time.Parse("20060102", "20070702")
	inv := Invoice{Number: 1503, Nit: 4189179011, Date: date, Amount: 2500}
	code := Generate(authCode, key, inv)

	expected := "6A-DC-53-05-14"

	if code != expected {
		t.Errorf("Generate() = %q, expected %q", code, expected)
	}
}
