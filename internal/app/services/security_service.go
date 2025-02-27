package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type DataEncryptorService struct {
	encryptionKey string
}

func NewDataEncryptorService() DataEncryptorService {
	return DataEncryptorService{"d8b74a074f6d77882d4e44a2f128ae22d2ac044277d45298b0ea81af4c9219d1"}
}

// MaskData encrypts the input data using AES encryption
func (encryptor DataEncryptorService) MaskData(data []byte) ([]byte, error) {
	key, _ := hex.DecodeString(encryptor.encryptionKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""), err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return []byte(""), err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte(""), err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (encryptor DataEncryptorService) UnmaskData(ciphertextBytes []byte) ([]byte, error) {

	key, _ := hex.DecodeString(encryptor.encryptionKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""), err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return []byte(""), err
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return []byte(""), err
	}

	return plaintext, nil
}
