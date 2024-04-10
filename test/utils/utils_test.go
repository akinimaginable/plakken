package utils_test

import (
	"errors"
	"testing"

	"git.gnous.eu/gnouseu/plakken/internal/utils"
)

func TestUtils(t *testing.T) {
	t.Parallel()

	testCheckCharNotRedundantTrue(t)
	testCheckCharNotRedundantFalse(t)
	testParseExpirationFull(t)
	testParseExpirationMissing(t)
	testParseExpirationWithCaps(t)
	testParseExpirationNull(t)
	testParseExpirationNegative(t)
	testParseExpirationInvalid(t)
	testParseExpirationInvalidRedundant(t)
	testParseExpirationInvalidTooHigh(t)
	testValidKey(t)
	testInValidKey(t)
}

func testCheckCharNotRedundantTrue(t *testing.T) { // Test CheckCharRedundant with redundant char
	t.Helper()
	want := true
	got := utils.CheckCharRedundant("2d1h3md4h7s", "h")
	if got != want {
		t.Fatal("Error in parseExpirationFull, want : ", want, "got : ", got)
	}
}

func testCheckCharNotRedundantFalse(t *testing.T) { // Test CheckCharRedundant with not redundant char
	t.Helper()
	want := false
	got := utils.CheckCharRedundant("2d1h3m47s", "h")
	if got != want {
		t.Fatal("Error in parseExpirationFull, want : ", want, "got : ", got)
	}
}

func testParseExpirationFull(t *testing.T) { // test parseExpirationFull with all valid separator
	t.Helper()
	result, _ := utils.ParseExpiration("2d1h3m47s")
	correctValue := 176627
	if result != correctValue {
		t.Fatal("Error in parseExpirationFull, want : ", correctValue, "got : ", result)
	}
}

func testParseExpirationMissing(t *testing.T) { // test parseExpirationFull with all valid separator
	t.Helper()
	result, _ := utils.ParseExpiration("1h47s")
	correctValue := 3647
	if result != correctValue {
		t.Fatal("Error in ParseExpirationFull, want : ", correctValue, "got : ", result)
	}
}

func testParseExpirationWithCaps(t *testing.T) { // test parseExpirationFull with all valid separator
	t.Helper()
	result, _ := utils.ParseExpiration("2D1h3M47s")
	correctValue := 176627
	if result != correctValue {
		t.Fatal("Error in parseExpirationFull, want : ", correctValue, "got : ", result)
	}
}

func testParseExpirationNull(t *testing.T) { // test ParseExpirationFull with all valid separator
	t.Helper()
	result, _ := utils.ParseExpiration("0")
	correctValue := 0
	if result != correctValue {
		t.Fatal("Error in ParseExpirationFull, want: ", correctValue, "got: ", result)
	}
}

func testParseExpirationNegative(t *testing.T) { // test ParseExpirationFull with all valid separator
	t.Helper()
	_, got := utils.ParseExpiration("-42h1m4s")
	want := &utils.ParseExpirationError{}
	if !errors.As(got, &want) {
		t.Fatal("Error in ParseExpirationFull, want : ", want, "got : ", got)
	}
}

func testParseExpirationInvalid(t *testing.T) { // test ParseExpirationFull with all valid separator
	t.Helper()
	_, got := utils.ParseExpiration("8h42h1m1d4s")
	want := &utils.ParseExpirationError{}
	if !errors.As(got, &want) {
		t.Fatal("Error in ParseExpirationFull, want : ", want, "got : ", got)
	}
}

func testParseExpirationInvalidRedundant(t *testing.T) { // test ParseExpirationFull with all valid separator
	t.Helper()
	_, got := utils.ParseExpiration("8h42h1m1h4s")
	want := &utils.ParseExpirationError{}
	if !errors.As(got, &want) {
		t.Fatal("Error in ParseExpirationFull, want : ", want, "got : ", got)
	}
}

func testParseExpirationInvalidTooHigh(t *testing.T) { // test ParseExpirationFull with all valid separator
	t.Helper()
	_, got := utils.ParseExpiration("2d1h3m130s")
	want := &utils.ParseExpirationError{}
	if !errors.As(got, &want) {
		t.Fatal("Error in ParseExpirationFull, want : ", want, "got : ", got)
	}
}

func testValidKey(t *testing.T) { // test ValidKey with a valid key
	t.Helper()
	got := utils.ValidKey("ab_a-C42")
	want := true

	if got != want {
		t.Fatal("Error in ValidKey, want : ", want, "got : ", got)
	}
}

func testInValidKey(t *testing.T) { // test ValidKey with invalid key
	t.Helper()
	got := utils.ValidKey("ab_?a-C42")
	want := false

	if got != want {
		t.Fatal("Error in ValidKey, want : ", want, "got : ", got)
	}
}
