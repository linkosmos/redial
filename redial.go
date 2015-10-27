package redial

import (
	"errors"
	"fmt"
	"net"

	"gopkg.in/fatih/pool.v2"

	"github.com/Sirupsen/logrus"
	"github.com/linkosmos/redial/tcpmap"
)

// -
var (
	ErrUseNewToInitialize = errors.New("Use redial.New to initalize")
)

// Dialer -
type Dialer struct {
	net.Dialer
	cPool pool.Pool
	aPool *tcpmap.TCPMap
}

// New - initalize dial.Dialer wrapper
func New(d net.Dialer, address, port string, cPoolInitial, cPoolMax int) (dialer *Dialer, err error) {
	dialer = &Dialer{
		Dialer: d,
	}
	p, err := tcpmap.Lookup(net.JoinHostPort(address, port))
	if err != nil {
		return nil, err
	}
	dialer.aPool = p
	cp, err := pool.NewChannelPool(cPoolInitial, cPoolMax,
		func() (net.Conn, error) { return dialer.dialFactory("tcp", address) })
	if err != nil {
		return nil, err
	}
	dialer.cPool = cp
	logrus.Debugf("net.Dialer initialized with %s", dialer)
	return dialer, nil
}

func (d *Dialer) dialFactory(network, address string) (net.Conn, error) {
	tcpAddr, err := d.aPool.GetRoundRobin()
	if err != nil {
		logrus.Warningf("aPool.GetRoundRobin error", err)
		return d.Dialer.Dial(network, address)
	}
	c, err := net.DialTCP(network, nil, tcpAddr)
	if err != nil {
		logrus.Warningf("Failed to setup DialTCP: %s, fallback dialer.Dial", err)
		return d.Dialer.Dial(network, address)
	}
	if d.KeepAlive != 0 {
		c.SetKeepAlive(true)
		c.SetKeepAlivePeriod(d.KeepAlive)
		c.SetLinger(0)
		c.SetNoDelay(true)
	}
	return c, err
}

// Dial - lightweight version of dialer.Dial, this has cached
// dns hostport and shorter TCP connection setup
func (d *Dialer) Dial(network, address string) (net.Conn, error) {
	if d.aPool == nil || d.cPool == nil {
		panic(ErrUseNewToInitialize)
	}
	return d.cPool.Get()
}

// Close - closes connection pool
func (d *Dialer) Close() {
	d.cPool.Close()
}

func (d *Dialer) String() string {
	return fmt.Sprintf(`
Timeout %f
Deadline %s
LocalAddr %s
DualStack %t
FallbackDelay %f
KeepAlive %f`, d.Timeout.Seconds(), d.Deadline, d.LocalAddr,
		d.DualStack, d.FallbackDelay.Seconds(), d.KeepAlive.Seconds())
}
