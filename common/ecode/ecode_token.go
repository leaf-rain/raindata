package ecode

import "errors"

var (
	ERR_TokenExpired     = errors.New("Token is expired")
	ERR_TokenNotValidYet = errors.New("Token not active yet")
	ERR_TokenMalformed   = errors.New("That's not even a token")
	ERR_TokenInvalid     = errors.New("Couldn't handle this token:")
)
