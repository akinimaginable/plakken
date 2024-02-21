package utils_test

import (
	"errors"
	"testing"

	"git.gnous.eu/gnouseu/plakken/internal/utils"
)

func TestCheckCharNotRedundantTrue(t *testing.T) { // Test CheckCharRedundant with redundant char
	want := true
	got := utils.CheckCharRedundant("2d1h3md4h7s", "h")
	if got != want {
		t.Fatal("Error in parseExpirationFull, want : ", want, "got : ", got)
	}
}

func TestCheckCharNotRedundantFalse(t *testing.T) { // Test CheckCharRedundant with not redundant char
	want := false
	got := utils.CheckCharRedundant("2d1h3m47s", "h")
	if got != want {
		t.Fatal("Error in parseExpirationFull, want : ", want, "got : ", got)
	}
}

func TestParseExpirationFull(t *testing.T) { // test parseExpirationFull with all valid separator
	result, _ := utils.ParseExpiration("2d1h3m47s")
	correctValue := 176627
	if result != correctValue {
		t.Fatal("Error in parseExpirationFull, want : ", correctValue, "got : ", result)
	}
}

func TestParseExpirationMissing(t *testing.T) { // test parseExpirationFull with all valid separator
	result, _ := utils.ParseExpiration("1h47s")
	correctValue := 3647
	if result != correctValue {
		t.Fatal("Error in ParseExpirationFull, want : ", correctValue, "got : ", result)
	}
}

func TestParseExpirationWithCaps(t *testing.T) { // test parseExpirationFull with all valid separator
	result, _ := utils.ParseExpiration("2D1h3M47s")
	correctValue := 176627
	if result != correctValue {
		t.Fatal("Error in parseExpirationFull, want : ", correctValue, "got : ", result)
	}
}

func TestParseExpirationNull(t *testing.T) { // test ParseExpirationFull with all valid separator
	result, _ := utils.ParseExpiration("0")
	correctValue := 0
	if result != correctValue {
		t.Fatal("Error in ParseExpirationFull, want: ", correctValue, "got: ", result)
	}
}

func TestParseExpirationNegative(t *testing.T) { // test ParseExpirationFull with all valid separator
	_, got := utils.ParseExpiration("-42h1m4s")
	want := &utils.ParseExpirationError{}
	if !errors.As(got, &want) {
		t.Fatal("Error in ParseExpirationFull, want : ", want, "got : ", got)
	}
}

func TestParseExpirationInvalid(t *testing.T) { // test ParseExpirationFull with all valid separator
	_, got := utils.ParseExpiration("8h42h1m1d4s")
	want := &utils.ParseExpirationError{}
	if !errors.As(got, &want) {
		t.Fatal("Error in ParseExpirationFull, want : ", want, "got : ", got)
	}

}

func TestParseExpirationInvalidRedundant(t *testing.T) { // test ParseExpirationFull with all valid separator
	_, got := utils.ParseExpiration("8h42h1m1h4s")
	want := &utils.ParseExpirationError{}
	if !errors.As(got, &want) {
		t.Fatal("Error in ParseExpirationFull, want : ", want, "got : ", got)
	}
}

func TestValidKey(t *testing.T) { // test ValidKey with a valid key
	got := utils.ValidKey("ab_a-C42")
	want := true

	if got != want {
		t.Fatal("Error in ValidKey, want : ", want, "got : ", got)
	}
}

func TestInValidKey(t *testing.T) { // test ValidKey with invalid key
	got := utils.ValidKey("ab_?a-C42")
	want := false

	if got != want {
		t.Fatal("Error in ValidKey, want : ", want, "got : ", got)
	}
}
