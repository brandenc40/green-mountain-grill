package grillclient

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// UDP connection constants
	_connReadDeadline  = 3 * time.Second
	_connWriteDeadline = 2 * time.Second
	_maxConnAttempts   = 5
)

// Client -
type Client interface {
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

// Params -
type Params struct {
	GrillIP   net.IP
	GrillPort int
	Logger    *logrus.Logger
}

// New -
func New(c Params) Client {
	client := &grillClient{
		grillAddr: &net.UDPAddr{
			IP:   c.GrillIP,
			Port: c.GrillPort,
		},
		logger: c.Logger,
	}
	if c.Logger == nil {
		client.logger = logrus.New()
	}
	return client
}

type grillClient struct {
	grillAddr *net.UDPAddr
	logger    *logrus.Logger
}

// GetState -
func (g *grillClient) GetState() (*State, error) {
	response, err := g.sendCommand(CommandGetInfo)
	if err != nil {
		return nil, err
	}
	return GetStateResponseToState(response), nil
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
		return nil, err
	}
	defer g.safeCloseConn(conn)

	// write the command to the udp connection
	cmd := command.Build(args...)
	n, err := conn.Write(cmd) // note: udp writes without confirmation of data transfer so this is non blocking
	if err != nil {
		return nil, err
	}
	g.logger.Debugf("%d bytes written: %v %#v %s", n, cmd, cmd, cmd)

	// read the response from the udp connection
	outBuf := make([]byte, 256)
	n, err = conn.Read(outBuf) // note: conn.Read() is blocking and will timeout after the _connReadDeadline duration
	if err != nil {
		return nil, err
	}
	outBuf = outBuf[:n] // trim the unused bytes
	g.logger.Debugf("%d bytes read: %v %#v", n, outBuf, outBuf)
	return outBuf, nil
}

func (g *grillClient) openConnectionWithRetries() (conn net.Conn, err error) {
	for i := 1; i <= _maxConnAttempts; i++ {
		g.logger.Debugf("udp conn attempt %d", i)
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
	if err = conn.SetReadDeadline(time.Now().Add(_connReadDeadline)); err != nil {
		return
	}
	if err = conn.SetWriteDeadline(time.Now().Add(_connWriteDeadline)); err != nil {
		return
	}
	g.logger.Debug("opened new udp connection")
	return
}

func (g *grillClient) safeCloseConn(conn net.Conn) {
	if err := conn.Close(); err != nil {
		g.logger.WithError(err).Error("err closing udp conn: ", err.Error())
	}
}
