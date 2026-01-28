package client

import (
	"sync"
)

type InMemoryKeystore struct {
	l sync.Mutex
	d map[string]JwkKey
}

func (ms *InMemoryKeystore) Get(kid string) (*JwkKey, error) {
	k, ok := ms.d[kid]
	if !ok {
		return nil, ErrNotFound
	}
	return &k, nil
}

func (ms *InMemoryKeystore) Set(k *JwkKey) error {
	ms.l.Lock()
	ms.d[k.Kid] = *k
	ms.l.Unlock()
	return nil
}

func (ms *InMemoryKeystore) Pop(kid string) error {
	ms.l.Lock()
	delete(ms.d, kid)
	ms.l.Unlock()
	return nil
}

func (ms *InMemoryKeystore) Clean() {
	ms.l.Lock()
	ms.d = make(map[string]JwkKey)
	ms.l.Unlock()
}

func (ms *InMemoryKeystore) Exist(s string) bool {
	ms.l.Lock()
	_, exists := ms.d[s]
	ms.l.Unlock()
	return exists
}
