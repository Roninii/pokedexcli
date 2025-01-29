package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	pokeapi "github.com/roninii/pokedexcli/PokeAPI"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     string
	Previous string
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Close the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Show available commands",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Show a paginated list of map locations; subsequent calls will show the next page of results.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show the previous page of map locations.",
			callback:    commandMapb,
		},
	}
}

func main() {
	config := &Config{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanInput := cleanInput(input)
		firstWord := cleanInput[0]

		command, ok := commands[firstWord]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(config)
		if err != nil {
			fmt.Printf("Error executing command: %s; %v\n", command.name, err)
		}

	}
}

func cleanInput(input string) []string {
	if input == "" {
		return []string{}
	}
	return strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Available commands:")
	fmt.Println("")

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func commandMap(config *Config) error {
	var url string
	if config.Next != "" {
		url = config.Next
	} else {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error fetching map data: %v", err)
	}

	defer res.Body.Close()
	var mapData pokeapi.Response
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&mapData)
	if err != nil {
		return fmt.Errorf("Error decoding map data: %v", err)
	}

	config.Next = mapData.Next
	if mapData.Previous != nil {
		config.Previous = *mapData.Previous
	}
	for _, location := range mapData.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(config *Config) error {
	if config.Previous == "" {
		return fmt.Errorf("Already at the beginning of the map!")
	}

	res, err := http.Get(config.Previous)
	if err != nil {
		return fmt.Errorf("Error fetching map data: %v", err)
	}

	defer res.Body.Close()
	var mapData pokeapi.Response
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&mapData)
	if err != nil {
		return fmt.Errorf("Error decoding map data: %v", err)
	}

	config.Next = mapData.Next
	if mapData.Previous != nil {
		config.Previous = *mapData.Previous
	} else {
		// if Previous is nil, we are back at the beginning and should clear this out
		config.Previous = ""
	}
	for _, location := range mapData.Results {
		fmt.Println(location.Name)
	}

	return nil
}
