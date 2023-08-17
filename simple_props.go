package simple_props

import (
	"errors"
	"os"
	"regexp"
	"strconv"
)

type Props struct {
	Props map[string]any
}

func (p *Props) Get(key string) string {
	return (p.Props[key]).(string)
}

func (p *Props) GetInt(key string, defaultValue int) int {
	x, err := strconv.Atoi(p.Get(key))

	if err != nil {
		return defaultValue
	}

	return x
}

func LoadProps(filepath string) (*Props, error) {

	props := &Props{
		Props: make(map[string]any, 1),
	}

	fileBytes, err := os.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	stringData := string(fileBytes)

	regLines, regError := regexp.Compile(`(?m)([^=]*)\s*=\s*(.*)\s*`)

	if regError != nil {
		return nil, regError
	}

	regParser, regError := regexp.Compile(`^\s*([^=]*)\s*=\s*(.*)\s*$`)

	if regError != nil {
		return nil, regError
	}

	matches := regLines.FindAllStringSubmatch(stringData, -1)

	if len(matches) == 0 {
		return nil, errors.New("no properties file lines contained property mappings")
	}

	for _, topMatch := range matches {

		toParse := topMatch[0]

		parseMatches := regParser.FindAllStringSubmatch(toParse, -1)

		key := parseMatches[0][1]
		value := parseMatches[0][2]

		props.Props[key] = value

	}

	return props, nil

}
