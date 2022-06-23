package auth

import (
	"context"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	dto "github.com/isd-sgcu/rnkm65-auth/src/app/dto/auth"
	"github.com/isd-sgcu/rnkm65-auth/src/app/model"
	"github.com/isd-sgcu/rnkm65-auth/src/app/model/auth"
	"github.com/isd-sgcu/rnkm65-auth/src/constant"
	mock "github.com/isd-sgcu/rnkm65-auth/src/mocks/auth"
	"github.com/isd-sgcu/rnkm65-auth/src/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"testing"
	"time"
)

type AuthServiceTest struct {
	suite.Suite
	Auth            *auth.Auth
	UserDto         *proto.User
	Credential      *proto.Credential
	UnauthorizedErr error
	NotFoundErr     error
	ServiceDownErr  error
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTest))
}

func (t *AuthServiceTest) SetupTest() {
	t.Auth = &auth.Auth{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		UserID:       faker.UUIDDigit(),
		Role:         constant.USER,
		RefreshToken: faker.Word(),
	}

	t.UserDto = &proto.User{
		Id:                    t.Auth.UserID,
		Firstname:             faker.FirstName(),
		Lastname:              faker.LastName(),
		Nickname:              faker.Name(),
		StudentID:             "63xxxxxx21",
		Faculty:               "Faculty of Engineering",
		Year:                  "3",
		Phone:                 faker.Phonenumber(),
		LineID:                faker.Word(),
		Email:                 faker.Email(),
		AllergyFood:           faker.Word(),
		FoodRestriction:       faker.Word(),
		AllergyMedicine:       faker.Word(),
		Disease:               faker.Word(),
		VaccineCertificateUrl: faker.URL(),
		ImageUrl:              faker.URL(),
	}

	t.Credential = &proto.Credential{
		AccessToken:  faker.Word(),
		RefreshToken: t.Auth.RefreshToken,
		ExpiresIn:    3600,
	}

	t.UnauthorizedErr = errors.New("unauthorized")
	t.NotFoundErr = errors.New("not found user")
	t.ServiceDownErr = errors.New("service is down")
}

func (t *AuthServiceTest) TestVerifyTicketSuccessFirstTimeLogin() {
	want := &proto.VerifyTicketResponse{
		Credential: t.Credential,
	}

	ticket := faker.Word()
	chulaSSORes := &dto.ChulaSSOCredential{
		UID:         faker.Word(),
		Username:    faker.Username(),
		Gecos:       faker.Username(),
		Email:       faker.Email(),
		Disable:     false,
		Roles:       []string{"student"},
		Firstname:   faker.FirstName(),
		Lastname:    faker.LastName(),
		FirstnameTH: faker.FirstName(),
		LastnameTH:  faker.LastName(),
		Ouid:        t.UserDto.StudentID,
	}

	a := &auth.Auth{
		UserID: t.UserDto.Id,
		Role:   t.Auth.Role,
	}

	repo := &mock.RepositoryMock{}
	repo.On("Create", a).Return(t.Auth, nil)

	chulaSSOClient := &mock.ChulaSSOClientMock{}
	chulaSSOClient.On("VerifyTicket", ticket, &dto.ChulaSSOCredential{}).Return(chulaSSORes, nil)

	in := &proto.User{
		StudentID: t.UserDto.StudentID,
		Faculty:   t.UserDto.Faculty,
		Year:      t.UserDto.Year,
	}

	userService := &mock.UserServiceMock{}
	userService.On("FindByStudentID", t.UserDto.StudentID).Return(nil, status.Error(codes.NotFound, "not found user"))
	userService.On("Create", in).Return(t.UserDto, nil)

	tokenService := &mock.TokenServiceMock{}
	tokenService.On("CreateOrUpdateCredentials", t.Auth).Return(t.Credential, nil)

	srv := NewService(repo, chulaSSOClient, tokenService, userService)
	actual, err := srv.VerifyTicket(context.Background(), &proto.VerifyTicketRequest{Ticket: ticket})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *AuthServiceTest) TestVerifyTicketSuccessNotFirstTimeLogin() {
	ticket := faker.Word()

	repo := &mock.RepositoryMock{}
	repo.On("FindByUserID", t.UserDto.Id).Return(t.Auth, nil)

	chulaSSOClient := &mock.ChulaSSOClientMock{}
	chulaSSOClient.On("VerifyTicket", ticket, &dto.ChulaSSOCredential{}).Return(nil, errors.New("Invalid Ticket"))

	userService := &mock.UserServiceMock{}
	userService.On("FindByStudentID", t.UserDto.Id).Return(t.UserDto, nil)

	tokenService := &mock.TokenServiceMock{}
	tokenService.On("CreateOrUpdateCredentials", t.Auth).Return(t.Credential, nil)

	srv := NewService(repo, chulaSSOClient, tokenService, userService)
	actual, err := srv.VerifyTicket(context.Background(), &proto.VerifyTicketRequest{Ticket: ticket})

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Unauthenticated, st.Code())
}

func (t *AuthServiceTest) TestVerifyTicketInvalid() {
	ticket := faker.Word()

	repo := &mock.RepositoryMock{}

	chulaSSOClient := &mock.ChulaSSOClientMock{}
	chulaSSOClient.On("VerifyTicket", ticket, &dto.ChulaSSOCredential{}).Return(nil, errors.New("Invalid Ticket"))

	userService := &mock.UserServiceMock{}
	userService.On("FindByStudentID", t.UserDto.Id).Return(nil, t.NotFoundErr)

	tokenService := &mock.TokenServiceMock{}

	srv := NewService(repo, chulaSSOClient, tokenService, userService)
	actual, err := srv.VerifyTicket(context.Background(), &proto.VerifyTicketRequest{Ticket: ticket})

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Unauthenticated, st.Code())
}

func (t *AuthServiceTest) TestVerifyTicketGrpcErr() {
	ticket := faker.Word()

	repo := &mock.RepositoryMock{}

	chulaSSOClient := &mock.ChulaSSOClientMock{}
	chulaSSOClient.On("VerifyTicket", ticket, &dto.ChulaSSOCredential{}).Return(nil, errors.New("Invalid Ticket"))

	userService := &mock.UserServiceMock{}
	userService.On("FindByStudentID", t.UserDto.Id).Return(nil, t.NotFoundErr)

	tokenService := &mock.TokenServiceMock{}

	srv := NewService(repo, chulaSSOClient, tokenService, userService)
	actual, err := srv.VerifyTicket(context.Background(), &proto.VerifyTicketRequest{Ticket: ticket})

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Unauthenticated, st.Code())
}
