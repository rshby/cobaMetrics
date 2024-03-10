package helper

import "golang.org/x/crypto/bcrypt"

type IHelperPassword interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type HelperPassword struct {
}

func NewHelperPassword() IHelperPassword {
	return &HelperPassword{}
}

func (h *HelperPassword) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (h *HelperPassword) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
