package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")

	for key, command := range commands {
		fmt.Printf("%s: %s\n", key, command.description)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands = map[string]cliCommand{}

func main() {
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

	for {
		fmt.Print("Pokedex > ")

		if scanner.Scan() {
			words := cleanInput(scanner.Text())

			if len(words) != 0 {
				command, exists := commands[words[0]]

				if !exists {
					fmt.Println("Unknown command")
				} else {
					err := command.callback()
					if err != nil {
						fmt.Printf("error: %v\n", err)
					}
				}
			}
		} else {
			err := scanner.Err()
			if err == nil {
				commandExit()
			} else {
				log.Fatal(err)
			}
		}
	}
}
