// Package secret implement all crypto utils like password hashing and secret generation
package secret

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"git.gnous.eu/gnouseu/plakken/internal/constant"
	"golang.org/x/crypto/argon2"
)

type argon2idHash struct {
	salt []byte
	hash []byte
}

// Argon2id config.
type config struct {
	saltLength uint8
	memory     uint32
	threads    uint8
	keyLength  uint32
	iterations uint32
}

// generateSecret for password hashing or token generation.
func generateSecret(length uint8) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, err
}

// GenerateToken generate hexadecimal token.
func GenerateToken() (string, error) {
	secret, err := generateSecret(constant.TokenLength)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(secret)

	return token, nil
}

// generateArgon2ID Generate an argon2id hash from source string and specified salt.
func (config config) generateArgon2ID(source string, salt []byte) []byte {
	hash := argon2.IDKey([]byte(source), salt, config.iterations, config.memory, config.threads, config.keyLength)

	return hash
}

// Password hash a source string with argon2id, return properly encoded hash.
func Password(password string) (string, error) {
	config := config{
		saltLength: constant.ArgonSaltSize,
		memory:     constant.ArgonMemory,
		threads:    constant.ArgonThreads,
		keyLength:  constant.ArgonKeyLength,
		iterations: constant.ArgonIterations,
	}

	salt, err := generateSecret(config.saltLength)
	if err != nil {
		return "", err
	}

	hash := config.generateArgon2ID(password, salt)

	base64Hash := base64.RawStdEncoding.EncodeToString(hash)
	base64Salt := base64.RawStdEncoding.EncodeToString(salt)

	formatted := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, config.memory, config.iterations, config.threads, base64Salt, base64Hash)

	return formatted, nil
}

// VerifyPassword check is source password and stored password is similar, take password and a properly encoded hash.
func VerifyPassword(password string, hash string) (bool, error) {
	argon2Hash, config, err := parseHash(hash)
	if err != nil {
		return false, err
	}

	result := config.generateArgon2ID(password, argon2Hash.salt)

	return bytes.Equal(result, argon2Hash.hash), nil
}

// parseHash parse existing encoded argon2id string.
func parseHash(source string) (argon2idHash, config, error) {
	separateItem := strings.Split(source, "$")
	if len(separateItem) != 6 { //nolint:mnd
		return argon2idHash{}, config{}, &parseError{message: "Hash format is not valid"}
	}

	if separateItem[1] != "argon2id" {
		return argon2idHash{}, config{}, &parseError{message: "Algorithm is not valid"}
	}

	separateParam := strings.Split(separateItem[3], ",")
	if len(separateParam) != 3 { //nolint:mnd
		return argon2idHash{}, config{}, &parseError{message: "Hash config is not valid"}
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(separateItem[4])
	if err != nil {
		return argon2idHash{}, config{}, err
	}

	var hash []byte
	hash, err = base64.RawStdEncoding.Strict().DecodeString(separateItem[5])
	if err != nil {
		return argon2idHash{}, config{}, err
	}

	saltLength := uint8(len(salt))
	keyLength := uint32(len(hash))

	var memory int
	memory, err = strconv.Atoi(strings.ReplaceAll(separateParam[0], "m=", ""))
	if err != nil {
		return argon2idHash{}, config{}, err
	}

	var iterations int
	iterations, err = strconv.Atoi(strings.ReplaceAll(separateParam[1], "t=", ""))
	if err != nil {
		return argon2idHash{}, config{}, err
	}

	var threads int
	threads, err = strconv.Atoi(strings.ReplaceAll(separateParam[2], "p=", ""))
	if err != nil {
		return argon2idHash{}, config{}, err
	}

	argon2idStruct := argon2idHash{
		salt: salt,
		hash: hash,
	}

	hashConfig := config{
		saltLength: saltLength,
		memory:     uint32(memory),
		threads:    uint8(threads),
		iterations: uint32(iterations),
		keyLength:  keyLength,
	}

	return argon2idStruct, hashConfig, nil
}
