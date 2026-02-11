package external

import (
	"sync"
	"time"
)

// CacheItem representa un item en el cache
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// Cache implementa un cache en memoria simple con TTL
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewCache crea una nueva instancia de Cache
func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		items: make(map[string]CacheItem),
		ttl:   ttl,
	}

	// Iniciar limpieza periódica de items expirados
	go cache.cleanup()

	return cache
}

// Set guarda un valor en el cache
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(c.ttl),
	}
}

// Get obtiene un valor del cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Verificar si expiró
	if time.Now().After(item.Expiration) {
		return nil, false
	}

	return item.Value, true
}

// Delete elimina un valor del cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear limpia todo el cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]CacheItem)
}

// Size retorna el número de items en el cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// cleanup elimina periódicamente items expirados
func (c *Cache) cleanup() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for range ticker.C {
		c.removeExpired()
	}
}

// removeExpired elimina items expirados del cache
func (c *Cache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.Expiration) {
			delete(c.items, key)
		}
	}
}
