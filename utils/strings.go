package utils

import (
	"regexp"
	"strings"
	"encoding/json"
	"net/url"
    "fmt"
)

/* 
    utils strings metods
*/



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
	re := regexp.MustCompile(`/+`)
	cleanedPath := re.ReplaceAllString(path, "/")

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

func RemoveSlice(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// print struct json format to console
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}



func EncodeURL(rawUrl string) string {
	//return url.QueryEscape(s)
	// Parse the URL
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return ""
	}

	// Get the base URL (without query params)
	baseUrl := parsedUrl.Scheme + "://" + parsedUrl.Host + parsedUrl.Path

	// Get the query parameters and encode them
	queryParams := parsedUrl.Query()
	encodedParams := url.Values{}

	// Loop through and encode each query parameter
	for key, values := range queryParams {
		for _, value := range values {
			encodedParams.Add(key, url.QueryEscape(value))
		}
	}

	// Combine the base URL with the encoded query parameters
	finalUrl := baseUrl + "?" + encodedParams.Encode()

	// Output the transformed URL
	//fmt.Println("Encoded URL:", finalUrl)
	return finalUrl
}
