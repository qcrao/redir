package main

import (
	"container/list"
	"sync"
	"time"
)

type item struct {
	k, v string
}

// lru is a naive thread-safe lru cache
type lru struct {
	cap   uint
	size  uint
	elems *list.List // of item

	mu sync.RWMutex
}

func newLRU(doexpire bool) *lru {
	l := &lru{
		cap:   32, // could do it with memory quota
		size:  0,
		elems: list.New(),
		mu:    sync.RWMutex{},
	}
	if doexpire {
		go l.clear()
	}
	return l
}

// clear clears the lru after a while, this is just a dirty
// solution to prevent if the database is updated but lru is
// not synced.
func (l *lru) clear() {
	t := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-t.C:
			l.mu.Lock()
			for e := l.elems.Front(); e != nil; e = e.Next() {
				l.elems.Remove(e)
			}
			l.mu.Unlock()
		}
	}
}

func (l *lru) Get(k string) (string, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for e := l.elems.Front(); e != nil; e = e.Next() {
		if e.Value.(*item).k == k {
			l.elems.MoveToFront(e)
			return e.Value.(*item).v, true
		}
	}
	return "", false
}

func (l *lru) Put(k, v string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// found from cache
	i := &item{k, v}
	for e := l.elems.Front(); e != nil; e = e.Next() {
		if e.Value.(*item).k == k {
			i.v = l.elems.Remove(e).(*item).v
			l.elems.PushFront(i)
			return
		}
	}

	// check if cache is full
	if l.size+1 > l.cap {
		l.elems.Remove(l.elems.Back())
	}
	l.elems.PushFront(i)
	l.size++
	return
}