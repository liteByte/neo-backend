package templater

import (
	"strings"
	"io/ioutil"
)

func Template(file string, values map[string]string) (string, error) {

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	result := string(b)

	state := 0
	key := ""
	for _, char := range b {
		switch state {
		case 0:
			if char == '{' {
				state = 1
			}
			break
		case 1:
			if char == '{' {
				state = 2
			} else {
				state = 0
			}
			break
		case 2:
			if char == '}' {
				state = 3
			} else {
				key += string(char)
			}
			break
		case 3:
			if char == '}' {
				if v, p := values[key]; p {
					r := strings.NewReplacer("{{" + key + "}}", v)
					result = r.Replace(result)
				}
			}
			state = 0
			key = ""
			break
		}
	}

	return result, nil
}
