/*
 * Copyright (C) 2024 STD DEEPANSHU <https://deepanshu.in>
 * STD BOTS - Telegram: @STD_DEEPANSHU, @STDBOTS
 *
 * This file is part of FeedVackBot.
 *
 * FeedVackBot is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, version 3 of the License.
 */

package database

import (
	"sync"
	"time"
)

type cacheItem struct {
	value      interface{}
	expiration int64
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

func NewCache() *Cache {
	c := &Cache{
		items: make(map[string]cacheItem),
	}
	go c.startCleanup()
	return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}
	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		return nil, false
	}
	return item.value, true
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:      value,
		expiration: expiration,
	}
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		now := time.Now().UnixNano()
		c.mu.Lock()
		for k, v := range c.items {
			if v.expiration > 0 && now > v.expiration {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}
