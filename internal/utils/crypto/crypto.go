package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/zelas91/metric-collector/internal/logger"
	"os"
)

var log = logger.New()

func LoadPublicKey(path string) *rsa.PublicKey {
	block, err := loadBlock(path)
	if err != nil {
		log.Errorf("error read file public key = %s  err: %v", path, err)
		return nil
	}
	pub, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Errorf("parse certificate err=%v", err)
		return nil
	}

	if key, ok := pub.PublicKey.(*rsa.PublicKey); ok {
		return key
	}
	return nil
}

func LoadPrivateKey(path string) *rsa.PrivateKey {
	block, err := loadBlock(path)
	if err != nil {
		log.Errorf("error read file private key = %s  err: %v", path, err)
		return nil
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Errorf("parse private key err=%v", err)
		return nil
	}
	pvKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		log.Errorf("cast key err (rsq.PrivateKey)")
		return nil
	}
	return pvKey
}

func loadBlock(path string) (*pem.Block, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Errorf(`read file private key ="%s" , err=%v`, path, err)
		return nil, err
	}
	block, _ := pem.Decode(data)
	return block, nil

}
func SplitData(data []byte, size int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(data)/size+1)
	for len(data) >= size {
		chunk, data = data[:size], data[size:]
		chunks = append(chunks, chunk)
	}
	if len(data) > 0 {
		chunks = append(chunks, data[:])
	}
	return chunks
}
