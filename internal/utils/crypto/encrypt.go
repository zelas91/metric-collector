package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
)

func Encrypt(key *rsa.PublicKey, data []byte) ([]byte, error) {
	if key == nil {
		return data, nil
	}
	splitData := SplitData(data, key.Size()-11)
	buf := bytes.NewBuffer([]byte{})
	for _, val := range splitData {
		encodedData, err := rsa.EncryptPKCS1v15(rand.Reader, key, val)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(encodedData)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
