package generate

import (
	"strings"
	"unicode"
)

// goKeywords contains Go reserved words that need escaping in identifiers.
var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
}

// ToPascalCase converts a snake_case or camelCase string to PascalCase.
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	// Split on underscores, hyphens, and camelCase boundaries
	words := splitWords(s)
	var result strings.Builder
	for _, word := range words {
		if word == "" {
			continue
		}
		// Handle common abbreviations
		upper := strings.ToUpper(word)
		if isCommonAbbreviation(upper) {
			result.WriteString(upper)
		} else {
			result.WriteString(strings.ToUpper(word[:1]))
			result.WriteString(strings.ToLower(word[1:]))
		}
	}
	return result.String()
}

// ToCamelCase converts a snake_case string to camelCase.
func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if pascal == "" {
		return ""
	}

	// Find the boundary of the first "word" to lowercase
	// For abbreviations like "ID", lowercase the whole thing
	runes := []rune(pascal)
	if len(runes) <= 1 {
		return strings.ToLower(pascal)
	}

	// If starts with consecutive uppercase (abbreviation), lowercase all consecutive uppercase
	// except the last one if it's followed by a lowercase letter.
	i := 0
	for i < len(runes) && unicode.IsUpper(runes[i]) {
		i++
	}
	if i > 1 {
		// Multiple consecutive uppercase letters
		lowEnd := i
		if i < len(runes) {
			// If there's a following lowercase letter, keep last uppercase as start of next word
			lowEnd = i - 1
		}
		for j := 0; j < lowEnd; j++ {
			runes[j] = unicode.ToLower(runes[j])
		}
	} else {
		runes[0] = unicode.ToLower(runes[0])
	}

	return string(runes)
}

// ToPackageName converts a contract name to a valid Go package name.
func ToPackageName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, " ", "")

	// Package names must start with a letter
	if len(s) > 0 && !unicode.IsLetter(rune(s[0])) {
		s = "pkg" + s
	}
	return s
}

// SafeGoName ensures an identifier is not a Go keyword.
func SafeGoName(name string) string {
	if goKeywords[name] {
		return name + "Val"
	}
	return name
}

// splitWords splits a string into words by underscore, hyphen, space, and camelCase boundaries.
func splitWords(s string) []string {
	var words []string
	var current strings.Builder

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == '_' || r == '-' || r == ' ' {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
			continue
		}

		// CamelCase boundary: lowercase followed by uppercase
		if i > 0 && unicode.IsUpper(r) && current.Len() > 0 {
			prevR := runes[i-1]
			if unicode.IsLower(prevR) {
				words = append(words, current.String())
				current.Reset()
			} else if unicode.IsUpper(prevR) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				// e.g. "XMLParser" -> "XML", "Parser"
				words = append(words, current.String())
				current.Reset()
			}
		}

		current.WriteRune(r)
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}

// isCommonAbbreviation returns true for common Go abbreviations that should be all-caps.
func isCommonAbbreviation(s string) bool {
	abbreviations := map[string]bool{
		"ID": true, "URL": true, "URI": true, "API": true, "HTTP": true,
		"HTTPS": true, "JSON": true, "XML": true, "SQL": true, "HTML": true,
		"CSS": true, "IP": true, "TCP": true, "UDP": true, "TLS": true,
		"SSL": true, "SSH": true, "RPC": true, "ABI": true, "SDK": true,
		"ARC": true, "MBR": true, "ASA": true, "NFT": true, "KV": true,
		"TX": true, "TXN": true,
	}
	return abbreviations[s]
}
