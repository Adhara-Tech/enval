package manifestchecker

import (
	"fmt"
	"regexp"
)

type RegexpVersionParser struct {
	regexpPattern string
	keys          []string
}

func NewRegexVersionParser(regexpPattern string, keys []string) *RegexpVersionParser {
	return &RegexpVersionParser{
		regexpPattern: regexpPattern,
		keys:          keys,
	}
}

func (parser RegexpVersionParser) isKeyPresent(key string) bool {
	for _, currentKey := range parser.keys {
		if key == currentKey {
			return true
		}
	}

	return false
}

func (parser *RegexpVersionParser) Parse(rawVersion string) (map[string]string, error) {
	pattern := regexp.MustCompile(parser.regexpPattern)

	match := pattern.FindStringSubmatch(rawVersion)
	resultMap := make(map[string]string)

	for counter, groupName := range pattern.SubexpNames() {
		if counter != 0 && parser.isKeyPresent(groupName) {
			resultMap[groupName] = match[counter]
		}
		//if counter != 0 && groupName == "version" {
		//	return match[counter], nil
		//}
	}

	for _, currentKey := range parser.keys {
		if _, ok := resultMap[currentKey]; !ok {
			return nil, fmt.Errorf("key [%s] not found as part of regexp to find version fields on version string [%s]", currentKey, rawVersion)
		}

	}

	return resultMap, nil

}
