package client

import (
	"fmt"
	"log"

	"github.com/RTLer/highdns/dns/config"
	libdns "github.com/miekg/dns"
)

type DNS struct {
	client *libdns.Client
}

func (a *DNS) Init() {
	a.client = new(libdns.Client)
}

func (a *DNS) Acquire(query *libdns.Msg) (err error) {

	hosts := config.Config.GetHosts(query.Question[0].Name).Hosts

	log.Println(query.Question[0].Name)
	for _, address := range hosts {
		r, _, err := a.client.Exchange(query, address)
		if err == nil {
			query.Answer = r.Answer
			return nil
		}
	}

	return fmt.Errorf("can not acquire any response for this query from %v", hosts)
}
