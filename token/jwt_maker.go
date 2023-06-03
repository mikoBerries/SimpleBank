package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

//JWTMaker
type JWTMaker struct {
	secretKey string
}

// CreateToken implements Maker.
func (m *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	//make token with header method and payload/claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	//finish jwt with signed it with secret key
	jwt, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", err
	}

	return jwt, nil
}

// VerifyToken implements Maker.
func (m *JWTMaker) VerifyToken(token string) (*Payload, error) {
	var keyfunc jwt.Keyfunc
	keyfunc = func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() { //checking vurnability hashing method
			return nil, fmt.Errorf("invalid signing algoritmn")
		}
		return []byte(m.secretKey), nil
		//get kid value from header
		// kid, ok := t.Header["kid"].(string)
		// if !ok {
		// 	return nil, fmt.Errorf("invalid key id")
		// }
		//check kid from database
		//exist? / expired ? / ETC
		// k, ok := keys[kid]
		// if !ok {
		// 	return nil, fmt.Errorf("invalid key id")
		// }
	}
	//parse to get *jwt.Token
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyfunc)

	//check jwt standart claim valid func
	if !jwtToken.Valid { // not valid signature with provided secret key
		if errors.Is(err, jwt.ErrTokenMalformed) {
			err = fmt.Errorf("that's not even a token")
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			// Invalid signature
			err = fmt.Errorf("invalid signature")
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			// Token is either expired or not active yet
			err = fmt.Errorf("token is expired")
		} else {
			err = fmt.Errorf("couldn't handle this token: %w", err)
		}

		return nil, err
	}
	return jwtToken.Claims.(*Payload), nil
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, nil
	}
	return &JWTMaker{secretKey}, nil
}
