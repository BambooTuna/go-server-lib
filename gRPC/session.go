package gRPC

import (
	"context"
	session "github.com/BambooTuna/go-server-lib/session_v2"
	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Session struct {
	Dao      session.SessionStorageDao
	Settings session.SessionSettings
}

func (s Session) SetSessionData(ctx context.Context, data string) (string, error) {
	id, err := s.Settings.IDGenerator()
	if err != nil {
		return "", status.Errorf(codes.Internal, "generate session token id failed")
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &jwt.StandardClaims{Id: id})
	signedToken, err := token.SignedString([]byte(s.Settings.Secret))
	if err != nil {
		return "", status.Errorf(codes.Internal, "generate jwt token failed")
	}

	if _, err := s.Dao.Store(id, data, s.Settings.ExpirationDate); err != nil {
		return "", status.Errorf(codes.Internal, "session store failed")
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(s.Settings.SetAuthHeaderName, signedToken)); err != nil {
		return "", status.Errorf(codes.Internal, "set session token failed")
	}
	return signedToken, nil
}

func (s Session) RequiredSession(ctx context.Context) (string, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "please set session token")
	}

	id, err := s.ParseClaimsId(token)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "invalid token: 4000")
	}

	storedData, err := s.Dao.Find(id)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "invalid token: 4004")
	}

	s.Dao.Refresh(id, s.Settings.ExpirationDate)
	return storedData, nil
}

func (s Session) InvalidateSession(ctx context.Context) error {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "please set session token")
	}

	id, err := s.ParseClaimsId(token)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	if _, err := s.Dao.Find(id); err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	if _, err := s.Dao.Remove(id); err != nil {
		return status.Errorf(codes.Internal, "remove session token failed")
	}

	return nil
}

func (s Session) ParseClaimsId(token string) (string, error) {
	var claims jwt.StandardClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Settings.Secret), nil
	})
	if err != nil {
		return "", err
	}
	return claims.Id, nil
}
