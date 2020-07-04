package kmutex

import "sync"

type Kmutex struct {
	c *sync.Cond
	l sync.Locker
	s map[interface{}]struct{}
}

func New() *Kmutex {
	l := sync.Mutex{}
	return &Kmutex{c: sync.NewCond(&l), l: &l, s: make(map[interface{}]struct{})}
}

func (km *Kmutex) locked(key interface{}) (ok bool) { _, ok = km.s[key]; return }

func (km *Kmutex) Unlock(key interface{}) {
	km.l.Lock()
	defer km.l.Unlock()
	delete(km.s, key)
	km.c.Broadcast()
}

func (km Kmutex) Lock(key interface{}) {
	km.l.Lock()
	defer km.l.Unlock()
	for km.locked(key) {
		km.c.Wait()
	}
	km.s[key] = struct{}{}
	return
}
