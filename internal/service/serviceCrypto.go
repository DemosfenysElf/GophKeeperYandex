package service

import "crypto/aes"

var key = []byte("Increment #9 key")

// CryptoToken шифрование данных
func CryptoToken(token []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cryptoString := make([]byte, len(token))
	aesBlock.Encrypt(cryptoString, token)
	return cryptoString, nil
}

// DeCryptoToken расшифрование данных
func DeCryptoToken(token []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	deCryptoString := make([]byte, len(token))
	aesBlock.Decrypt(deCryptoString, token)
	return deCryptoString, nil
}
