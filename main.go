package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

type (
	cliCommand struct {
		name        string
		description string
		callback    func(config *config) error
	}

	config struct {
		Previous string
		Next     string
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
	if config.Next == "" {
		return errors.New("Next is empty")
	}

	res, err := http.Get(config.Next)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	var data locationAreaResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	config.Next = data.Next
	config.Previous = data.Previous

	return nil
}

func commandMapb(config *config) error {
	if config.Previous == "" {
		return errors.New("Previous is empty")
	}

	res, err := http.Get(config.Previous)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	var data locationAreaResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	config.Next = data.Next
	config.Previous = data.Previous

	return nil
}

var commands = map[string]cliCommand{}

func main() {
	config := config{
		Previous: "",
		Next:     "https://pokeapi.co/api/v2/location-area/?limit=20&offset=0",
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
