package mock

import (
	"cobaMetrics/app/config"
	"github.com/stretchr/testify/mock"
)

type ConfigMock struct {
	Mock *mock.Mock
}

func NewConfigMock() *ConfigMock {
	return &ConfigMock{&mock.Mock{}}
}

func (c *ConfigMock) Config() *config.ConfigApp {
	args := c.Mock.Called()

	return args.Get(0).(*config.ConfigApp)
}
