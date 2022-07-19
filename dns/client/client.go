package client

import (
	libdns "github.com/miekg/dns"
)

type Client interface {
	Init()
	Acquire(query *libdns.Msg) (err error)
}
