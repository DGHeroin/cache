package lru

import (
    "container/list"
    "log"
    "sync"
)

const (
    RemoveTypeFullEntries = RemoveReason(0)
    RemoveTypeByUser      = RemoveReason(1)
)

type Cache struct {
    MaxEntries uint64
    OnEvicted  func(key string, value interface{}, fullType RemoveReason)
    ll         *list.List
    cache      map[interface{}]*list.Element
    mu         sync.RWMutex
}

type RemoveReason int

type entry struct {
    key   string
    value interface{}
}

func New(maxEntries uint64) *Cache {
    if maxEntries == 0 {
        maxEntries = ^uint64(0)
    }
    s := &Cache{
        MaxEntries: maxEntries,
        ll:         list.New(),
        cache:      make(map[interface{}]*list.Element),
    }
    log.Println("=============", s.ll)
    return s
}

func (c *Cache) Add(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.cache == nil {
        c.cache = make(map[interface{}]*list.Element)
        c.ll = list.New()
    }

    if ee, ok := c.cache[key]; ok {
        c.ll.MoveToFront(ee)
        ee.Value.(*entry).value = value
        return
    }
    ele := c.ll.PushFront(&entry{key, value})
    c.cache[key] = ele
    if uint64(c.ll.Len()) > c.MaxEntries {
        c.removeOldestLocked(RemoveTypeFullEntries)
    }
}

func (c *Cache) Get(key string) (value interface{}, ok bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.cache == nil {
        return
    }
    if ele, hit := c.cache[key]; hit {
        c.ll.MoveToFront(ele)
        return ele.Value.(*entry).value, true
    }
    return
}

func (c *Cache) Remove(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.cache == nil {
        return
    }
    if ele, hit := c.cache[key]; hit {
        c.removeElement(ele, RemoveTypeByUser)
    }
}

func (c *Cache) RemoveOldest() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.removeOldestLocked(RemoveTypeByUser)
}

func (c *Cache) removeOldestLocked(removeType RemoveReason) {
    if c.cache == nil {
        return
    }
    ele := c.ll.Back()
    if ele != nil {
        c.removeElement(ele, removeType)
    }
}

func (c *Cache) removeElement(e *list.Element, removeType RemoveReason) {
    c.ll.Remove(e)
    kv := e.Value.(*entry)
    delete(c.cache, kv.key)
    if c.OnEvicted != nil {
        c.OnEvicted(kv.key, kv.value, removeType)
    }
}

func (c *Cache) Len() uint64 {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.lenLocked()
}

func (c *Cache) lenLocked() uint64 {
    if c.cache == nil {
        return 0
    }
    return uint64(c.ll.Len())
}

func (c *Cache) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.OnEvicted != nil {
        for _, e := range c.cache {
            kv := e.Value.(*entry)
            c.OnEvicted(kv.key, kv.value, RemoveTypeFullEntries)
        }
    }
    c.ll = list.New()
    c.cache = make(map[interface{}]*list.Element)
}
func (c *Cache) Range(fn func(key string, value interface{})) {
    c.mu.Lock()
    defer c.mu.Unlock()
    log.Println(">>>",c.ll)
    if c.ll.Len() == 0 {
        return
    }
    e := c.ll.Front()
    for e != nil {
        kv := e.Value.(*entry)
        fn(kv.key, kv.value)
        e = e.Next()
    }
}

func (r RemoveReason) String() string {
    switch r {
    case RemoveTypeFullEntries:
        return "Remove by full entries"
    case RemoveTypeByUser:
        return "Remove by user"
    }
    return "Unknown remove reason"
}
