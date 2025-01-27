package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(input string) []string {
	if input == "" {
		return []string{}
	}
	return strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
}
