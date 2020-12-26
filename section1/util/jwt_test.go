package util_test

import (
	"testing"
	"go-kit-demo/util"
)

func TestJWTCreate(t *testing.T) {
	jwtToken, err := util.JWTCreate("sai",2)
	if err != nil {
		t.Error(err)
	}
	t.Log(jwtToken)
}

func TestParseToken(t *testing.T) {
	jwtToken, _ := util.JWTCreate("sai",2)
	jwtInfo, err := util.ParseToken(jwtToken)
	if err != nil {
		t.Error(err)
	}
	t.Log(jwtInfo)
}