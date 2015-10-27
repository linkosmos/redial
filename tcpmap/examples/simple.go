package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/linkosmos/redial/tcpmap"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	m, err := tcpmap.Lookup("www.example.org:80")
	if err != nil {
		fmt.Println("Lookup error", err)
		return
	}
	ip, err := m.Get(false)
	if err != nil {
		fmt.Println("TCMMap.Get error", err)
		return
	}
	fmt.Println("m.Get got", ip)
}
