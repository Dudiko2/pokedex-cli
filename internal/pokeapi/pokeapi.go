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
	client http.Client
	cache  pokecache.Cache
}

func NewClient() *Client {
	c := Client{
		client: http.Client{},
		cache:  *pokecache.NewCache(5 * time.Minute),
	}
	return &c
}

func getData(urlString string) (body []byte, err error) {
	res, err := http.Get(urlString)
	if err != nil {
		return body, err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		errMsg := fmt.Sprintf("Request failed (code %v)", res.StatusCode)
		return body, errors.New(errMsg)
	}
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return body, readErr
	}
	return body, nil
}

func (c *Client) GetLocationAreas(p GetLocationAreasPayload) (LocationAreasRes,
	error) {
	path := fmt.Sprintf("location-area?offset=%v&limit=%v", p.Offset, p.Limit)
	key := baseURL + path
	body, err := func() ([]byte, error) {
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
	}()
	if err != nil {
		return LocationAreasRes{}, err
	}
	locations, err := parseLocationAreas(body)
	if err != nil {
		return LocationAreasRes{}, err
	}
	return locations, nil
}

func (c *Client) GetLocationArea(id string) {}

type GetLocationAreasPayload struct {
	Offset int
	Limit  int
}

func parseLocationAreas(data []byte) (LocationAreasRes, error) {
	locations := LocationAreasRes{}
	err := json.Unmarshal(data, &locations)
	if err != nil {
		return LocationAreasRes{}, err
	}
	return locations, nil
}

type LocationAreasEntry struct {
	Name string  `json:"name"`
	Url  *string `json:"url"`
}

type LocationAreasRes struct {
	Count    int                  `json:"count"`
	Next     *string              `json:"next"`
	Previous *string              `json:"previous"`
	Results  []LocationAreasEntry `json:"results"`
}
