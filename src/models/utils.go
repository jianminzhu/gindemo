package models

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"gin/src/util"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
	"io"
	"time"
)

// JWT : HEADER PAYLOAD SIGNATURE
const (
	SecretKEY              string = "server"
	DEFAULT_EXPIRE_SECONDS int    = 600 // default expired 10 minutes
	PasswordHashBytes             = 16
)

// This struct is the payload
type MyCustomClaims struct {
	UserID         int64  `json:"userID"`
	Mobile         string `json:"mobile"`
	InvitationCode string `json:"invitationCode"`
	jwt.StandardClaims
}

// This struct is the parsing of token payload
type JwtPayload struct {
	Username       string `json:"username"`
	UserID         int64  `json:"userID"`
	InvitationCode string `json:"invitationCode"`
	IssuedAt       int64  `json:"iat"`
	ExpiresAt      int64  `json:"exp"`
}

func GenerateEncToken(loginInfo *LoginRequest, userID int64, expiredSeconds int) (tokenString string, err error) {
	token, err := genToken(loginInfo, userID, expiredSeconds)
	return util.EncodePassword(token, util.DefaultSalt), err

}

// generate token
func genToken(loginInfo *LoginRequest, userID int64, expiredSeconds int) (tokenString string, err error) {
	if expiredSeconds == 0 {
		expiredSeconds = DEFAULT_EXPIRE_SECONDS
	}

	// Create the Claims
	mySigningKey := []byte(SecretKEY)
	expireAt := time.Now().Add(time.Second * time.Duration(expiredSeconds)).Unix()
	logs.Info("Token will be expired at ", time.Unix(expireAt, 0))

	user := *loginInfo
	claims := MyCustomClaims{
		userID,
		user.Mobile,
		util.InviteCodeDefault.IdToCode(uint64(userID)),
		jwt.StandardClaims{
			Issuer:    user.Username,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireAt,
		},
	}

	// Create the token using your claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Signs the token with a secret
	tokenStr, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("error: failed to generate token")
	}

	return tokenStr, nil
}

//先解密，后验证token
func ValidateEncToken(encTokenString string) (*JwtPayload, error) {
	tokenString := util.DecodePassword(encTokenString, util.DefaultSalt)
	return validateToken(tokenString)
}

/**验证token */
func validateToken(tokenString string) (*JwtPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKEY), nil
		})

	claims, ok := token.Claims.(*MyCustomClaims)
	if ok && token.Valid {
		logs.Info("%v %v", claims.UserID, claims.StandardClaims.ExpiresAt)
		logs.Info("Token will be expired at ", time.Unix(claims.StandardClaims.ExpiresAt, 0))

		return &JwtPayload{
			Username:  claims.StandardClaims.Issuer,
			UserID:    claims.UserID,
			IssuedAt:  claims.StandardClaims.IssuedAt,
			ExpiresAt: claims.StandardClaims.ExpiresAt,
		}, nil
	} else {
		logs.Info(err.Error())
		return nil, errors.New("登录已失效")
	}
}

//update token
func RefreshToken(tokenString string) (newTokenString string, err error) {
	// get previous token
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKEY), nil
		})

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok || !token.Valid {
		return "", err
	}

	mySigningKey := []byte(SecretKEY)
	expireAt := time.Now().Add(time.Second * time.Duration(DEFAULT_EXPIRE_SECONDS)).Unix() //new expired
	newClaims := MyCustomClaims{
		claims.UserID,
		claims.Mobile,
		util.InviteCodeDefault.IdToCode(uint64(claims.UserID)),
		jwt.StandardClaims{
			Issuer:    claims.StandardClaims.Issuer, //name of token issue
			IssuedAt:  time.Now().Unix(),            //time of token issue
			ExpiresAt: expireAt,
		},
	}
	// generate new token with new claims
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	// sign the token with a secret
	tokenStr, err := newToken.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("error: failed to generate new fresh json web token")
	}

	return tokenStr, nil
}

// generate salt
func GenerateSalt() (salt string, err error) {
	buf := make([]byte, PasswordHashBytes)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", errors.New("error: failed to generate user's salt")
	}

	return fmt.Sprintf("%x", buf), nil
}

// generate password hash
func GeneratePassHash(password string, salt string) (hash string, err error) {
	h, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, PasswordHashBytes)
	if err != nil {
		return "", errors.New("error: failed to generate password hash")
	}

	return fmt.Sprintf("%x", h), nil
}

func ToJson(tempMap interface{}) string {
	data, _ := json.Marshal(tempMap)
	return string(data)
}

func ToFormatJson(tempMap interface{}) string {
	data, _ := json.MarshalIndent(tempMap, "", "\t")
	return string(data)
}
