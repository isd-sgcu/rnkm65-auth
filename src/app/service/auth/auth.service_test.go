package auth

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type AuthServiceTest struct {
	suite.Suite
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTest))
}
