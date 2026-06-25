package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if scanner.Scan() {
			text := scanner.Text()

			words := cleanInput(text)

			fmt.Printf("Your command was: %s\n", words[0])
		} else {
			err := scanner.Err()

			if err == nil {
				break
			} else {
				log.Fatal(err)
			}
		}
	}
}
