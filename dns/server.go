package dns

import (
	"log"

	"github.com/RTLer/highdns/dns/cache"
	"github.com/RTLer/highdns/dns/config"
	libdns "github.com/miekg/dns"
)

type DNSServer struct {
	Cache *cache.Cache
}

func (s *DNSServer) Serve() {
	libdns.HandleFunc(".", s.handleDnsRequest)

	server := &libdns.Server{Addr: config.Config.Server.Address, Net: config.Config.Server.Net}
	log.Printf("Starting at %s\n", config.Config.Server.Address)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}

func (s *DNSServer) handleDnsRequest(w libdns.ResponseWriter, r *libdns.Msg) {
	m := new(libdns.Msg)
	for _, q := range r.Question {
		m.SetQuestion(q.Name, q.Qtype)
	}

	switch r.Opcode {
	case libdns.OpcodeQuery:

		err := s.Cache.Remember(m)
		if err != nil {
			log.Printf("s.Cache.Remember: %+v", err)
			break
		}
	}
	m.SetReply(r)
	m.RecursionDesired = true

	w.WriteMsg(m)
}
