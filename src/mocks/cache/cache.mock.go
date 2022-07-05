package cache

import (
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
	V map[string]interface{}
}

func (t *RepositoryMock) SaveCache(key string, v interface{}, ttl int) error {
	args := t.Called(key, v, ttl)

	log.Print(key)

	t.V[key] = v

	return args.Error(0)
}

func (t *RepositoryMock) GetCache(key string, v interface{}) error {
	args := t.Called(key, v)

	v = t.V[key]

	return args.Error(0)
}
