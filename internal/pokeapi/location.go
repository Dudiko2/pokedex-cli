package pokeapi

import "fmt"

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
