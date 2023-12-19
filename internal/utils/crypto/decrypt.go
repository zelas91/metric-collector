package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
)

func Decrypt(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	splitData := SplitData(data, key.Size())
	var buf bytes.Buffer

	for _, val := range splitData {
		decryptChunk, err := rsa.DecryptPKCS1v15(rand.Reader, key, val)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(decryptChunk)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
