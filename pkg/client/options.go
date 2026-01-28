package client

import (
	"context"
	"time"
)

type UaaOption func(c *UaaClient)

func WithTimeout(d time.Duration) func(c *UaaClient) {
	return func(c *UaaClient) {
		c.timeout = d
	}
}

func WithContext(ctx context.Context) func(*UaaClient) {
	return func(c *UaaClient) {
		c.ctx = ctx
	}
}

func WithStore(store KeyStore) func(*UaaClient) {
	return func(c *UaaClient) {
		c.s = store
	}
}
