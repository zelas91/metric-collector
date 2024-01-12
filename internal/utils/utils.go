package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
)

var ErrInvalidKey = errors.New("invalid key")

func GenerateHash(body []byte, key string) (*string, error) {
	if key == "" {
		return nil, nil
	}
	k, err := base64.StdEncoding.DecodeString(key)

	if err != nil {
		return nil, fmt.Errorf("generate hash decode key err:%w", ErrInvalidKey)
	}
	h := hmac.New(sha256.New, k)
	_, err = h.Write(body)
	if err != nil {
		return nil, fmt.Errorf("generate hash err:%w", err)
	}
	hash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return &hash, nil
}

func GetSubnet(subnet string) (*net.IPNet, bool) {
	_, network, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, false
	}
	return network, true
}
