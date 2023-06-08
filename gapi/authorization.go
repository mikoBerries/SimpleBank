package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/MikoBerries/SimpleBank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "Authorization"
	authorizationBearer = "bearer"
)

func (server *server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	//extract meta data from incoming context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	value := md.Get(authorizationHeader)
	if len(value) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}
	//get header param (authHeader) value
	authHeader := value[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	//check (auth-type) -> authorizationBearer
	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}
	//verify token (Paseto)
	accessToken := fields[1]
	payload, err := server.token.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	return payload, nil
}
