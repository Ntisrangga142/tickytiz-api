package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashConfig struct {
	Memory  uint32
	Time    uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

func (h *HashConfig) GenHash(password string) (any, any) {
	panic("unimplemented")
}

func NewHashConfig() *HashConfig {
	return &HashConfig{}
}

func (h *HashConfig) SetConfig(memory, time, keylen, saltlen uint32, threads uint8) {
	h.KeyLen = keylen
	h.SaltLen = saltlen
	h.Memory = memory
	h.Time = time
	h.Threads = threads
}

func (h *HashConfig) UseRecommended() {
	h.KeyLen = 32
	h.SaltLen = 16
	h.Memory = 64 * 1024
	h.Time = 2
	h.Threads = 1
}

func (h *HashConfig) GenerateHash(password string) (string, error) {
	salt, err := h.generateSalt()
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, h.Time, h.Memory, h.Threads, h.KeyLen)

	// Format Penulisan Hash
	// $jenisKey$versiKey$konfigurasi(memory, time, thread)$salt$hash
	version := argon2.Version
	saltStr := base64.RawStdEncoding.EncodeToString(salt)
	hashStr := base64.RawStdEncoding.EncodeToString(hash)
	hashedPwd := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", version, h.Memory, h.Time, h.Threads, saltStr, hashStr)
	return hashedPwd, nil
}

func (h *HashConfig) generateSalt() ([]byte, error) {
	salt := make([]byte, h.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func (h *HashConfig) CompareHashAndPassword(password, hashedPassword string) (bool, error) {
	result := strings.Split(hashedPassword, "$")
	// Cek Panjang Hasil Split
	if len(result) != 6 {
		return false, errors.New("invalid hash format")
	}

	// Cek Cryptography
	if result[1] != "argon2id" {
		return false, errors.New("invalid crypto method")
	}

	// Cek Versi
	var version int
	if _, err := fmt.Sscanf(result[2], "v=%d", &version); err != nil {
		return false, errors.New("invalid version format")
	}

	// Ambil Konfigurasi Hasing dari HashedPassword
	if _, err := fmt.Sscanf(result[3], "m=%d,t=%d,p=%d", &h.Memory, &h.Time, &h.Threads); err != nil {
		return false, errors.New("invalid format")
	}

	// Ambil Nilai Salt
	salt, err := base64.RawStdEncoding.DecodeString(result[4])
	if err != nil {
		return false, err
	}
	h.SaltLen = uint32(len(salt))

	// Ambil Nilai Hash
	hash, err := base64.RawStdEncoding.DecodeString(result[5])
	if err != nil {
		return false, err
	}
	h.KeyLen = uint32(len(hash))

	// Comparison
	// Generate Hash dari password
	hashPwd := argon2.IDKey([]byte(password), salt, h.Time, h.Memory, h.Threads, h.KeyLen)

	// Compare Hasil Hash dengan Waktu Konstan (Lebih Aman dari Timing Attack di Hash)
	if subtle.ConstantTimeCompare(hash, hashPwd) == 0 {
		return false, nil
	}
	return true, nil
}
