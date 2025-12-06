package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

func GenerateKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

func Encrypt(publicKey *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)
}

func Decrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

func ParseRSAPublicKeyFromPEM(pubKeyStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubKeyStr))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not a valid RSA public key")
	}

	return rsaPub, nil
}

func NormalizePublicKey(pubKey string) string {
	pubKey = strings.ReplaceAll(pubKey, "-----BEGIN PUBLIC KEY-----", "")
	pubKey = strings.ReplaceAll(pubKey, "-----END PUBLIC KEY-----", "")
	pubKey = strings.ReplaceAll(pubKey, "\n", "")
	pubKey = strings.TrimSpace(pubKey)

	var buf strings.Builder
	buf.WriteString("-----BEGIN PUBLIC KEY-----\n")
	for len(pubKey) > 64 {
		buf.WriteString(pubKey[:64])
		buf.WriteString("\n")
		pubKey = pubKey[64:]
	}
	if len(pubKey) > 0 {
		buf.WriteString(pubKey)
		buf.WriteString("\n")
	}
	buf.WriteString("-----END PUBLIC KEY-----\n")
	return buf.String()
}
