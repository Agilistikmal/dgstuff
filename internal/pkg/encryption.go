package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/argon2"
)

const (
	saltSize  = 16
	keySize   = 32
	nonceSize = 12
)

func GenerateKeyFromSecret(secretKey, salt []byte) []byte {
	return argon2.IDKey(secretKey, salt, 1, 64*1024, 4, keySize)
}

func Encrypt(value string, secretKey string) (string, error) {
	if secretKey == "" {
		logrus.Error("Encryption failed: secret key is required")
		return "", errors.New("secret key is required")
	}

	data := []byte(value)
	secretKeyBytes := []byte(secretKey)

	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		logrus.Errorf("Encryption failed: failed to generate salt: %v", err)
		return "", errors.New("failed to generate salt")
	}

	key := GenerateKeyFromSecret(secretKeyBytes, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		logrus.Errorf("Encryption failed: failed to create cipher block: %v", err)
		return "", errors.New("failed to create cipher block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logrus.Errorf("Encryption failed: failed to create GCM: %v", err)
		return "", errors.New("failed to create GCM")
	}

	nonce := make([]byte, nonceSize)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logrus.Errorf("Encryption failed: failed to generate nonce: %v", err)
		return "", errors.New("failed to generate nonce")
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)

	result := append(salt, nonce...)
	result = append(result, ciphertext...)

	return base64.StdEncoding.EncodeToString(result), nil
}

func Decrypt(encodedValue string, secretKey string) (string, error) {
	if secretKey == "" {
		logrus.Error("Decryption failed: secret key is required")
		return "", errors.New("secret key is required")
	}

	data, err := base64.StdEncoding.DecodeString(encodedValue)
	if err != nil {
		logrus.Errorf("Decoding failed: illegal base64 data: %v", err)
		return "", fmt.Errorf("failed to decode value: %w", err)
	}

	if len(data) < saltSize+nonceSize {
		logrus.Error("Decryption failed: data too short after decoding")
		return "", errors.New("invalid encrypted data format")
	}

	salt := data[:saltSize]
	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	key := GenerateKeyFromSecret([]byte(secretKey), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		logrus.Errorf("Decryption failed: failed to create cipher block: %v", err)
		return "", errors.New("failed to create cipher block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logrus.Errorf("Decryption failed: failed to create GCM: %v", err)
		return "", errors.New("failed to create GCM")
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logrus.Errorf("Decryption failed: authentication failed or data corrupted: %v", err)
		return "", errors.New("authentication failed or data corrupted")
	}

	return string(plaintext), nil
}
