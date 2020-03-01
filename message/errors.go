package message

import (
	"fmt"
)

type ParserNotFoundError struct {
	msgType string
	message string
}

func (e ParserNotFoundError) Error() string {
	return fmt.Sprintf("parser for message type %s not found: %s", e.msgType, e.message)
}

type ParserExistsError struct {
	msgType string
}

func (e ParserExistsError) Error() string {
	return fmt.Sprintf("parser for message type %s already exists", e.msgType)
}

type InvalidFormatError struct {
	msg string
}

func (e InvalidFormatError) Error() string {
	return fmt.Sprintf("cannot parse message: %s", e.msg)
}
