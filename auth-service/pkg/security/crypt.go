package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"hash"

	"github.com/mfajri11/family-catering-micro-service/auth-service/config"
)

var (
	hashF      hash.Hash       = sha256.New()
	privateKey *rsa.PrivateKey = config.Cfg.App.Security.PrivKey()
	publicKey  *rsa.PublicKey  = config.Cfg.App.Security.PubKey()
)

func (s security) Decrypt(cipher string) (string, error) {
	// assume data is encoded with base64 encoding
	decodedCipher, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		return "", err // TODO: wrap error
	}
	plaintextBytes, err := rsa.DecryptOAEP(hashF, rand.Reader, privateKey, decodedCipher, nil)
	if err != nil {
		return "", err // TODO: wrap error
	}

	return string(plaintextBytes), nil
}

func encrypt(plaintext string) ([]byte, error) {
	cipher, err := rsa.EncryptOAEP(hashF, rand.Reader, publicKey, []byte(plaintext), nil)
	if err != nil {
		return nil, err
	}

	return cipher, nil
}

func (s security) EncryptWithURLEncode(plaintext string) (string, error) {
	cipher, err := encrypt(plaintext)
	if err != nil {
		return "", err // TODO: wrap error
	}

	return base64.URLEncoding.EncodeToString(cipher), nil
}
