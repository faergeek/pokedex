package main

import (
	"strings"
	"unicode"
)

func cleanInput(text string) []string {
	result := []string{}

	item := ""
	for _, c := range strings.ToLower(text) {
		if unicode.IsSpace(c) {
			if len(item) != 0 {
				result = append(result, item)
				item = ""
			}
		} else {
			item += string(c)
		}
	}

	if len(item) != 0 {
		result = append(result, item)
	}

	return result
}
