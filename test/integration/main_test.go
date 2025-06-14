package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestMain(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
