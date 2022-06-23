package auth

import (
	"context"
	dto "github.com/isd-sgcu/rnkm65-auth/src/app/dto/auth"
	model "github.com/isd-sgcu/rnkm65-auth/src/app/model/auth"
	"github.com/isd-sgcu/rnkm65-auth/src/app/utils"
	"github.com/isd-sgcu/rnkm65-auth/src/constant"
	"github.com/isd-sgcu/rnkm65-auth/src/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repo           IRepository
	chulaSSOClient IChulaSSOClient
	tokenService   ITokenService
	userService    IUserService
}

type IRepository interface {
	FindByRefreshToken(string, *model.Auth) error
	FindByUserID(string, *model.Auth) error
	Create(*model.Auth) error
	Update(string, *model.Auth) error
}

type IChulaSSOClient interface {
	VerifyTicket(string, *dto.ChulaSSOCredential) error
}

type IUserService interface {
	FindByStudentID(string) (*proto.User, error)
	Create(*proto.User) (*proto.User, error)
}

type ITokenService interface {
	CreateOrUpdateCredentials(*model.Auth) (*proto.Credential, error)
}

func NewService(
	repo IRepository,
	chulaSSOClient IChulaSSOClient,
	tokenService ITokenService,
	userService IUserService,
) *Service {
	return &Service{
		repo:           repo,
		chulaSSOClient: chulaSSOClient,
		tokenService:   tokenService,
		userService:    userService,
	}
}

func (s *Service) VerifyTicket(_ context.Context, req *proto.VerifyTicketRequest) (res *proto.VerifyTicketResponse, err error) {
	ssoData := dto.ChulaSSOCredential{}
	auth := model.Auth{}

	err = s.chulaSSOClient.VerifyTicket(req.Ticket, &ssoData)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	user, err := s.userService.FindByStudentID(ssoData.Ouid)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				year, err := utils.CalYearFromID(ssoData.Ouid)
				if err != nil {
					return nil, err
				}

				faculty, err := utils.GetFacultyFromID(ssoData.Ouid)
				if err != nil {
					return nil, err
				}

				in := &proto.User{
					StudentID: ssoData.Ouid,
					Year:      year,
					Faculty:   faculty.FacultyEN,
				}

				user, err = s.userService.Create(in)
				if err != nil {
					return nil, status.Error(codes.Unauthenticated, st.Message())
				}

				auth = model.Auth{
					Role:   constant.USER,
					UserID: user.Id,
				}

				err = s.repo.Create(&auth)
				if err != nil {
					return nil, status.Error(codes.Unavailable, st.Message())
				}

			default:
				return nil, status.Error(codes.Unavailable, st.Message())
			}
		} else {
			return nil, status.Error(codes.Unavailable, "Service is down")
		}
	} else {
		err := s.repo.FindByUserID(user.Id, &auth)
		if err != nil {
			return nil, status.Error(codes.NotFound, "not found user")
		}
	}

	credential, err := s.tokenService.CreateOrUpdateCredentials(&auth)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.VerifyTicketResponse{Credential: credential}, err
}

func (s *Service) Validate(context.Context, *proto.ValidateRequest) (res *proto.ValidateResponse, err error) {
	return
}

func (s *Service) RefreshToken(context.Context, *proto.RefreshTokenRequest) (res *proto.RefreshTokenResponse, err error) {
	return
}
