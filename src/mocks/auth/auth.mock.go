package auth

import (
	"context"
	dto "github.com/isd-sgcu/rnkm65-auth/src/app/dto/auth"
	model "github.com/isd-sgcu/rnkm65-auth/src/app/model/auth"
	"github.com/isd-sgcu/rnkm65-auth/src/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) FindByStudentID(id string, result *model.Auth) error {
	args := r.Called(id, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*model.Auth)
	}

	return args.Error(1)
}

func (r *RepositoryMock) FindByRefreshToken(id string, result *model.Auth) error {
	args := r.Called(id, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*model.Auth)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Create(in *model.Auth) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*model.Auth)
	}

	return args.Error(1)
}

type ChulaSSOClientMock struct {
	mock.Mock
}

func (c *ChulaSSOClientMock) VerifyTicket(ticket string, result *dto.ChulaSSOCredential) error {
	args := c.Called(ticket, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*dto.ChulaSSOCredential)
	}

	return args.Error(1)
}

type UserClientMock struct {
	mock.Mock
}

func (c *UserClientMock) FindOne(_ context.Context, in *proto.FindOneUserRequest, _ ...grpc.CallOption) (res *proto.FindOneUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindOneUserResponse)
	}

	return res, args.Error(1)
}

func (c *UserClientMock) Create(_ context.Context, in *proto.CreateUserRequest, _ ...grpc.CallOption) (res *proto.CreateUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.CreateUserResponse)
	}

	return res, args.Error(1)
}

func (c *UserClientMock) CreateOrUpdate(_ context.Context, in *proto.CreateOrUpdateUserRequest, _ ...grpc.CallOption) (res *proto.CreateOrUpdateUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.CreateOrUpdateUserResponse)
	}

	return res, args.Error(1)
}
