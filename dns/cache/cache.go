package cache

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
	libdns "github.com/miekg/dns"
)

type Cache struct {
	db      *badger.DB
	Queue   AcquireQueue
	Acquire func(query *libdns.Msg) (err error)
}

func (c *Cache) Open() error {
	var err error
	c.Queue = AcquireQueue{
		queue:   []Payload{},
		acquire: c.Acquire,
	}
	c.db, err = badger.Open(badger.DefaultOptions("").WithInMemory(true))
	return err
}

func (c *Cache) Close() error {
	return c.db.Close()
}

func (c *Cache) Remember(query *libdns.Msg) error {
	var err error

	recordKey := c.generateKey(query)

	val, err := c.get(recordKey)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return c.runQuery(query)
		}
		return err

	}

	tr, err := c.decode(val)
	rr := []libdns.RR{}
	for _, t := range tr {
		t.rr.Header().Ttl = uint32(t.ttl - uint64(time.Now().Unix()))
		if t.ttl <= uint64(time.Now().Unix()) {
			t.rr.Header().Ttl = 1
		}
		rr = append(rr, t.rr)
	}
	if len(tr) > 0 && tr[0].ttl <= uint64(time.Now().Unix()) {
		c.Queue.Add(recordKey, *query)
	}
	query.Answer = rr
	return err
}
func (c *Cache) set(recordKey string, val []byte) error {
	return c.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(recordKey), []byte(val)).WithTTL(time.Hour * 24)
		return txn.SetEntry(e)

	})

}
func (c *Cache) get(recordKey string) ([]byte, error) {
	var val []byte
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(recordKey))
		if err != nil {
			return err
		}

		return item.Value(func(v []byte) error {
			val = append([]byte{}, v...)
			return nil
		})
	})

	return val, err
}

func (c *Cache) encode(val []libdns.RR) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	var pval = [][]byte{}

	for _, i := range val {
		msg := make([]byte, 108)
		binary.LittleEndian.PutUint64(msg, uint64(time.Now().Unix()+int64(i.Header().Ttl)))
		libdns.PackRR(i, msg, 8, make(map[string]int), false)

		pval = append(pval, msg)
	}

	if err := enc.Encode(pval); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil

}

type TimedRR struct {
	rr  libdns.RR
	ttl uint64
}

func (c *Cache) decode(val []byte) ([]TimedRR, error) {
	buf := bytes.NewBuffer(val)
	dec := gob.NewDecoder(buf)

	var pval = [][]byte{}
	if err := dec.Decode(&pval); err != nil {
		return nil, err
	}

	res := []TimedRR{}
	for _, msg := range pval {
		r, _, err := libdns.UnpackRR(msg, 8)
		if err != nil {
			return nil, err
		}

		res = append(res, TimedRR{
			rr:  r,
			ttl: binary.LittleEndian.Uint64(msg[:8]),
		})
	}

	return res, nil
}

func (c *Cache) generateKey(query *libdns.Msg) string {
	return fmt.Sprintf("%d;%d;%s", query.Question[0].Qclass, query.Question[0].Qtype, query.Question[0].Name)
}

func (c *Cache) runQuery(query *libdns.Msg) error {
	recordKey := c.generateKey(query)

	err := c.Acquire(query)
	if err != nil {
		return err
	}
	res, err := c.encode(query.Answer)
	if err != nil {
		return err
	}

	return c.set(recordKey, res)
}
