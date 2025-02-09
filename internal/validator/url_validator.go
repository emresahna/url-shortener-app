package validator

import "regexp"

var urlRegex = regexp.MustCompile(`^https?:\/\/([\w-]+\.)+[a-zA-Z]{2,63}(:\d{1,5})?(\/[^\s]*)?$`)

func ValidateURL(url string) bool {
	return urlRegex.MatchString(url)
}
