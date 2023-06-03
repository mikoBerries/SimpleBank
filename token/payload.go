package token

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

//Payload contains the payload data of token (Claims in jwt)
type Payload struct {
	jwt.RegisteredClaims
	UserName string
	Role     string
}

//Payload Builder
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	// fmt.Println(uuid.)
	// jwt.RegisteredClaims{
	// 	Issuer:    "localhost:8080",
	// 	Audience:  jwt.ClaimStrings{},
	// 	Subject:   "session",
	// 	ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
	// 	NotBefore: jwt.NewNumericDate(time.Now()),
	// 	IssuedAt:  jwt.NewNumericDate(time.Now()),
	// 	ID:        "1",
	// },
	payload := &Payload{
		RegisteredClaims: jwt.RegisteredClaims{
			// Audience:  jwt.ClaimStrings{},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			ID:        uuid.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "localhost:8080",
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "session",
		},
		UserName: username,
		Role:     "user",
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt.Time) {
		return ErrExpiredToken
	}
	return nil
}
