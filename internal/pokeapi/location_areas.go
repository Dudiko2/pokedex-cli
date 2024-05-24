package pokeapi

import "fmt"

type GetLocationAreasPayload struct {
	Offset int
	Limit  int
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
