package redial

import (
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/facebookgo/httpcontrol"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func TestDial(t *testing.T) {
	u, _ := url.Parse("http://www.google.com")
	dialer, err := New(net.Dialer{
		KeepAlive: 3 * time.Second,
	}, u.Host, "80", 2, 4)
	assert.Nil(t, err, "Got error on init %s", err)

	transport := &httpcontrol.Transport{}
	transport.Dial = dialer.Dial

	request, _ := http.NewRequest("GET", u.String(), nil)

	_, err = transport.RoundTrip(request)
	assert.Nil(t, err, "Got error %s", err)
}
