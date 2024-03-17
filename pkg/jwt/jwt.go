package jwt

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrExpiredToken  = errors.New("token is expired")
)

func NewTokenizer(salt string) *Tokenizer {
	return &Tokenizer{
		jwt: jwtauth.New("HS256", []byte(salt), nil),
	}
}

type Tokenizer struct {
	jwt *jwtauth.JWTAuth
}

// Метод для создания Middleware, аутентификации пользователя
//
// Предоставляет токен записанный в запросе одним из следующих методов:
// - Bearer <token> в заголовке Authorization
// - Query параметр token
// - Cookie с именем access_token
func (tokenizer *Tokenizer) Authentificator(action func(*http.Request, jwt.Token, error) *http.Request) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := tokenizer.findToken(r)
			token, err := tokenizer.parseClaims(tokenString)
			action(r, token, err)
			h.ServeHTTP(w, r)
		})
	}
}

// Метод для создания нового access токена
//
// метод добавляет дополнительные поля в карту в остальном
// идентичен NewToken
// для инвалидации токена спустя 10 минут после его создания
func (tokenizer *Tokenizer) NewAccessToken(values map[string]interface{}) (string, error) {
	values["exp"] = time.Now().Add(time.Minute * 10)
	return tokenizer.NewToken(values)
}

// Метод для создания нового токена
func (tokenizer *Tokenizer) NewToken(values map[string]interface{}) (tokenString string, err error) {
	_, tokenString, err = tokenizer.jwt.Encode(values)
	return
}

// Метод для проверки токена полученного от пользователя
func (tokenizer *Tokenizer) ParseToken(tokenString string) (jwt.Token, error) {
	return tokenizer.parseClaims(tokenString)
}

func (tokenizer *Tokenizer) findToken(r *http.Request) string {
	fn := []func(r *http.Request) string{
		func(r *http.Request) string {
			// Bearer <token>
			header := r.Header.Get("Authorization")
			if len(header) > 7 && header[:7] == "Bearer " {
				return header[7:]
			}
			return ""
		},
		func(r *http.Request) string {
			return r.URL.Query().Get("token")
		},
		func(r *http.Request) string {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				return ""
			}
			return cookie.Value
		},
	}

	for _, f := range fn {
		token := f(r)
		if token != "" {
			return token
		}
	}

	return ""
}

func (tokenizer *Tokenizer) parseClaims(tokenString string) (jwt.Token, error) {
	if tokenString == "" {
		return nil, ErrTokenNotFound
	}

	token, err := tokenizer.jwt.Decode(tokenString)
	if err != nil {
		return nil, err
	}

	exp := token.Expiration()
	if exp.Before(time.Now()) && !exp.IsZero() {
		return nil, ErrExpiredToken
	}

	return token, nil
}
