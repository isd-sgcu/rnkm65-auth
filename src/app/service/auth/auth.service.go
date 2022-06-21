package auth

import (
	"context"
	model "github.com/isd-sgcu/rnkm65-auth/src/app/model/auth"
	"github.com/isd-sgcu/rnkm65-auth/src/proto"
)

type Service struct {
	repo IRepository
}

type IRepository interface {
	FindByStudentID(string, *model.Auth) error
	FindByRefreshToken(string, *model.Auth) error
	Create(*model.Auth) error
}

func NewService(repo IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) VerifyTicket(context.Context, *proto.VerifyTicketRequest) (res *proto.VerifyTicketResponse, err error) {
	return
}

func (s *Service) Validate(context.Context, *proto.ValidateRequest) (res *proto.ValidateResponse, err error) {
	return
}

func (s *Service) RefreshToken(context.Context, *proto.RefreshTokenRequest) (res *proto.RefreshTokenResponse, err error) {
	return
}
