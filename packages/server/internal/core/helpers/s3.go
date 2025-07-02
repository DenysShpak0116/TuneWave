package helpers

import "strings"

func ExtractS3Key(fullURL string) string {
	const baseURL = "https://tunewavebucket.s3.eu-west-3.amazonaws.com/"
	return strings.TrimPrefix(fullURL, baseURL)
}
