package auth

import "golang.org/x/crypto/bcrypt"

// generate a hash for the user password
// @params: p(string) user password
func HashPassword(p string) (string, error) {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(hashByte), err
}

// compare the hash with the user input password
// @params: h(string) hashed password, p(string) user input password
func CompareHashPassword(h, p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	return err == nil
}
