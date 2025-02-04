package pokeapi

const (
	BaseURL         = "https://pokeapi.co/api/v2"
	LocationAreaURL = BaseURL + "/location-area/"
	PokemonURL      = BaseURL + "/pokemon/"
)

type Response struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous *string   `json:"previous"`
	Results  []Results `json:"results"`
}

type Results struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ExploreResponse struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonEncounters struct {
	Pokemon PokemonEncounter `json:"pokemon"`
}

type PokemonResponse struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}
