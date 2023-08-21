package simple_props

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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

// GetBool returns bool property value mapped to the provided string; TRUE values include [TRUE,true,YES,yes,Y,y,1]; all other values represent FALSE.
func (p *Props) GetBool(key string, defaultValue bool) bool {

	stringValue := strings.ToUpper(p.Get(key))

	switch stringValue {
	case "TRUE":
		return true
	case "YES":
		return true
	case "Y":
		return true
	case "1":
		return true
	default:
		return defaultValue
	}

}

// GetDate returns a date value mapped to the provided properties key/string. Format directive is expected; e.g., dateProperty=08/21/2023 format:MM/DD/YYYY
func (p *Props) GetDate(key string) (time.Time, error) {

	stringValue := p.Get(key)

	if stringValue == "" {
		return time.Now(), fmt.Errorf("error reading date property from properties file - no such key %s", key)
	}

	dateReg, regErr := regexp.Compile(`(?i)^(.*?)\s*format:(.*)$`)

	if regErr != nil {
		return time.Now(), regErr
	}

	matches := dateReg.FindAllStringSubmatch(stringValue, -1)

	if len(matches) == 0 || len(matches[0]) < 3 {
		return time.Now(), fmt.Errorf("attempted to match improperly formed date property for key: %s", key)
	}

	dateVal := strings.TrimSpace(matches[0][1])
	dateFormat := strings.TrimSpace(matches[0][2])

	switch dateFormat {
	case "YYYY-MM-DD":
		return time.Parse("2006-01-02", dateVal)
	case "YYYY-M-D":
		return time.Parse("2006-1-2", dateVal)
	case "YYYYMMDD":
		return time.Parse("20060102", dateVal)
	case "MM/DD/YYYY":
		return time.Parse("01/02/2006", dateVal)
	case "M/D/YYYY":
		return time.Parse("1/2/2006", dateVal)

	default:
		return time.Now(), fmt.Errorf("failed to match date format %s for key %s", dateFormat, key)
	}

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
