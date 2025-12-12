package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// Number of iterations for PBKDF2
	pbkdf2Iterations = 100000
	// Key size for AES-256
	keySize = 32
	// Salt size for key derivation
	saltSize = 16
)

var (
	// ErrInvalidCiphertext is returned when the ciphertext is too short or invalid
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	// ErrDecryptionFailed is returned when decryption fails
	ErrDecryptionFailed = errors.New("decryption failed")
)

// GetMachineID generates a machine-specific identifier for key derivation.
// This ensures that encrypted data is tied to the specific machine.
func GetMachineID() (string, error) {
	// Use hostname as the primary identifier
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("failed to get hostname: %w", err)
	}

	// Combine hostname with OS and architecture for uniqueness
	machineID := fmt.Sprintf("%s-%s-%s", hostname, runtime.GOOS, runtime.GOARCH)

	return machineID, nil
}

// DeriveKey derives a cryptographic key from a machine ID using PBKDF2.
// The salt is stored with the ciphertext to allow key derivation during decryption.
func DeriveKey(machineID string, salt []byte) []byte {
	return pbkdf2.Key([]byte(machineID), salt, pbkdf2Iterations, keySize, sha256.New)
}

// Encrypt encrypts plaintext using AES-256-GCM with a machine-specific key.
// The output format is: [salt(16 bytes)][nonce(12 bytes)][ciphertext+tag]
// Returns base64-encoded result for safe storage in database.
func Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// Get machine ID for key derivation
	machineID, err := GetMachineID()
	if err != nil {
		return "", fmt.Errorf("failed to get machine ID: %w", err)
	}

	// Generate random salt
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive encryption key
	key := DeriveKey(machineID, salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Combine salt, nonce, and ciphertext
	result := make([]byte, 0, saltSize+len(nonce)+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	// Encode to base64 for safe storage
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts ciphertext that was encrypted with Encrypt.
// The input must be base64-encoded and contain: [salt][nonce][ciphertext+tag]
func Decrypt(ciphertextBase64 string) (string, error) {
	if ciphertextBase64 == "" {
		return "", nil
	}

	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Check minimum length: salt + nonce + tag (at least)
	if len(data) < saltSize+12+16 {
		return "", ErrInvalidCiphertext
	}

	// Extract salt
	salt := data[:saltSize]

	// Get machine ID for key derivation
	machineID, err := GetMachineID()
	if err != nil {
		return "", fmt.Errorf("failed to get machine ID: %w", err)
	}

	// Derive decryption key
	key := DeriveKey(machineID, salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce and ciphertext
	nonceSize := gcm.NonceSize()
	if len(data) < saltSize+nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	// Decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	return string(plaintext), nil
}

// IsEncrypted checks if a value appears to be encrypted (base64 with correct length).
// This is a heuristic check and may have false positives for very short values.
func IsEncrypted(value string) bool {
	if value == "" {
		return false
	}

	// Try to decode as base64
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return false
	}

	// Check if it has at least the minimum length for encrypted data
	// salt(16) + nonce(12) + tag(16) = 44 bytes minimum
	return len(data) >= saltSize+12+16
}
