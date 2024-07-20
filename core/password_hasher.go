package core

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	// HashPassword generates a hash from the given password.
	//
	// It should return the hashed password and an error if it occurs
	HashPassword(plainPassword string) (string, error)

	// VerifyPassword compares the possible plaintext password with the provided hash.
	//
	// It should return a boolean value indicating if plain text password matches with hash and an error if it occurs
	VerifyPassword(plainPassword, hashedPassword string) (bool, error)
}

type bcriptPasswordHasher struct {
	cost int
}

func (h *bcriptPasswordHasher) HashPassword(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), h.cost)
	return string(hash), err
}

func (h *bcriptPasswordHasher) VerifyPassword(plainPassword, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil, nil
}

func NewBcriptPasswordHasher() PasswordHasher {
	return &bcriptPasswordHasher{cost: bcrypt.DefaultCost}
}

type argon2IDPasswordHasher struct {
	format            string
	version           int
	memoryCost        uint32
	iterations        uint32
	parallelismFactor uint8
	hashLength        uint32
	keyLen            uint32
}

func (h *argon2IDPasswordHasher) HashPassword(plainPassword string) (string, error) {
	salt := make([]byte, h.hashLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(plainPassword), salt, h.iterations, h.memoryCost, h.parallelismFactor, h.keyLen)
	hash := fmt.Sprintf(h.format, h.version, h.memoryCost, h.iterations, h.parallelismFactor, base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(key))

	return hash, nil
}

func (h *argon2IDPasswordHasher) VerifyPassword(plainPassword, hashedPassword string) (bool, error) {
	hashParts := strings.Split(hashedPassword, "$")

	if len(hashParts) != 6 {
		return false, errors.New("the encoded hash is not in the correct format")
	}

	var version int
	if _, err := fmt.Sscanf(hashParts[2], "v=%d", &version); err != nil {
		return false, err
	}

	if version != h.version {
		return false, errors.New("incompatible version of argon2")
	}

	var memory, time, threads uint
	_, err := fmt.Sscanf(hashParts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(hashParts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(hashParts[5])
	if err != nil {
		return false, err
	}

	keyLen := len(decodedHash)
	hashToCompare := argon2.IDKey([]byte(plainPassword), salt, uint32(time), uint32(memory), uint8(threads), uint32(keyLen))

	return subtle.ConstantTimeCompare(decodedHash, hashToCompare) == 1, nil
}

func NewArgon2IDPasswordHasher() PasswordHasher {
	return &argon2IDPasswordHasher{
		format:            "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		version:           argon2.Version,
		memoryCost:        64 * 1024,
		iterations:        1,
		parallelismFactor: 4,

		hashLength: 16,
		keyLen:     32,
	}
}
