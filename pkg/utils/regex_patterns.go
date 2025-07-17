package utils

import "regexp"

var URLPattern = `^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`

func IsValidURL(url string) bool {
	re := regexp.MustCompile(URLPattern)
	return re.MatchString(url)
}
