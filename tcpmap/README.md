# tcpmap

[![Build Status](https://travis-ci.org/linkosmos/redial/tcpmap.svg?branch=master)](https://travis-ci.org/linkosmos/redial/tcpmap)
[![GoDoc](http://godoc.org/github.com/linkosmos/redial/tcpmap?status.svg)](http://godoc.org/github.com/linkosmos/redial/tcpmap)
[![BSD License](http://img.shields.io/badge/license-BSD-blue.svg)](http://opensource.org/licenses/BSD-3-Clause)

## Cached DNS IP records with simple, round-robin OR random access

### Methods

 - Lookup(hostport string) (*TCPMap, error)
 - Exist(ip net.IP) bool
 - Add(ip net.IP, port int)
 - Size() int
 - GetRoundRobin() (*net.TCPAddr, error)
 - Get(randomize bool) (*net.TCPAddr, error)

### Example

```go
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
		fmt.Println("TCPMap.Get error", err)
		return
	}
	fmt.Println("m.Get got", ip)
}
```

### Custon DNS nameservers can be defined in resolv.conf

```resolv.conf
nameserver 8.8.8.8
nameserver 8.8.4.4

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
