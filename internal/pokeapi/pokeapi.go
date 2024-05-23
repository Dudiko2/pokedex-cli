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

func (c *Client) GetLocationAreas(p GetLocationAreasPayload) (LocationAreasRes,
	error) {
	path := fmt.Sprintf("location-area?offset=%v&limit=%v", p.Offset, p.Limit)
	key := baseURL + path
	body, err := c.cachedGetData(key)
	if err != nil {
		return LocationAreasRes{}, err
	}
	locations, err := parseJSON[LocationAreasRes](body)
	if err != nil {
		return LocationAreasRes{}, err
	}
	return locations, nil
}

func (c *Client) GetLocationArea(id string) (LocationRes, error) {
	path := fmt.Sprintf("location-area/%s", id)
	key := baseURL + path
	body, err := c.cachedGetData(key)
	if err != nil {
		return LocationRes{}, err
	}
	location, err := parseJSON[LocationRes](body)
	if err != nil {
		return LocationRes{}, err
	}
	return location, nil
}

type GetLocationAreasPayload struct {
	Offset int
	Limit  int
}

func parseJSON[T any](data []byte) (T, error) {
	var parsed T
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return parsed, err
	}
	return parsed, nil
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

type LocationRes struct {
	GameIndex int             `json:"game_index"`
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Names     []LocationNames `json:"names"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	PokemonEncounters    []LocationPokemonEncounter `json:"pokemon_encounters"`
	EncounterMethodRates []EncounterMethodRate      `json:"encounter_method_rates"`
}

type EncounterMethodRate struct {
	EncounterMethod struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"encounter_method"`
	VersionDetails []struct {
		Rate    int     `json:"rate"`
		Version Version `json:"version"`
	} `json:"version_details"`
}

type LocationPokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
	VersionDetails []PokemonEncounterVersionDetails `json:"version_details"`
}

type LocationNames struct {
	Language LanguageInfo `json:"language"`
	Name     string       `json:"name"`
}

type LanguageInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonEncounterVersionDetails struct {
	EncounterDetails []struct {
		Chance          int   `json:"chance"`
		ConditionValues []any `json:"condition_values"`
		MaxLevel        int   `json:"max_level"`
		Method          struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"method"`
		MinLevel int `json:"min_level"`
	} `json:"encounter_details"`
	MaxChance int     `json:"max_chance"`
	Version   Version `json:"version"`
}

type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
