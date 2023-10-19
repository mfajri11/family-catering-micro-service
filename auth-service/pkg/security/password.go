package security

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s security) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.Wrap(err, "fail to hash password")
		return "", err // TODO: wrap error
	}

	return string(hashedPassword), nil
}

func (s security) CompareHashPassword(password, hashedPassword string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err // TODO: wrap error
	}

	return nil
}
