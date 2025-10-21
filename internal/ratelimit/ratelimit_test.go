package ratelimit_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type rateLimitSuite struct {
	suite.Suite
}

func TestRateLimitSuite(t *testing.T) {
	suite.Run(t, &rateLimitSuite{})
}
