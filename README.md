# redial

Golang augmented net.Dialer with connection pooling, cached DNS record IP's with round-robin &amp; random access.

Notice: you must use single dialer per host

[![Build Status](https://travis-ci.org/linkosmos/redial.svg?branch=master)](https://travis-ci.org/linkosmos/redial)
[![GoDoc](http://godoc.org/github.com/linkosmos/redial?status.svg)](http://godoc.org/github.com/linkosmos/redial)
[![Go Report Card](http://goreportcard.com/badge/linkosmos/redial)](http://goreportcard.com/report/linkosmos/redial)
[![BSD License](http://img.shields.io/badge/license-BSD-blue.svg)](http://opensource.org/licenses/BSD-3-Clause)

### Methods
 - New(d net.Dialer, address, port string, cPoolInitial, cPoolMax int) (*Dialer, error)
 - Dial(network, address string) (net.Conn, error)
 - Close()
 - String() string

### Example

```go
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
```


### License

Copyright (c) 2015, linkosmos
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of redial nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
