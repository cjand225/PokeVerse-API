package validation

import (
	"regexp"
	"strings"
)

// isValidLanguageCode checks whether the given language code is a valid ISO 639-1 or ISO 639-3 language code.
//
// Parameters:
//   code: A string representing the language code to be validated.
//
// Return Value:
//   bool: A boolean value indicating whether the language code is valid or not. The function returns true
//         if the language code matches either the ISO 639-1 or ISO 639-3 pattern. Otherwise, it returns false.
//
// Example:
//   validCode := isValidLanguageCode("en")    // Returns true
//   invalidCode := isValidLanguageCode("english") // Returns false
func IsValidLanguageCode(code string) bool {
	iso6391Pattern := `^[a-z]{2}$`
	iso6393Pattern := `^[a-z]{3}$`

	code = strings.ToLower(code)

	// Check if the language code matches the ISO 639-1 pattern.
	matched, err := regexp.MatchString(iso6391Pattern, code)
	if err != nil {
		return false
	}

	// If there was no match with the first pattern, try matching with the ISO 639-3 pattern.
	if !matched {
		matched, err = regexp.MatchString(iso6393Pattern, code)
		if err != nil {
			return false
		}
	}

	return matched
}
