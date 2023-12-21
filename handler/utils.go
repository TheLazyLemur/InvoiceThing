package handler

import (
	"bytes"
	"log/slog"
	"sync"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type cache struct {
	lock  sync.RWMutex
	cache map[string][]byte
}

var che = &cache{
	cache: make(map[string][]byte),
	lock:  sync.RWMutex{},
}

func (c *cache) get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.cache[key], c.cache[key] != nil
}

func (c *cache) set(key string, val []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cache[key] = val
}

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func renderWithCache(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html; charset=utf-8")

	p := c.Path()
	data, ok := che.get(p)

	if ok {
		if _, err := c.Write(data); err != nil {
			return err
		}
		slog.Debug("from cache")
		return nil
	}

	bytes := bytes.NewBuffer(make([]byte, 0))

	err := component.Render(c.Context(), bytes)
	if err != nil {
		return err
	}

	if _, err = c.Write(bytes.Bytes()); err != nil {
		return err
	}

	che.set(p, bytes.Bytes())

	return nil
}
