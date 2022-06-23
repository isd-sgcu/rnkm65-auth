package token

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	model "github.com/isd-sgcu/rnkm65-auth/src/app/model/auth"
	"github.com/isd-sgcu/rnkm65-auth/src/proto"
)

type Service struct {
	jwtService IJwtService
}

type IJwtService interface {
	SignAuth(*model.Auth) (string, error)
	VerifyAuth(string) (*jwt.Token, error)
}

func (s *Service) CreateOrUpdateCredentials(*model.Auth) (*proto.Credential, error) {
	return nil, nil
}

func createRefreshToken() string {
	return uuid.New().String()
}
