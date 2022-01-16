package main

var alphabets = map[string]string{
	"binary":    "01",
	"lower":     alphabetLower(),
	"upper":     alphabetUpper(),
	"numbers":   alphabetNumbers(),
	"alnum":     alphabetLower() + alphabetNumbers(),
	"hex":       alphabetNumbers() + "ABCDEF",
	"non-space": alphabetNonSpace(),
}

func alphabetLower() (res string) {
	for c := 'a'; c <= 'z'; c++ {
		res += string(c)
	}

	return res
}

func alphabetUpper() (res string) {
	for c := 'A'; c <= 'Z'; c++ {
		res += string(c)
	}

	return res
}

func alphabetNumbers() (res string) {
	for c := '0'; c <= '9'; c++ {
		res += string(c)
	}

	return res
}

func alphabetNonSpace() (res string) {
	for c := 33; c <= 127; c++ {
		res += string(c)
	}

	return res
}
