package lru

import (
	"container/list"
	"sync"
)

type LRUCacher interface {
	Add(key, value string) bool
	Get(key string) (string, bool)
	Remove(key string) bool
	Len() int
}

type Node struct {
	Data   string
	KeyPtr *list.Element
}

type LRUCache struct {
	Queue    *list.List
	Items    map[string]*Node
	Capacity int
	lock     sync.RWMutex
}

func (l *LRUCache) Len() int {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.Queue.Len()
}

func (l *LRUCache) Add(key, value string) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	if item, ok := l.Items[key]; !ok {
		if l.Capacity == len(l.Items) {
			back := l.Queue.Back()
			l.Queue.Remove(back)
			delete(l.Items, back.Value.(string))
		}
		l.Items[key] = &Node{Data: value, KeyPtr: l.Queue.PushFront(key)}
	} else {
		item.Data = value
		l.Items[key] = item
		l.Queue.MoveToFront(item.KeyPtr)
	}
	return true
}

func (l *LRUCache) Get(key string) (string, bool) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if item, ok := l.Items[key]; ok {
		l.Queue.MoveToFront(item.KeyPtr)
		return item.Data, true
	}
	return "", false
}

func (l *LRUCache) Remove(key string) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	if item, exists := l.Items[key]; exists {
		delete(l.Items, key)
		l.Queue.Remove(item.KeyPtr)
		return true
	}
	return false
}

func NewLRUCache(n int) LRUCacher {
	return &LRUCache{
		Queue:    list.New(),
		Items:    make(map[string]*Node),
		Capacity: n,
	}
}
