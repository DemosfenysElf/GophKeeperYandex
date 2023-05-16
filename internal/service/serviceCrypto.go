package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"

	"github.com/pkg/errors"
)

func EnCrypt(message []byte, key string) ([]byte, error) {
	sum := md5.Sum([]byte(key))
	key16 := sum[:]
	block, err := aes.NewCipher(key16)

	if err != nil {
		return nil, err
	}

	b := message
	b = PKCS5Padding(b, aes.BlockSize)
	encMessage := make([]byte, len(b))
	iv := key16[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encMessage, b)

	return encMessage, nil
}

func DeCrypt(encMessage []byte, key string) ([]byte, error) {
	sum := md5.Sum([]byte(key))
	key16 := sum[:]
	iv := key16[:aes.BlockSize]
	block, err := aes.NewCipher(key16)

	if err != nil {
		return nil, err
	}

	if len(encMessage) < aes.BlockSize {
		return nil, errors.New("encMessage слишком короткий")
	}

	decrypted := make([]byte, len(encMessage))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, encMessage)

	return PKCS5UnPadding(decrypted), nil
}

func PKCS5Padding(cipher []byte, blockSize int) []byte {
	padding := blockSize - len(cipher)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipher, padText...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:(length - unPadding)]
}
