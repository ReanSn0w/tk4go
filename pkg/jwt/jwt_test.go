package jwt_test

import (
	"testing"
	"time"

	"github.com/ReanSn0w/tk4go/pkg/jwt"
)

func Test_Tokenizer(t *testing.T) {
	tokenizer := jwt.NewTokenizer("some_salt")
	tokenString, err := tokenizer.NewAccessToken(map[string]interface{}{
		"key": "value",
	})

	if err != nil {
		t.Log("new token failed with error", err)
		t.FailNow()
	}

	token, err := tokenizer.ParseToken(tokenString)
	if err != nil {
		t.Log("token parse failed with error", err)
		t.FailNow()
	}

	val, ok := token.Get("key")
	if !ok || val != "value" {
		t.Log("token not contain tested value")
		t.Fail()
	}
}

func Test_TokenizerExpiredToken(t *testing.T) {
	tokenizer := jwt.NewTokenizer("some_salt")
	tokenString, err := tokenizer.NewToken(map[string]interface{}{
		"key": "value",
		"exp": time.Now().Add(-1 * time.Second),
	})

	if err != nil {
		t.Log("new token failed with error", err)
		t.FailNow()
	}

	_, err = tokenizer.ParseToken(tokenString)
	if err != jwt.ErrExpiredToken {
		t.Log("token parse failed with error", err)
		t.FailNow()
	}
}

func Test_TokenizerEmptyToken(t *testing.T) {
	tokenizer := jwt.NewTokenizer("some_salt")
	_, err := tokenizer.ParseToken("")
	if err != jwt.ErrTokenNotFound {
		t.Log("token parse failed with error", err)
		t.FailNow()
	}
}

func Test_TokenizerFailedToken(t *testing.T) {
	tokenizer := jwt.NewTokenizer("some_salt")
	_, err := tokenizer.ParseToken("unvalid_token_string")
	if err == nil {
		t.Log("token parse failed with error", err)
		t.FailNow()
	}
}
