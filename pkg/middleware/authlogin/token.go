package authlogin

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"

	"github.com/google/uuid"
)

type Metadata struct {
	AppID       uuid.UUID
	ThirdUserID string
	Third       string
	UserID      uuid.UUID
	ClientIP    net.IP
	UserAgent   string
}

func MetadataFromContext(ctx context.Context) (*Metadata, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("fail get metadata")
	}

	clientIP := ""
	if forwards, ok := meta["x-forwarded-for"]; ok {
		if len(forwards) > 0 {
			clientIP = strings.Split(forwards[0], ",")[0]
		}
	}

	userAgent := ""
	if agents, ok := meta["grpcgateway-user-agent"]; ok {
		if len(agents) > 0 {
			userAgent = agents[0]
		}
	}

	return &Metadata{
		ClientIP:  net.ParseIP(clientIP),
		UserAgent: userAgent,
	}, nil
}

func (meta *Metadata) ToJWTClaims() jwt.MapClaims {
	claims := jwt.MapClaims{}

	claims["app_id"] = meta.AppID
	claims["user_id"] = meta.UserID
	claims["third_user_id"] = meta.ThirdUserID
	claims["third"] = meta.Third
	claims["client_ip"] = meta.ClientIP
	claims["user_agent"] = meta.UserAgent

	return claims
}

func createToken(meta *Metadata) (string, error) {
	tokenAccessSecret, ok := os.LookupEnv("LOGIN_TOKEN_ACCESS_SECRET") // set to common-secret
	if !ok {
		// TODO
	}
	if tokenAccessSecret == "" {
		return "", fmt.Errorf("invalid login token access secret")
	}

	claims := meta.ToJWTClaims()
	candidate := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := candidate.SignedString([]byte(tokenAccessSecret))
	if err != nil {
		return "", fmt.Errorf("fail sign jwt claims: %v", err)
	}

	return token, nil
}
