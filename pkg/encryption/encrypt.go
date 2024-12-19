package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type EncConfig struct {
	Key    string
	XORKey int64
}

var (
	aesKey string
	xorKey = int64(123)
)

func SetEncConfig(config EncConfig) EncConfig {
	aesKey = config.Key
	xorKey = config.XORKey
	return config
}

func GetXorKey() int64 {
	return xorKey
}

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Encrypt(text string) string {
	if text == "" {
		return ""
	}
	block, err := aes.NewCipher(createHash(aesKey))
	if err != nil {
		panic(err)
	}
	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return encodeBase64(ciphertext)
}

func EncryptNumber(num float64) float64 {
	return float64(int64(num) ^ xorKey)
}

func DecryptNumber(encryptedNum float64) float64 {
	return float64(int64(encryptedNum) ^ xorKey)
}
