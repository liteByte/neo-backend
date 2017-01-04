package templater

import (
	"io/ioutil"
)

func Template(file string, values map[string]string) (string, error) {

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	result := ""

	state := 0
	key := ""
	for _, char := range b {
		result += string(char)
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
					result = result[:len(result) - len("{{" + key + "}}")] + v
				}
			}
			if char == '{' {
				state = 1
			} else {
				state = 0
			}
			key = ""
			break
		}
	}

	return result, nil
}
