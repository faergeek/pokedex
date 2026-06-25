package internal

import (
	"io"
	"net/http"
	"time"
)

type Client struct {
	cache *Cache
}

func NewClient(lifetime time.Duration) *Client {
	client := Client{
		cache: NewCache(lifetime),
	}

	return &client
}

func (c *Client) Get(url string) ([]byte, error) {
	if val, cached := c.cache.Get(url); cached {
		return val, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	val, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	c.cache.Add(url, val)
	return val, nil
}
