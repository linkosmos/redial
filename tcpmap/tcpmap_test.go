package tcpmap

import (
	"net"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func testIPs() (output []net.IP) {
	output = append(output, net.IPv4(127, 0, 1, 2))
	output = append(output, net.IPv4(128, 0, 1, 2))
	output = append(output, net.IPv4(129, 0, 1, 2))
	output = append(output, net.IPv4(130, 0, 1, 2))
	return output
}

func TestAdd(t *testing.T) {
	tcpMap := &TCPMap{}

	for _, ip := range testIPs() {
		tcpMap.Add(ip, 443)
	}
	for _, ip := range testIPs() {
		assert.True(t, tcpMap.Exist(ip), "Expected %s to exist", ip)
	}
}

func TestGetRobinRobing(t *testing.T) {
	tcpMap := &TCPMap{}

	ip1 := testIPs()[0]
	ip2 := testIPs()[1]
	ip3 := testIPs()[2]

	tcpMap.Add(ip1, 80)
	tcpMap.Add(ip2, 80)
	tcpMap.Add(ip3, 80)

	// Should return ip1
	got, _ := tcpMap.GetRoundRobin()
	assert.True(t, got.IP.Equal(ip1), "Expected %s, got %s", ip1, got.IP)
	// Should return ip2
	got, _ = tcpMap.GetRoundRobin()
	assert.True(t, got.IP.Equal(ip2), "Expected %s, got %s", ip2, got.IP)
	// Should return ip3
	got, _ = tcpMap.GetRoundRobin()
	assert.True(t, got.IP.Equal(ip3), "Expected %s, got %s", ip3, got.IP)
	// Should return ip1
	got, _ = tcpMap.GetRoundRobin()
	assert.True(t, got.IP.Equal(ip1), "Expected %s, got %s", ip1, got.IP)
}
