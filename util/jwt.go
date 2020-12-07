package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type CustomToken struct {
	Name string
	Id int
	jwt.StandardClaims
}

func JWTCreate(name string, id int) (string, error) {
	mySigningKey := []byte(viper.GetString("jwt.sign"))
	token := CustomToken{
		Name: name,
		Id: id,
	}
	token.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("jwt.ttl")) * time.Minute).Unix(),
		Issuer:    viper.GetString("jwt.issue"),
	}

	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
	return tokenClaim.SignedString(mySigningKey)
}


// 解析Token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
    // Claims := &CustomToken{}
    // token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
    //     return  []byte(viper.GetString("token.sign")), nil
    // })
	// return token, err
	
	jwtToken, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(viper.GetString("jwt.sign")), nil
	})
	if err != nil || jwtToken == nil {
		return nil, err
	}
	claim, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok && jwtToken.Valid {
		return claim, nil
	} else {
		return nil, nil
	}
}