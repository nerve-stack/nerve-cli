package cases

import (
	"regexp"
	"strings"
)

func ToCamelCase(s string) string {
	// Define exceptions map for special cases.
	exceptions := map[string]string{
		"id":   "ID",
		"uuid": "UUID",
		"url":  "URL",
	}

	// Convert the string to lower case and split into words.
	words := splitAndClean(s)

	// Iterate through the words and convert them to CamelCase.
	for i, word := range words {
		// If the word is in the exceptions map, replace it with the correct value.
		if replacement, found := exceptions[strings.ToLower(word)]; found {
			words[i] = replacement
		} else {
			// Capitalize the first letter of each word, lowercase the rest.
			words[i] = Capitalize(word)
		}
	}

	// Join the words back together and return the result.
	return strings.Join(words, "")
}

func Capitalize(s string) string {
	if s == "" {
		return ""
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

// splitAndClean splits the input string by non-alphabetical characters and cleans it.
func splitAndClean(s string) []string {
	// Remove leading and trailing spaces, then split the string by non-alphabetic characters.
	reg := regexp.MustCompile("[^a-zA-Z]+")
	words := reg.Split(s, -1)

	// Filter out any empty strings caused by multiple non-alphabetic characters.
	var cleanedWords []string

	for _, word := range words {
		if word != "" {
			cleanedWords = append(cleanedWords, word)
		}
	}

	return cleanedWords
}
