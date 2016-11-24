package config

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceConfigSuite struct {
	suite.Suite
	service *ServiceConfig
}

func TestServiceConfigSuite(t *testing.T) {
	suite.Run(t, new(ServiceConfigSuite))
}

func (s *ServiceConfigSuite) SetupTest() {
	s.service = NewServiceConfig()

}
