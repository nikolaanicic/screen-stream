package util

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(data []byte) string {
	hash, _ := bcrypt.GenerateFromPassword(data,bcrypt.DefaultCost)
	return string(hash)
}

func CompareHash(hash string, plain string) error{
	return bcrypt.CompareHashAndPassword([]byte(hash),[]byte(plain))
}