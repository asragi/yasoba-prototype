package utils

func SplitString(s string) []string {
	runes := []rune(s)
	stringArray := make([]string, len(runes))
	for i, r := range runes {
		stringArray[i] = string(r)
	}
	return stringArray
}
