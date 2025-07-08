package models

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	// ErrInvalidHash represents an error when the hash format is invalid
	ErrInvalidHash = errors.New("the encoded hash is not in the correct format")

	// ErrIncompatibleVersion represents an error when the argon2 version is incompatible
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

// ArgonParams holds the configuration for Argon2i hashing
type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultArgonParams returns the default Argon2i parameters
// These are secure defaults suitable for production use
func DefaultArgonParams() *ArgonParams {
	return &ArgonParams{
		Memory:      64 * 1024, // 64 MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// HashPassword creates an Argon2i hash of a password using the provided parameters
func HashPassword(password string, params *ArgonParams) (string, error) {
	if params == nil {
		params = DefaultArgonParams()
	}

	// Generate a cryptographically secure random salt
	salt, err := generateRandomBytes(params.SaltLength)
	if err != nil {
		return "", err
	}

	// Pass the plaintext password, salt and parameters to the argon2.Key
	// function. This will generate a hash of the password using the Argon2i
	// variant with the provided parameters.
	hash := argon2.Key([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	// Base64 encode the salt and hashed password (using standard encoding, not raw)
	b64Salt := base64.StdEncoding.EncodeToString(salt)
	b64Hash := base64.StdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation
	encodedHash := fmt.Sprintf("$argon2i$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, params.Memory, params.Iterations, params.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword performs password verification by comparing a password with its hash
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded password hash
	params, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters
	otherHash := argon2.Key([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	// Check that the contents of the hashed passwords are identical
	// Use the subtle.ConstantTimeCompare() function for this to help prevent
	// timing attacks
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// generateRandomBytes generates cryptographically secure random bytes
func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// decodeHash extracts the parameters, salt and derived key from an encoded hash
func decodeHash(encodedHash string) (params *ArgonParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	params = &ArgonParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.StdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	hash, err = base64.StdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(hash))

	return params, salt, hash, nil
}
