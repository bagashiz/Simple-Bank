package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	// header for grpc-gateway
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwardedForHeader        = "x-forwarded-for"
	// header for grpc
	userAgentHeader = "user-agent"
)

// Metadata contains information about the client user agent and IP address.
type Metadata struct {
	UserAgent string
	ClientIp  string
}

// extractMetadata extracts the metadata from the context.
func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if clientIps := md.Get(xForwardedForHeader); len(clientIps) > 0 {
			mtdt.ClientIp = clientIps[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIp = p.Addr.String()
	}

	return mtdt
}
