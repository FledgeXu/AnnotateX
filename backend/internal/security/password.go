package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// ====== Argon2 Default Config ======
const (
	argonVersion = argon2.Version
	saltLength   = 16
)

var defaultParams = Params{
	Time:    3,
	Memory:  128 * 1024,
	Threads: 4,
	KeyLen:  32,
}

// Params holds tunable Argon2id parameters.
type Params struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// ====== Sentinel Errors ======
var (
	ErrInvalidHashFormat   = errors.New("invalid hash format")
	ErrUnsupportedVersion  = errors.New("unsupported argon2 version")
	ErrInvalidParameters   = errors.New("invalid hash parameters")
	ErrInvalidSaltEncoding = errors.New("invalid salt encoding")
	ErrInvalidHashEncoding = errors.New("invalid hash encoding")
)

// HashPassword generates a new hashed password string in encoded format.
func HashPassword(password string) (string, error) {
	return HashPasswordWithParams(password, defaultParams)
}

// HashPasswordWithParams hashes a password using provided Argon2 parameters.
func HashPasswordWithParams(password string, p Params) (string, error) {
	salt, err := generateSalt(saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, p.KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argonVersion, p.Memory, p.Time, p.Threads, b64Salt, b64Hash,
	), nil
}

// VerifyPassword checks if a password matches the given encoded hash.
func VerifyPassword(password, encodedHash string) (bool, bool, error) {
	parsed, err := decodeHash(encodedHash)
	if err != nil {
		return false, false, err
	}

	computed := argon2.IDKey([]byte(password), parsed.Salt, parsed.Params.Time, parsed.Params.Memory, parsed.Params.Threads, uint32(len(parsed.Hash)))
	match := subtle.ConstantTimeCompare(parsed.Hash, computed) == 1

	needsRehash := parsed.Params != defaultParams
	return match, needsRehash, nil
}

// RehashIfNeeded verifies a password and rehashes it if parameters are outdated.
func RehashIfNeeded(password, encodedHash string) (string, bool, error) {
	match, needsRehash, err := VerifyPassword(password, encodedHash)
	if err != nil || !match {
		return "", false, err
	}
	if !needsRehash {
		return encodedHash, false, nil
	}
	newHash, err := HashPassword(password)
	if err != nil {
		return "", false, err
	}
	return newHash, true, nil
}

// ====== Internal ======

type parsedHash struct {
	Params Params
	Salt   []byte
	Hash   []byte
}

func decodeHash(encoded string) (*parsedHash, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 5 {
		return nil, ErrInvalidHashFormat
	}

	var version int
	var memory uint32
	var time uint32
	var threads uint8

	if _, err := fmt.Sscanf(parts[1], "v=%d", &version); err != nil || version != argonVersion {
		return nil, ErrUnsupportedVersion
	}
	if _, err := fmt.Sscanf(parts[2], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return nil, ErrInvalidParameters
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return nil, ErrInvalidSaltEncoding
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, ErrInvalidHashEncoding
	}

	return &parsedHash{
		Params: Params{
			Time:    time,
			Memory:  memory,
			Threads: threads,
			KeyLen:  uint32(len(hash)),
		},
		Salt: salt,
		Hash: hash,
	}, nil
}

// generateSalt returns a securely generated random salt of given length.
func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	return salt, err
}
