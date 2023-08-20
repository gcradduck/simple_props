package simple_props

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Props struct captures a map of key strings to any values.
type Props struct {
	Props map[string]any
}

// Get returns string property value mapped to the provided string.
func (p *Props) Get(key string) string {
	v := p.Props[key]

	if v == nil {
		return ""
	}

	return (p.Props[key]).(string)
}

// GetCleanFilePath returns a "cleaned" filepath, ensuring system file separator is never duplicated.
func (p *Props) GetCleanFilePath(key string) string {
	filePath := p.Get(key)

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return ""
	}

	doubleSeparator := string(os.PathSeparator) + string(os.PathSeparator)

	present := strings.Contains(filePath, doubleSeparator)

	for present {

		filePath = strings.ReplaceAll(filePath, doubleSeparator, string(os.PathSeparator))
		present = strings.Contains(filePath, doubleSeparator)

	}

	return filePath
}

// GetInt returns int property value mapped to the provided string
func (p *Props) GetInt(key string, defaultValue int) int {
	x, err := strconv.Atoi(p.Get(key))

	if err != nil {
		return defaultValue
	}

	return x
}

// LoadProps loads a Props struct by way of a provided filepath/name.
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

	regParser, regError := regexp.Compile(`(?sm)^\s*([^\n]*)\s*=\s*([^\n]*)\s*$`)

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

		props.Props[key] = strings.TrimSpace(value)

	}

	return props, nil

}
