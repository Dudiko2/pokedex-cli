package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dudiko2/pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2/"

type Client struct {
	cache pokecache.Cache
}

func NewClient() *Client {
	c := Client{
		cache: *pokecache.NewCache(5 * time.Minute),
	}
	return &c
}

func getData(urlString string) (body []byte, err error) {
	res, err := http.Get(urlString)
	if err != nil {
		return body, err
	}
	defer res.Body.Close()
	// XXX add custom errors
	if res.StatusCode > 299 {
		errMsg := fmt.Sprintf("request failed (code %v)", res.StatusCode)
		return body, errors.New(errMsg)
	}
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return body, readErr
	}
	return body, nil
}

func (c *Client) cachedGetData(key string) ([]byte, error) {
	d, found := c.cache.Get(key)
	if found {
		return d, nil
	}
	d, err := getData(key)
	if err != nil {
		return []byte{}, err
	}
	c.cache.Add(key, d)
	return d, nil
}

func parseJSON[T any](data []byte) (T, error) {
	var parsed T
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return parsed, err
	}
	return parsed, nil
}
