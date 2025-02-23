/**
 * File: pool.go
 * Author: Ming Cheng<mingcheng@outlook.com>
 *
 * Created Date: Tuesday, June 21st 2022, 6:03:26 pm
 * Last Modified: Thursday, July 7th 2022, 6:47:39 pm
 *
 * http://www.opensource.org/licenses/MIT
 */

package socks5lb

import (
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

type Pool struct {
	backends []*Backend
	current  uint64
}

func (b *Pool) Add(backend *Backend) {
	b.backends = append(b.backends, backend)
}

func (b *Pool) NextIndex() int {
	return int(atomic.AddUint64(&b.current, uint64(1)) % uint64(len(b.backends)))
}

func (b *Pool) Next() *Backend {
	// loop entire backends to find out an Alive backend
	next := b.NextIndex()
	l := len(b.backends) + next // start from next and move a full cycle
	for i := next; i < l; i++ {
		idx := i % len(b.backends)   // take an index by modding
		if b.backends[idx].Alive() { // if we have an alive backend, use it and store if its not the original one
			if i != next {
				atomic.StoreUint64(&b.current, uint64(idx))
			}
			return b.backends[idx]
		}
	}

	return nil
}

func (b *Pool) Check() {
	for _, b := range b.backends {
		err := b.Check()
		if err != nil {
			log.Errorf("check backend %s is failed, error %v", b.Addr, err)
		} else {
			log.Debugf("check backend %s is successful", b.Addr)
		}
	}
}

func NewPool() *Pool {
	return &Pool{
		backends: []*Backend{},
	}
}
