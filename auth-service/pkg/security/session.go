package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// func generateRandomString(n int) (string, error) {
// 	s, err := generateRandomBytes(n)
// 	if err != nil {
// 		return "", err
// 	}

// 	return hex.EncodeToString(s), nil
// }

func generateRandomBytes(n int) ([]byte, error) {
	s := make([]byte, n)
	_, err := rand.Read(s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s security) GenerateSID() (string, error) {
	sid, err := generateRandomBytes(16)
	fmt.Println(sid)
	if err != nil {
		return "", err // TODO: wrap error
	}
	encodedSID := base64.StdEncoding.EncodeToString(sid)

	return encodedSID, nil

}
