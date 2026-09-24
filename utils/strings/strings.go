package strings_utils

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strings"
)

/*
   utils strings metods
*/

var multipleSlashes = regexp.MustCompile(`/+`)

func CleanUrlPath(url string) string {
	// Split the protocol from the rest of the URL
	protocolSplit := strings.SplitN(url, "://", 2)

	// Check if we have a valid split (contains "://")
	if len(protocolSplit) < 2 {
		return url // If there's no "://", return as is (not a valid URL)
	}

	protocol := protocolSplit[0]
	path := protocolSplit[1]

	// Replace multiple slashes in the path part (ignores the protocol)
	cleanedPath := multipleSlashes.ReplaceAllString(path, "/")

	// Recombine the protocol and cleaned path
	return protocol + "://" + cleanedPath
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func ContainsWord(s, word string) bool {
	// Convert both the string and the word to lowercase for case-insensitive comparison
	s = strings.ToLower(s)
	word = strings.ToLower(word)

	// Check if the string contains the word
	return strings.Contains(s, word)
}

// RemoveSlice returns a copy of slice without the element at index s; the
// input is not modified. An out-of-range index returns an unchanged copy.
func RemoveSlice(slice []string, s int) []string {
	out := make([]string, 0, len(slice))
	for i, v := range slice {
		if i != s {
			out = append(out, v)
		}
	}
	return out
}

// print struct json format to console
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// EncodeURL re-encodes the query string of rawUrl so every parameter value is
// escaped exactly once; scheme, userinfo, host, path and fragment are kept.
// It returns "" when rawUrl cannot be parsed.
func EncodeURL(rawUrl string) string {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return ""
	}
	parsedUrl.RawQuery = parsedUrl.Query().Encode()
	return parsedUrl.String()
}

// ReverseString reverses a string and returns it.
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// CountWords returns the number of words in a given string.
func CountWords(s string) int {
	words := strings.Fields(s)
	return len(words)
}
