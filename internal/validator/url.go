package validator

import "regexp"

func ValidateURL(url string) bool {
	regex := `^(https?:\/\/)?[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(url)
}
