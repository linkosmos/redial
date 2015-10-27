package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/linkosmos/redial"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	url, _ := url.Parse("http://www.example.org")
	port := "80"
	initialConnectionPool := 1
	maxConnectionsInPool := 3

	dialer, err := redial.New(net.Dialer{}, url.Host, port, initialConnectionPool, maxConnectionsInPool)
	if err != nil {
		fmt.Println("redial.New error", err)
		return
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		fmt.Println("http.NewRequest", err)
		return
	}

	response, err := transport.RoundTrip(request)
	if err != nil {
		fmt.Println("transport.RoundTrip", err)
		return
	}
	defer response.Body.Close()

	fmt.Println("Response status", response.Status, "for", url)
}
