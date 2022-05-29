package gmg

import (
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
)

const (
	// default UDP connection constants
	_connReadDeadline  = 2 * time.Second
	_connWriteDeadline = time.Second
	_maxConnAttempts   = 5
)

// Client - Green Mountain Grill client interface definition
type Client interface {
	IsAvailable() bool
	GetState() (*State, error)
	GetID() ([]byte, error)
	GetFirmware() ([]byte, error)
	SetGrillTemp(temp int) error
	SetProbe1Target(temp int) error
	SetProbe2Target(temp int) error
	PowerOn() error
	PowerOnColdSmoke() error
	PowerOff() error
}

type Option interface {
	apply(*grillClient)
}

// WithReadTimeout -
func WithReadTimeout(timeout time.Duration) Option {
	return readTimeoutOption{readTimeout: timeout}
}

type readTimeoutOption struct{ readTimeout time.Duration }

func (r readTimeoutOption) apply(g *grillClient) { g.readTimeout = r.readTimeout }

// WithWriteTimeout -
func WithWriteTimeout(timeout time.Duration) Option {
	return writeTimeoutOption{writeTimeout: timeout}
}

type writeTimeoutOption struct{ writeTimeout time.Duration }

func (w writeTimeoutOption) apply(g *grillClient) { g.writeTimeout = w.writeTimeout }

// WithMaxConnectionAttempts -
func WithMaxConnectionAttempts(attempts int) Option {
	return maxConnAttempts{attempts: attempts}
}

type maxConnAttempts struct{ attempts int }

func (m maxConnAttempts) apply(g *grillClient) { g.maxConnAttempts = m.attempts }

// WithZapLogger -
func WithZapLogger(l *zap.Logger) Option {
	return withZapLogger{logger: l}
}

type withZapLogger struct{ logger *zap.Logger }

func (l withZapLogger) apply(g *grillClient) {
	if l.logger != nil {
		g.logger = l.logger
	}
}

// New -
func New(grillIP net.IP, grillPort int, options ...Option) (Client, error) {
	client := grillClient{
		grillAddr:       &net.UDPAddr{IP: grillIP, Port: grillPort},
		readTimeout:     _connReadDeadline,
		writeTimeout:    _connWriteDeadline,
		maxConnAttempts: _maxConnAttempts,
	}
	for _, option := range options {
		option.apply(&client)
	}
	if client.logger == nil {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		client.logger = logger
	}
	return &client, nil
}

type grillClient struct {
	grillAddr       *net.UDPAddr
	logger          *zap.Logger
	readTimeout     time.Duration
	writeTimeout    time.Duration
	maxConnAttempts int
}

// IsAvailable -
func (g *grillClient) IsAvailable() bool {
	_, err := g.GetState()
	return err == nil
}

// GetState -
func (g *grillClient) GetState() (*State, error) {
	response, err := g.sendCommand(CommandGetInfo)
	if err != nil {
		return nil, err
	}
	if len(response) != 36 {
		return nil, fmt.Errorf("expected 36 bytes, got %d", len(response))
	}
	return BytesToState(response)
}

// GetID -
func (g *grillClient) GetID() ([]byte, error) {
	response, err := g.sendCommand(CommandGetGrillID)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetFirmware -
func (g *grillClient) GetFirmware() ([]byte, error) {
	response, err := g.sendCommand(CommandGetGrillFirmware)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// SetGrillTemp -
func (g *grillClient) SetGrillTemp(temp int) error {
	_, err := g.sendCommand(CommandSetGrillTemp, temp)
	return err
}

// SetProbe1Target -
func (g *grillClient) SetProbe1Target(temp int) error {
	_, err := g.sendCommand(CommandSetProbe1Temp, temp)
	return err
}

// SetProbe2Target -
func (g *grillClient) SetProbe2Target(temp int) error {
	_, err := g.sendCommand(CommandSetProbe2Temp, temp)
	return err
}

// PowerOn -
func (g *grillClient) PowerOn() error {
	_, err := g.sendCommand(CommandPowerOn)
	return err
}

// PowerOnColdSmoke -
func (g *grillClient) PowerOnColdSmoke() error {
	_, err := g.sendCommand(CommandPowerOnColdSmoke)
	return err
}

// PowerOff -
func (g *grillClient) PowerOff() error {
	_, err := g.sendCommand(CommandPowerOff)
	return err
}

func (g *grillClient) sendCommand(command Command, args ...interface{}) ([]byte, error) {
	// open a new udp connection
	conn, err := g.openConnectionWithRetries()
	if err != nil {
		g.logger.Error("grill is unreachable", zap.Error(err))
		return nil, GrillUnreachableErr{Err: err}
	}
	defer g.safeCloseConn(conn)

	// write the command to the udp connection
	cmd := command.Build(args...)
	n, err := conn.Write(cmd) // note: udp writes without confirmation of data transfer so this is non blocking
	if err != nil {
		g.logger.Error("unable to write to udp conn", zap.Error(err))
		return nil, GrillUnreachableErr{Err: err}
	}
	g.logger.Debug(fmt.Sprintf("%d bytes written: %v %#v %s", n, cmd, cmd, cmd))

	// read the response from the udp connection
	outBuf := make([]byte, 64)
	n, err = conn.Read(outBuf) // note: conn.Read() is blocking and will timeout after the _connReadDeadline duration
	if err != nil {
		g.logger.Error("unable to read from udp conn", zap.Error(err))
		return nil, GrillUnreachableErr{Err: err}
	}
	outBuf = outBuf[:n] // trim the unused bytes
	g.logger.Debug(fmt.Sprintf("%d bytes read: %v %#v", n, outBuf, outBuf))
	return outBuf, nil
}

func (g *grillClient) openConnectionWithRetries() (conn net.Conn, err error) {
	for i := 1; i <= _maxConnAttempts; i++ {
		g.logger.Debug(fmt.Sprintf("udp conn attempt: %d", i))
		conn, err = g.openConnection()
		if err == nil || i == _maxConnAttempts {
			return
		}
		continue
	}
	return
}

func (g *grillClient) openConnection() (conn net.Conn, err error) {
	if conn, err = net.DialUDP("udp4", nil, g.grillAddr); err != nil {
		return
	}
	now := time.Now()
	if err = conn.SetReadDeadline(now.Add(_connReadDeadline)); err != nil {
		return
	}
	if err = conn.SetWriteDeadline(now.Add(_connWriteDeadline)); err != nil {
		return
	}
	g.logger.Debug("opened new udp connection")
	return
}

func (g *grillClient) safeCloseConn(conn net.Conn) {
	if err := conn.Close(); err != nil {
		g.logger.Error("err closing udp conn", zap.Error(err))
	}
}
