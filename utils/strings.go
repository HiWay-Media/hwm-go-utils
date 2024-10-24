package utils

import (
	"regexp"
	"strings"
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