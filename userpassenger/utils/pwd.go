package utils

import (
	"crypto/sha512"
	"fmt"
	"strings"

	"github.com/anaskhan96/go-password-encoder"
)

var (
	passMethod = "pbkdf2-sha512"
)

func Encrypt(pwd string) string {
	options := &password.Options{SaltLen: 10, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(pwd, options)
	pwd = fmt.Sprintf("$%s$%s$%s", passMethod, salt, encodedPwd)
	return pwd
}

func GetEncoded(encryptedPwd string) (method string, salt string, encodedPwd string) {
	var s = strings.Split(encryptedPwd, "$")[1:]
	method = s[0]
	salt = s[1]
	encodedPwd = s[2]
	return
}

func Verify(rawPwd, salt, encodedPwd string) bool {
	options := &password.Options{SaltLen: 10, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	return password.Verify(rawPwd, salt, encodedPwd, options)
}
