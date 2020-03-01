package message

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	testParserValue = "ok"
)

func testParser(msg string) (interface{}, error) {
	return testParserValue, nil
}

type ParserTestSuite struct {
	suite.Suite
}

func (s *ParserTestSuite) SetupTest() {
	registeredParsers = parsers{}
}

func (s *ParserTestSuite) TearDownTest() {
	registeredParsers = parsers{}
}

func (s *ParserTestSuite) TestRegisterType() {
	err := Handle("TEST-1", testParser)
	s.Nil(err)
}

func (s *ParserTestSuite) TestRegisterTypeDuplicate() {
	err := Handle("TEST-1", testParser)
	s.Nil(err)
	err = Handle("TEST-1", testParser)
	s.IsType(ParserExistsError{}, err)
}

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}
