package main

import (
	"bufio"
	"fmt"
	"os"

	pokecmd "github.com/roninii/pokedexcli/internal/commands"
)

func main() {
	config := &pokecmd.Config{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanInput := pokecmd.CleanInput(input)
		firstWord := cleanInput[0]

		command, ok := pokecmd.Commands[firstWord]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.Callback(config, cleanInput[1:])
		if err != nil {
			fmt.Printf("Error executing command: %s; %v\n", command.Name, err)
		}
	}
}
