package integration

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(Suite))
}
