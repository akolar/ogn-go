package message

// Parse parses the message into one of the message interfaces. If the parser
// for this message type is not registered, Parse returns an error.
func Parse(message string) (interface{}, error) {
	var m interface{}

	msgType, err := extractMessageType(message)
	if err != nil {
		return m, err
	}

	if p, ok := registeredParsers[msgType]; ok {
		m, err = p(message)
	} else {
		err = ParserNotFoundError{msgType: msgType}
	}

	return m, err
}

// ParserFunc is the type of parser functions. It takes the message as its only
// argument and returns a parsed message and a potential error.
type ParserFunc func(string) (interface{}, error)

type parsers map[string]ParserFunc

var registeredParsers = parsers{}

// Handle registers a new parser for the given typeName. It returns an
// error if a parser for this type already exists.
func Handle(typeName string, parser ParserFunc) error {
	if _, ok := registeredParsers[typeName]; ok {
		return ParserExistsError{msgType: typeName}
	}

	registeredParsers[typeName] = parser
	return nil
}
