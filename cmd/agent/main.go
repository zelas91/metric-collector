package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/zelas91/metric-collector/internal/logger"
	"os"
)

var privPEMData = []byte(`
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)

func main() {
	//conf := NewConfig()
	//ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//_ = cancel
	//var pubKey []byte
	//agent.Run(ctx, conf.PollInterval, conf.ReportInterval, conf.BaseURL, conf.Key, *conf.RateLimit, pubKey)
	//log.Info("start agent")
	//<-ctx.Done()
	//stop()
	readPublicKey("key.pub")
}
func readPublicKey(path string) []byte {
	readKey, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("READ FILE public key = %s , err=%v", path, err)
	}
	block, _ := pem.Decode(readKey)
	pub, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	var encryptedData []byte
	originalData := []byte("Привет, мир!")
	if key, ok := pub.PublicKey.(*rsa.PublicKey); ok {
		encryptedData, err = rsa.EncryptPKCS1v15(rand.Reader, key, originalData)
		fmt.Printf("Got a %T, with remaining data:\n\n\n", pub.PublicKey)
	}
	//fmt.Println(string(encryptedData))
	b, _ := pem.Decode(privPEMData)
	priv, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, priv, encryptedData)
	if err != nil {
		fmt.Println("Не удалось расшифровать данные:", err)
	}
	fmt.Println(string(decryptedData))
	return nil

}
func stop() {
	log.Info("stop agent")
	logger.Shutdown()
	os.Exit(0)
}
