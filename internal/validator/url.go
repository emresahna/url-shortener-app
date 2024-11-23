package validator

import "regexp"

func ValidateURL(url string) bool {
	regex := `^(https?|ftp):\/\/[^\s/$.?#].[^\s]*$`
	re := regexp.MustCompile(regex)
	return re.MatchString(url)
}
