package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/bagashiz/Simple-Bank/token"
	"google.golang.org/grpc/metadata"
)

// authorization constants
const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

// authorizeUser extracts the access token from the context and verifies it to authorize the user.
func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type %s", authType)
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %v", err)
	}

	return payload, nil
}
