package vos

type Token struct {
	AccessToken   string
	TokenType     string
	ExpiresAtUnix int64
}