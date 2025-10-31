package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {

	tokenPatterns := map[string]*regexp.Regexp{
		"Constant":         regexp.MustCompile(`\b[0-9]+\b`),
		"int":              regexp.MustCompile(`\bint\b`),
		"void":             regexp.MustCompile(`\bvoid\b`),
		"return":           regexp.MustCompile(`\breturn\b`),
		"openParenthesis":  regexp.MustCompile(`\(`),
		"closeParenthesis": regexp.MustCompile(`\)`),
		"openBrace":        regexp.MustCompile(`\{`),
		"closeBrace":       regexp.MustCompile(`\}`),
		"semicolon":        regexp.MustCompile(`;`),
	}

	data, err := os.ReadFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	text := string(data)

	for name, pattern := range tokenPatterns {
		matches := pattern.FindAllString(text, -1)
		if len(matches) > 0 {
			fmt.Printf("Token: %-15s  Matches: %v\n", name, matches)
		}
	}

}
