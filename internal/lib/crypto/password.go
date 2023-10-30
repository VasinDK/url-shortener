package crypto

import "golang.org/x/crypto/bcrypt"

func HashPass(pass string) (string, error) {
	byte, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	return string(byte), err
}

func CheckPassHash(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}