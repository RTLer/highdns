package cache

import (
	"log"
	"sync"
	"time"

	libdns "github.com/miekg/dns"
)

type Payload struct {
	key string
	msg *libdns.Msg
}
type AcquireQueue struct {
	queue   []Payload
	lock    sync.Mutex
	acquire func(query *libdns.Msg) (err error)
}

func (aq *AcquireQueue) Run(c *Cache) {
	for {
		if len(aq.queue) > 0 {
			aq.lock.Lock()
			payload := aq.queue[len(aq.queue)-1]
			aq.queue = aq.queue[:len(aq.queue)-1]
			if err := c.runQuery(payload.msg); err != nil {
				log.Printf("c.runQuery(payload.msg) error: %s", err)
			}
			aq.lock.Unlock()

			continue
		}
		time.Sleep(time.Millisecond * 10)
	}

}
func (aq *AcquireQueue) Add(key string, msg libdns.Msg) {
	aq.lock.Lock()
	defer aq.lock.Unlock()

	for _, i := range aq.queue {
		if i.key == key {
			return
		}
	}
	aq.queue = append(aq.queue, Payload{
		key: key,
		msg: &msg,
	})
}
