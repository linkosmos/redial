package tcpmap

import (
	"errors"
	"math/rand"
	"net"
	"strconv"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/bogdanovich/dns_resolver"
)

// -
var (
	ErrLookupFirst = errors.New("Lookup IP's first")
	ErrEmptyIPS    = errors.New("Given hostport has no IP's")
)

// RESOLVCONFIG - config with public DNS servers
// if file is found then will be used
const RESOLVCONFIG = "resolv.conf"

// Lookup - lookup host and initialize host IP's
func Lookup(hostport string) (*TCPMap, error) {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return nil, err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	resolver, err := dns_resolver.NewFromResolvConf(RESOLVCONFIG)
	if err != nil {
		resolver = dns_resolver.New([]string{"8.8.8.8", "8.8.4.4"})
	} else {
		logrus.Debugf("Found %s DNS config, attempting to use it", RESOLVCONFIG)
	}
	ips, err := resolver.LookupHost(host)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, ErrEmptyIPS
	}
	tcp := &TCPMap{}
	for _, ip := range ips {
		tcp.Add(ip, portInt)
	}
	return tcp, nil
}

// TCPMap - map of host+port and resolved IP's, added some syntactic
// sugar and convenience methods
type TCPMap struct {
	addresses   []*net.TCPAddr
	pointer     int
	pointerLock *sync.Mutex
}

func (t *TCPMap) lock() {
	if t.pointerLock == nil {
		t.pointerLock = &sync.Mutex{}
	}
	t.pointerLock.Lock()
}

func (t *TCPMap) unlock() {
	t.pointerLock.Unlock()
}

// Exist -
func (t *TCPMap) Exist(ip net.IP) bool {
	for _, address := range t.addresses {
		if address.IP.Equal(ip) {
			return true
		}
	}
	return false
}

// Add -
func (t *TCPMap) Add(ip net.IP, port int) {
	if t.addresses == nil {
		t.addresses = make([]*net.TCPAddr, 0, 0)
	}
	if !t.Exist(ip) {
		t.addresses = append(t.addresses, &net.TCPAddr{IP: ip, Port: port})
	}
}

// Size - count of cached addresses
func (t *TCPMap) Size() int {
	t.lock()
	defer t.unlock()
	return len(t.addresses)
}

// GetRoundRobin -
func (t *TCPMap) GetRoundRobin() (address *net.TCPAddr, err error) {
	if t.addresses == nil {
		return nil, ErrLookupFirst
	}
	size := t.Size()
	if size == 0 {
		return nil, ErrEmptyIPS
	}
	if size == 1 {
		address = t.addresses[0]
		logrus.Debugf("Found single address: %s", address)
		return t.addresses[0], nil
	}
	t.lock()
	defer t.unlock()
	if t.pointer >= size {
		t.pointer = 0
	}
	address = t.addresses[t.pointer]
	logrus.Debugf("Returning RoundRobin (pointer %d): %s", t.pointer, address)
	t.pointer++
	return address, nil
}

// Get - returns random IP address if more than 1
func (t *TCPMap) Get(randomize bool) (address *net.TCPAddr, err error) {
	if t.addresses == nil {
		return nil, ErrLookupFirst
	}
	size := t.Size()
	if size == 0 {
		return nil, ErrEmptyIPS
	}
	if size == 1 {
		address = t.addresses[0]
		logrus.Debugf("Found single address: %s", address)
		return t.addresses[0], nil
	}
	if randomize {
		address = t.addresses[rand.Intn(size)]
		logrus.Debugf("Returning randomly: %s", address)
		return address, nil
	}
	return t.addresses[0], nil
}
