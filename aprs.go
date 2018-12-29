package ogn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net/textproto"
	"regexp"
	"runtime"
	"time"
)

const (
	statusVerified   = "verified"
	statusUnverified = "unverified"
)

const (
	readTimeout         = 15 * time.Second
	maxReconnectRetries = 6
)

const (
	keepaliveInterval = 180 * time.Second
	keepaliveString   = "#keepalive"
)

const (
	bufferSize = 1000
)

var (
	ReadTimeoutError     = errors.New("read from APRS server timed out")
	ReconnectFailedError = errors.New("failed to reconnect to APRS server")
	ReadFailedError      = errors.New("read from APRS server failed")
	InvalidResponseError = errors.New("received invalid resposne from APRS server")
	LoginFailedError     = errors.New("login attempt failed")
)

var (
	loginConfimationPattern = regexp.MustCompile(`# logresp (?P<user>\w+) (?P<status>(verified)|(unverified)), server (?P<server>\w+)`)
)

type Urler interface {
	URL() string
}

type Auther interface {
	AuthString() string
	HasPassword() bool
}

func NewServer(host string, port int) Urler {
	return server{host: host, port: port}
}

func NewSettings(username, password, filter string) Auther {
	return settings{username: username, password: password, filter: filter}
}

func Connect(server Urler, settings Auther, withReconnect bool) AprsConnection {
	aprs := AprsConnection{
		server:        server,
		settings:      settings,
		withReconnect: withReconnect,
		buffer:        make(chan string),
	}

	return aprs
}

type server struct {
	host string
	port int
}

func (s server) URL() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

type settings struct {
	username string
	password string
	filter   string
}

func (s settings) AuthString() string {
	str := fmt.Sprintf("user %s pass %s ver %s %s", s.username, s.password, Package, Version)

	if filter := s.filter; filter != "" {
		str += fmt.Sprintf(" %s", filter)
	}

	return str
}

func (s settings) HasPassword() bool {
	return s.password != "" && s.password != "-1"
}

type AprsConnection struct {
	server        Urler
	settings      Auther
	withReconnect bool

	buffer chan string
}

func (ac *AprsConnection) Read() string {
	return <-ac.buffer
}

func (ac *AprsConnection) Receive(ctx context.Context) error {
	conn, err := ac.dial()

	for err == nil {
		err := ac.receiveInConnection(ctx, conn)
		log.Printf("read exited with error: %s", err)
		fmt.Println(runtime.NumGoroutine())

		select {
		case <-ctx.Done():
			return err
		default:
		}

		if !ac.withReconnect {
			return err
		}

		conn, err = ac.reconnect()
	}

	return err
}

func (ac *AprsConnection) receiveInConnection(ctx context.Context, conn *textproto.Conn) error {
	ctxConn := context.Background()
	ctxConn, cancel := context.WithCancel(ctxConn)
	defer cancel()

	msgCh := make(chan string)
	errCh := make(chan error)

	go keepAlive(ctxConn, conn)

	for {
		go receiveMessage(ctxConn, conn, msgCh, errCh)

		select {
		case msg := <-msgCh:
			ac.buffer <- msg
		case e := <-errCh:
			return e
		case <-time.After(readTimeout):
			return ReadTimeoutError
		case <-ctx.Done():
			return nil
		}
	}

	return nil
}

func receiveMessage(ctx context.Context, conn *textproto.Conn, msgCh chan<- string, errCh chan<- error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered panic in receiveMessage(): ", r)
		}
	}()

	msg, err := conn.ReadLine()

	select {
	case <-ctx.Done():
		return
	default:
		if err != nil {
			errCh <- err
			return
		}
		msgCh <- msg
	}
}

func (ac *AprsConnection) reconnect() (*textproto.Conn, error) {
	for i := 0; i < maxReconnectRetries; i++ {
		wait := math.Pow(2, float64(i))
		log.Printf("Attempting to reconnect in %.0f s", wait)
		time.Sleep(time.Duration(i) * time.Second)

		conn, err := ac.dial()

		if err == nil {
			log.Println("Successfully reconnected")
			return conn, err
		}
	}

	log.Println("Failed to reconnect")
	return nil, ReconnectFailedError
}

func (ac *AprsConnection) dial() (*textproto.Conn, error) {
	conn, err := textproto.Dial("tcp", ac.server.URL())
	if err != nil {
		return nil, err
	}

	line, err := conn.ReadLine()
	if err != nil {
		return nil, ReadFailedError
	}
	log.Printf("Received message: %s", line)

	conn.PrintfLine(ac.settings.AuthString())
	line, err = conn.ReadLine()
	if err != nil {
		return nil, ReadFailedError
	}
	data := unpackResponse(line, loginConfimationPattern)
	if data == nil {
		return nil, InvalidResponseError
	}
	log.Printf("Logged in as %s (%s), using server %s", data["user"], data["status"], data["server"])

	if ac.settings.HasPassword() && data["status"] != statusVerified {
		log.Println("Login attempt failed")
		return nil, LoginFailedError
	}

	return conn, nil
}

func keepAlive(ctx context.Context, conn *textproto.Conn) {
	ticker := time.NewTicker(keepaliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Println("sending keepalive beacon")
			conn.PrintfLine(keepaliveString)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func unpackResponse(line string, format *regexp.Regexp) map[string]string {
	result := make(map[string]string)

	matches := format.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	groups := format.SubexpNames()

	for i, group := range groups {
		result[group] = matches[i]
	}

	return result
}
