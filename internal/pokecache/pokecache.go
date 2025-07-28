package pokecache

import (
    "sync"
    "time"
)

type Cache struct {
    entries map[string]cacheEntry
    mutex *sync.Mutex
    interval time.Duration
}

type cacheEntry struct {
    createdAt time.Time
    val []byte
}

func NewCache(interval time.Duration) Cache {
    c := Cache{
        entries: make(map[string]cacheEntry),
        mutex: &sync.Mutex{},
        interval: interval,
    }
    go c.reapLoop()
    return c
}

func (c Cache) Add(key string, val []byte) {
    c.mutex.Lock()
    c.entries[key] = cacheEntry{
        createdAt: time.Now(),
        val: val,
    }
    c.mutex.Unlock()
}

func (c Cache) Get(key string) ([]byte, bool) {
    c.mutex.Lock()
    val, exists := c.entries[key]
    c.mutex.Unlock()
    return val.val, exists
}

func (c Cache) reapLoop() {
    ticker := time.NewTicker(c.interval)
    defer ticker.Stop()
    for {
        select {
        case t := <-ticker.C:
            c.mutex.Lock()
            for key, val := range c.entries {
                if val.createdAt.Add(c.interval).Before(t) {
                    delete(c.entries, key)
                }
            }
            c.mutex.Unlock()
        }
    }
}
