package commands

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/roninii/pokedexcli/internal/pokeapi"
	"github.com/roninii/pokedexcli/internal/pokecache"
	"github.com/roninii/pokedexcli/internal/pokedex"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Config, []string) error
}

type Config struct {
	Next     string
	Previous string
}

var Commands map[string]CliCommand
var cache pokecache.Cache

func init() {
	cache = pokecache.NewCache(5 * time.Second)
	Commands = map[string]CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Close the Pokedex",
			Callback:    CommandExit,
		},
		"help": {
			Name:        "help",
			Description: "Show available commands",
			Callback:    CommandHelp,
		},
		"map": {
			Name:        "map",
			Description: "Show a paginated list of map locations; subsequent calls will show the next page of results.",
			Callback:    CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Show the previous page of map locations.",
			Callback:    CommandMapb,
		},
		"explore": {
			Name:        "explore",
			Description: "Show a list of Pokemon in a given location.",
			Callback:    CommandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Attempt to catch the specified Pokemon.",
			Callback:    CommandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect a caught Pokemon.",
			Callback:    CommandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "List all caught Pokemon.",
			Callback:    CommandPokedex,
		},
	}
}

func CleanInput(input string) []string {
	if input == "" {
		return []string{}
	}
	return strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
}

func CommandExit(config *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func CommandHelp(config *Config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Available commands:")
	fmt.Println("")

	for _, command := range Commands {
		fmt.Printf("%s: %s\n", command.Name, command.Description)
	}

	return nil
}

func CommandMap(config *Config, args []string) error {
	var url string
	if config.Next != "" {
		url = config.Next
	} else {
		url = pokeapi.LocationAreaURL
	}

	var mapData pokeapi.Response

	if val, exists := cache.Get(url); exists {
		err := json.Unmarshal(val, &mapData)
		if err != nil {
			return fmt.Errorf("Error decoding map data: %v", err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error fetching map data: %v", err)
		}

		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&mapData)
		if err != nil {
			return fmt.Errorf("Error decoding map data: %v", err)
		}

		responseBytes, err := json.Marshal(mapData)
		if err != nil {
			fmt.Printf("Error adding response to cache: %v\n", err)
		} else {
			cache.Add(url, responseBytes)
		}
	}

	config.Next = mapData.Next
	if mapData.Previous != nil {
		config.Previous = *mapData.Previous
	}

	printEntries(mapData.Results)

	return nil
}

func CommandMapb(config *Config, args []string) error {
	if config.Previous == "" {
		return fmt.Errorf("Already at the beginning of the map!")
	}

	var mapData pokeapi.Response

	if val, exists := cache.Get(config.Previous); exists {
		err := json.Unmarshal(val, &mapData)
		if err != nil {
			return fmt.Errorf("Error decoding map data: %v", err)
		}
	} else {
		res, err := http.Get(config.Previous)
		if err != nil {
			return fmt.Errorf("Error fetching map data: %v", err)
		}

		defer res.Body.Close()
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&mapData)
		if err != nil {
			return fmt.Errorf("Error decoding map data: %v", err)
		}

		responseBytes, err := json.Marshal(mapData)
		if err != nil {
			fmt.Printf("Error adding response to cache: %v\n", err)
		} else {
			cache.Add(config.Previous, responseBytes)
		}

	}

	config.Next = mapData.Next
	if mapData.Previous != nil {
		config.Previous = *mapData.Previous
	} else {
		// if Previous is nil, we are back at the beginning and should clear this out
		config.Previous = ""
	}

	printEntries(mapData.Results)

	return nil
}

func CommandExplore(config *Config, args []string) error {
	location := args[0]
	url := fmt.Sprintf("%s%s", pokeapi.LocationAreaURL, location)

	var areaData pokeapi.ExploreResponse

	if val, exists := cache.Get(url); exists {
		err := json.Unmarshal(val, &areaData)
		if err != nil {
			return fmt.Errorf("Error decoding Pokemon data at location %s: %v", location, err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error fetching Pokemon data at location %s: %v", location, err)
		}

		defer res.Body.Close()
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&areaData)
		if err != nil {
			return fmt.Errorf("Error decoding Pokemon data at location %s: %v", location, err)
		}

		responseBytes, err := json.Marshal(areaData)
		if err != nil {
			fmt.Printf("Error adding response to cache: %v\n", err)
		} else {
			cache.Add(config.Previous, responseBytes)
		}
	}

	fmt.Println("")
	for _, encounter := range areaData.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}

func CommandCatch(config *Config, args []string) error {
	pokemon := args[0]
	url := fmt.Sprintf("%s%s", pokeapi.PokemonURL, pokemon)
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	var pokemonData pokeapi.Pokemon
	if val, exists := cache.Get(url); exists {
		err := json.Unmarshal(val, &pokemonData)
		if err != nil {
			return fmt.Errorf("Error decoding data for %s: %v", pokemon, err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error fetching Pokemon data for %s: %v", pokemon, err)
		}

		defer res.Body.Close()
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&pokemonData)
		if err != nil {
			return fmt.Errorf("Error decoding Pokemon data for %s: %v", pokemon, err)
		}

		responseBytes, err := json.Marshal(pokemonData)
		if err != nil {
			fmt.Printf("Error adding response to cache: %v\n", err)
		} else {
			cache.Add(config.Previous, responseBytes)
		}
	}

	baseCatchRate := math.Max(10, float64(100-pokemonData.BaseExperience))
	roll := rand.Float64() * 100

	if roll <= baseCatchRate {
		fmt.Printf("%s was caught!\n", pokemon)
		fmt.Printf("Adding %s to the Pokedex...\n", pokemon)
		fmt.Printf("Done! You may now view details about %s with the inspect command.\n", pokemon)
		pokedex.AddPokemon(pokemonData)
	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}

	return nil
}

func CommandInspect(config *Config, args []string) error {
	pokedex := pokedex.Pokedex
	name := args[0]
	if pokemon, ok := pokedex[name]; ok {
		// TODO: loop through values in the name, height, wieght, stats, and types and print them out
		for _, key := range []string{"Name", "Height", "Weight", "Stats", "Types"} {
			switch key {
			case "Name":
				fmt.Printf("Name: %s\n", pokemon.Name)
			case "Height":
				fmt.Printf("Height: %d\n", pokemon.Height)
			case "Weight":
				fmt.Printf("Weight: %d\n", pokemon.Weight)
			case "Stats":
				fmt.Println("Stats:")
				for _, stat := range pokemon.Stats {
					fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
				}
			case "Types":
				fmt.Println("Types:")
				for _, t := range pokemon.Types {
					fmt.Printf("  - %s\n", t.Type.Name)
				}

			}
		}

		return nil
	}

	return fmt.Errorf("%s has not been caught.\n", name)
}

func CommandPokedex(config *Config, args []string) error {
	for name := range pokedex.Pokedex {
		fmt.Printf("  - %s\n", name)
		return nil
	}

	return fmt.Errorf("No Pokemon have been caught yet.")
}

func printEntries(entries []pokeapi.Results) {
	fmt.Println("")
	for _, location := range entries {
		fmt.Println(location.Name)
	}
}
