package pokedex

import (
	"github.com/roninii/pokedexcli/internal/pokeapi"
)

type Pokemon = pokeapi.Pokemon

var Pokedex = map[string]Pokemon{}

func AddPokemon(p Pokemon) {
	Pokedex[p.Name] = p
}
