package token

import "time"

// Maker interface for managin token
type Maker interface {
	//CreateToken crete new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	//VerifyToken check token valid
	VerifyToken(token string) (*Payload, error)
}

var _ Maker = (*JWTMaker)(nil)
var _ Maker = (*PasetoMaker)(nil)
