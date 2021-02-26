package util

import (
	"github.com/pyihe/secret"
)

var DefaultSalt = "1234567812345678"

func EncodePassword(password string, salt string) string {
	c := secret.NewCipher()
	var request = &secret.SymRequest{
		Key:         []byte(salt),
		Type:        secret.SymTypeAES,
		ModeType:    secret.BlockModeECB,
		PaddingType: secret.PaddingTypeZeros,
	}
	request.PlainData = password
	cipherString, _ := c.SymEncryptToString(request)
	return cipherString
}

func DecodePassword(encodePassword string, salt string) string {
	c := secret.NewCipher()
	var request = &secret.SymRequest{
		Key:         []byte(salt),
		Type:        secret.SymTypeAES,
		ModeType:    secret.BlockModeECB,
		PaddingType: secret.PaddingTypeZeros,
	}
	request.CipherData = encodePassword
	plainText, _ := c.SymDecrypt(request)
	return string(plainText)
}
