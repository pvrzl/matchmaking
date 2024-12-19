package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io/ioutil"

	logHelper "app/pkg/log"
)

func createHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

func Decrypt(text string) string {
	if text == "" {
		return ""
	}
	block, err := aes.NewCipher(createHash(aesKey))
	if err != nil {
		logHelper.Logger.Debugf("[pkg][DecryptFile] failed during decrypting %s: %s\n", text, err)
		return text
	}
	ciphertext := decodeBase64(text)
	cfb := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext)
}

func decrypt(data []byte, passphrase string) []byte {
	key := createHash(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func DecryptFile(filename string, passphrase string) []byte {
	pp := []byte{108, 111, 99, 97, 108, 112, 114, 105, 100, 101}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logHelper.Logger.Errorln("[pkg][DecryptFile] error reading file:", err)
	}
	return decrypt(data, string(pp))
}
