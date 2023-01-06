package interceptors

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// BasicAuthInterceptor creates a gRPC interceptor for the Basic Auth
// we need some way of authorizing to this API, as we do not want to deploy
// open API without any restriction
// TODO: restrict to specific group in LDAP
func BasicAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	auth, err := extractHeader(ctx, "authorization")
	if err != nil {
		return ctx, err
	}

	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return ctx, status.Error(codes.Unauthenticated, `missing "Basic " prefix in "Authorization" header`)
	}

	decodedAuthBytes, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, `invalid base64 in header`)
	}

	decodedAuthStr := string(decodedAuthBytes)
	s := strings.IndexByte(decodedAuthStr, ':')
	if s < 0 {
		return ctx, status.Error(codes.Unauthenticated, `invalid basic auth format`)
	}

	username, password := decodedAuthStr[:s], decodedAuthStr[s+1:]
	if username != "user" || password != "psswd" {
		return ctx, status.Error(codes.Unauthenticated, "invalid user or password")
	} else {
		fmt.Printf("Authenticated the user name %s", username)
	}

	return handler(ctx, req)
}

func extractHeader(ctx context.Context, header string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Unauthenticated request.")
	}

	authHeaders, ok := md[header]
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Unauthenticated request.")
	}

	if len(authHeaders) != 1 {
		return "", status.Error(codes.Unauthenticated, "more than 1 header in request")
	}

	return authHeaders[0], nil
}
