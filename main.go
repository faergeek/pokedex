package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faergeek/pokedex/internal"
)

type (
	cliCommand struct {
		name        string
		description string
		callback    func(config *config) error
	}

	config struct {
		client         *internal.Client
		previous, next string
	}

	locationAreaResponse struct {
		Count          int
		Next, Previous string
		Results        []struct {
			Name, Url string
		}
	}
)

func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(config *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")

	for key, command := range commands {
		fmt.Printf("%s: %s\n", key, command.description)
	}

	return nil
}

func commandMap(config *config) error {
	if config.next == "" {
		return errors.New("Next is empty")
	}

	val, err := config.client.Get(config.next)
	if err != nil {
		return err
	}

	var data locationAreaResponse
	if err := json.Unmarshal(val, &data); err != nil {
		return err
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	config.previous = config.next
	config.next = data.Next

	return nil
}

func commandMapb(config *config) error {
	if config.previous == "" {
		return errors.New("Previous is empty")
	}

	val, err := config.client.Get(config.previous)
	if err != nil {
		return err
	}

	var data locationAreaResponse
	if err := json.Unmarshal(val, &data); err != nil {
		return err
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	config.next = config.previous
	config.previous = data.Previous

	return nil
}

var commands = map[string]cliCommand{}

func main() {
	config := config{
		client:   internal.NewClient(time.Minute),
		previous: "",
		next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
	}

	scanner := bufio.NewScanner(os.Stdin)

	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}

	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	commands["map"] = cliCommand{
		name:        "map",
		description: "Displays the names of the next 20 location areas in the Pokemon world",
		callback:    commandMap,
	}

	commands["mapb"] = cliCommand{
		name:        "map",
		description: "Displays the names of the previous 20 location areas in the Pokemon world",
		callback:    commandMapb,
	}

	for {
		fmt.Print("Pokedex > ")

		if scanner.Scan() {
			words := cleanInput(scanner.Text())

			if len(words) != 0 {
				command, exists := commands[words[0]]

				if !exists {
					fmt.Println("Unknown command")
				} else {
					err := command.callback(&config)
					if err != nil {
						fmt.Printf("error: %v\n", err)
					}
				}
			}
		} else {
			err := scanner.Err()
			if err == nil {
				commandExit(&config)
			} else {
				log.Fatal(err)
			}
		}
	}
}
