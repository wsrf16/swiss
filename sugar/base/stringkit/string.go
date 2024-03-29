package stringkit

import (
	// "ios"
	"strings"
	"unicode/utf8"
)

func CountInString(s string) int {
	return utf8.RuneCountInString(s)
}

func ToRune(s string) []rune {
	r := []rune(s)
	return r
}

func HasAnyPrefix(s string, prefixes ...string) (bool, string) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true, prefix
		}
	}
	return false, ""
}

func HasAnySuffix(s string, suffixes ...string) (bool, string) {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true, suffix
		}
	}
	return false, ""
}

func TrimAllPrefixes(s string, prefixes ...string) string {
	ret := s
	for {
		if has, prefix := HasAnyPrefix(ret, prefixes...); has {
			ret = strings.TrimPrefix(ret, prefix)
		} else {
			return ret
		}
	}
}

func TrimAllSuffixes(s string, suffixes ...string) string {
	ret := s
	for {
		if has, suffix := HasAnySuffix(ret, suffixes...); has {
			ret = strings.TrimSuffix(ret, suffix)
		} else {
			return ret
		}
	}
}

func TrimAll(s string, sides ...string) string {
	ret := s
	ret = TrimAllSuffixes(ret, sides...)
	ret = TrimAllPrefixes(ret, sides...)
	return ret
}

func SplitPath(s string, sep string) []string {
	splits := strings.Split(s, sep)
	for i, segment := range splits {
		splits[i] = TrimAll(segment, " ")
	}
	return splits
}

func JoinURL(parts ...string) string {
	url := ""
	for i, part := range parts {
		if i == 0 {
			url += TrimAllSuffixes(part, "/")
		} else {
			url += "/" + TrimAll(part, "/")
		}

	}
	return url
}

// func JoinURL(url ...string) string {
// 	base, next := "", ""
// 	for _, v := range url {
// 		next = TrimAll(v, "/")
// 		base = fmt.Sprintf("%s/%s", base, next)
// 	}
// 	return base
// }

// func RuneCountInString(s string) int {return 1}

// func Base64Bytes(p []byte, w ios.Writer) (int, error) {
//	encoder := base64.NewEncoder(base64.StdEncoding, w)
//	return encoder.Write(p)
// }

// func Base64StringTo(s string, w ios.Writer) (int, error) {
//	input := []byte(s)
//	return Base64Bytes(input, w)
// }
