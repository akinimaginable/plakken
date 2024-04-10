package secret_test

import (
	"fmt"
	"regexp"
	"testing"

	"git.gnous.eu/gnouseu/plakken/internal/constant"
	"git.gnous.eu/gnouseu/plakken/internal/secret"
	"golang.org/x/crypto/argon2"
)

func TestSecret(t *testing.T) {
	t.Parallel()

	testPasswordFormat(t)
	testVerifyPassword(t)
	testVerifyPasswordInvalid(t)
}

func testPasswordFormat(t *testing.T) {
	t.Helper()
	regex := fmt.Sprintf("\\$argon2id\\$v=%d\\$m=%d,t=%d,p=%d\\$[A-Za-z0-9+/]*\\$[A-Za-z0-9+/]*$", argon2.Version, constant.ArgonMemory, constant.ArgonIterations, constant.ArgonThreads) //nolint:lll

	got, err := secret.Password("Password!")
	if err != nil {
		t.Fatal(err)
	}

	result, _ := regexp.MatchString(regex, got)
	if !result {
		t.Fatal("Error in Password, format is not valid "+": ", got)
	}
}

func testVerifyPassword(t *testing.T) {
	t.Helper()
	result, err := secret.VerifyPassword("Password!", "$argon2id$v=19$m=65536,t=2,p=4$A+t5YGpyy1BHCbvk/LP1xQ$eNuUj6B2ZqXlGi6KEqep39a7N4nysUIojuPXye+Ypp0") //nolint:lll
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Fatal("Error in VerifyPassword, got:", result, "want: ", true)
	}
}

func testVerifyPasswordInvalid(t *testing.T) {
	t.Helper()
	result, err := secret.VerifyPassword("notsamepassword", "$argon2id$v=19$m=65536,t=2,p=4$A+t5YGpyy1BHCbvk/LP1xQ$eNuUj6B2ZqXlGi6KEqep39a7N4nysUIojuPXye+Ypp0") //nolint:lll
	if err != nil {
		t.Fatal(err)
	}

	if result {
		t.Fatal("Error in VerifyPassword, got:", result, "want: ", false)
	}
}
