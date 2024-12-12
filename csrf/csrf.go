package csrf

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	CSRF_SECRET_LENGTH = 32
)

func GenerateCSRFSecret() (string, error) {
	secret := make([]byte, CSRF_SECRET_LENGTH)

	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}

	encodedSecret := hex.EncodeToString(secret)

	return encodedSecret, nil
}

func EncryptToken(unixTime int64, key []byte) (string, error) {
	var encryptedString string

	unixTimeBytes := []byte(strconv.FormatInt(unixTime, 10))

	encodedKey := base64.StdEncoding.EncodeToString(key)

	joinedData := strings.Join([]string{string(encodedKey), string(unixTimeBytes)}, ":")

	paddedData := pad([]byte(joinedData), aes.BlockSize)

	block, err := aes.NewCipher(key)
	if err != nil {
		return encryptedString, err
	}

	cipherText := make([]byte, aes.BlockSize+len(paddedData))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return encryptedString, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(cipherText[aes.BlockSize:], paddedData)

	encryptedString = base64.StdEncoding.EncodeToString(cipherText)

	return encryptedString, nil
}

func DecryptToken(encryptedStr string, userToken []byte) ([]byte, int64, error) {
	var userKey []byte
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return userKey, 0, err
	}

	block, err := aes.NewCipher(userToken)
	if err != nil {
		return userKey, 0, err
	}

	if len(encryptedData) < aes.BlockSize {
		return userKey, 0, errors.New("encrypted data is too short")
	}

	iv := encryptedData[:aes.BlockSize]
	cipherText := encryptedData[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(cipherText, cipherText)

	unpaddedData, err := unpad(cipherText)
	if err != nil {
		return userKey, 0, err
	}

	parts := strings.Split(string(unpaddedData), ":")
	if len(parts) != 2 {
		return userKey, 0, errors.New("invalid decrypted data format")
	}

	userKey, err = base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return userKey, 0, err
	}

	unixTime, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return userKey, 0, err
	}
	return userKey, unixTime, nil
}

func pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("unpad error: input data is empty")
	}
	unpadding := int(data[length-1])
	if unpadding > length || unpadding == 0 {
		return nil, errors.New("unpad error: unpadding value is invalid")
	}
	return data[:(length - unpadding)], nil
}

func ValidateCSRFToken(isUsed bool, csrfToken string, userToken []byte) error {
	decryptedKey, decryptedUnixTime, err := DecryptToken(csrfToken, userToken)
	if err != nil {
		fmt.Printf("ERROR DECRYPTING TOKEN: %+v\n", err)
		return err
	}

	if !bytes.Equal(decryptedKey, userToken) {
		return errors.New("tokens don't match")
	}

	// Unix time validation
	if time.Now().Unix() > decryptedUnixTime {
		return errors.New("token expired")
	}

	// Check if used
	if isUsed {
		return errors.New("token already used")
	}

	return nil
}
