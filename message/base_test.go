package message

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BaseTestSuite struct {
	suite.Suite
}

func (s *BaseTestSuite) TestExtractMessageType() {
	extracted, err := extractMessageType("FLRDD89C9>TYPE,qAS")
	s.Equal(extracted, "TYPE")
	s.Nil(err)
}

func (s *BaseTestSuite) TestExtractMessageTypeMissingStart() {
	_, err := extractMessageType("FLRDD89C9TYPE,qAS")
	s.IsType(err, InvalidFormatError{})
}

func (s *BaseTestSuite) TestExtractMessageTypeMissingEnd() {
	_, err := extractMessageType("FLRDD89C9>TYPEqAS")
	s.IsType(err, InvalidFormatError{})
}

func TestBase(t *testing.T) {
	suite.Run(t, new(BaseTestSuite))
}
