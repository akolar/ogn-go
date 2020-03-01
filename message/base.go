package message

import (
	"strings"
	"time"
)

type Value interface{}

// AircraftBase is the fundation for all APRS messages originating from
// airplane beacons.
// It is implemented by all derived types (which may provide additional
// methods).
type AircraftBase interface {
	// Sender retunrs the unique identifier of the sender of the message.
	Sender() string

	// Type returns the format type using which the message was encoded.
	Type() string

	// Relayer returns the unique identifier of the aircraft relaying the
	// message. If Relayer returns an empty string, the message was not
	// relayed.
	Relayer() string

	// Receiver returns the name of the ground station that received the
	// message.
	Receiver() string

	// Time returns the timestamp of the message.
	Time() time.Time

	// Latitude returns the latitude of the aircraft in decimal degrees, so
	// that -90 <= Latitude() <= 90.
	Latitude() float64

	// Latitude returns the longitude of the aircraft in decimal degrees, so
	// that -180 <= Longitude() <= 180.
	Longitude() float64

	// Heading returns the aircraft's heading in degrees, so that 1 <=
	// Heading() <= 360. Heading=0 indicates no data
	Heading() int

	// Speed returns the speed of the aircraft (as reported by it) in m/s.
	Speed() float64

	// Altitude returns the altitude of the aircraft (as reported by it) in m.
	Altitude() int

	// AsMap returns a generic mapping between the message's fields and their
	// corresponding values.
	AsMap() map[string]Value

	// AsAPRS encodes the message as an APRS string.
	AsAPRS() string
}

type ReceiverBase interface {
	// Name returns the name of the ground station.
	Name() string

	// Type returns the format type using which the message was encoded.
	Type() string

	// Server returns the identifier of the server that relayed the message.
	Server() string

	// Latitude returns the latitude of the station in decimal degrees, so that
	// -90 <= Latitude() <= 90.
	Latitude() float64

	// Latitude returns the longitude of the station in decimal degrees, so
	// that -180 <= Longitude() <= 180.
	Longitude() float64

	// Altitude returns the altitude of the station in m.
	Altitude() int
}

func extractMessageType(message string) (string, error) {
	start := strings.Index(message, ">")
	if start < 0 {
		return "", InvalidFormatError{msg: message}
	}

	end := strings.Index(message[start:], ",")
	if end < 0 {
		return "", InvalidFormatError{msg: message}
	}

	return message[start+1 : start+end], nil
}
