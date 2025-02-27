package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskData_Success(t *testing.T) {
	// Given
	dataEncryptor := NewDataEncryptorService()
	plaintext := []byte("SensitiveData")

	// When
	encryptedData, err := dataEncryptor.MaskData(plaintext)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedData)
	assert.NotEqual(t, plaintext, encryptedData) // Ensure data is encrypted
}

func TestUnmaskData_Success(t *testing.T) {
	// Given
	dataEncryptor := NewDataEncryptorService()
	plaintext := []byte("SensitiveData")

	// Encrypt the data
	encryptedData, err := dataEncryptor.MaskData(plaintext)
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedData)

	// When
	decryptedData, err := dataEncryptor.UnmaskData(encryptedData)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decryptedData) // Ensure decrypted text matches the original
}

func TestMaskData_InvalidKey(t *testing.T) {
	// Given
	invalidDataEncryptor := DataEncryptorService{
		encryptionKey: "invalidkey", // Key that is not a valid AES key length
	}
	plaintext := []byte("SensitiveData")

	// When
	encryptedData, err := invalidDataEncryptor.MaskData(plaintext)

	// Then
	assert.Error(t, err)
	assert.Empty(t, encryptedData) // Encryption should fail
}

func TestUnmaskData_InvalidKey(t *testing.T) {
	// Given
	dataEncryptor := NewDataEncryptorService()
	plaintext := []byte("SensitiveData")

	// Encrypt data with valid DataEncryptorService
	encryptedData, err := dataEncryptor.MaskData(plaintext)
	assert.NoError(t, err)

	// Attempt decryption with invalid key
	invalidDataEncryptor := DataEncryptorService{
		encryptionKey: "invalidkey", // Key that differs from the one used during encryption
	}

	// When
	decryptedData, err := invalidDataEncryptor.UnmaskData(encryptedData)

	// Then
	assert.Error(t, err)
	assert.Empty(t, decryptedData) // Decryption should fail
}

func TestUnmaskData_InvalidCiphertext(t *testing.T) {
	// Given
	dataEncryptor := NewDataEncryptorService()
	invalidCiphertext := []byte("InvalidCiphertextData") // Corrupted or invalid ciphertext

	// When
	decryptedData, err := dataEncryptor.UnmaskData(invalidCiphertext)

	// Then
	assert.Error(t, err)
	assert.Empty(t, decryptedData) // Decryption should fail
}
