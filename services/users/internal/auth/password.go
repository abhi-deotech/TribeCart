package auth

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
	// ErrInvalidHash is returned when the encoded hash is not in the correct format
	ErrInvalidHash = errors.New("the encoded hash is not in the correct format")
	// ErrIncompatibleVersion is returned when the hash was created with a different version of argon2
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type params struct {
	memory      uint32
time        uint32
	threads    uint8
	keyLen     uint32
	saltLength uint32
}

var defaultParams = &params{
	memory:      64 * 1024, // 64 MB
	time:        3,
	threads:     2,
	keyLen:      32,
	saltLength:  16,
}

// HashPassword generates a new Argon2id hash of the password using the provided parameters
func HashPassword(password string) (string, error) {
	// Generate a cryptographically secure random salt
	salt := make([]byte, defaultParams.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Generate the hash using Argon2id
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		defaultParams.time,
		defaultParams.memory,
		defaultParams.threads,
		defaultParams.keyLen,
	)

	// Encode the hash and salt as a string
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		defaultParams.memory,
		defaultParams.time,
		defaultParams.threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encodedHash, nil
}

// VerifyPassword compares a password with a hash to see if they match
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Extract the parameters, salt, and hash from the encoded hash
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Generate a hash using the same parameters and salt
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		p.time,
		p.memory,
		p.threads,
		p.keyLen,
	)

	// Compare the hashes in constant time to avoid timing attacks
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

// decodeHash decodes the encoded hash into its components
func decodeHash(encodedHash string) (*params, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	// Check the algorithm
	if parts[1] != "argon2id" {
		return nil, nil, nil, fmt.Errorf("unsupported algorithm: %s", parts[1])
	}

	// Parse the version
	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	// Parse the parameters
	p := &params{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.time, &p.threads)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode the salt
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	// Decode the hash
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLen = uint32(len(hash))

	return p, salt, hash, nil
}

// GenerateRandomPassword generates a random password of the given length
func GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+{}[]|;:,.<>?"
	if length < 12 {
		length = 12
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	for i := range bytes {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}

	return string(bytes), nil
}
