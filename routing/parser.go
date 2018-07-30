package routing

import (
	"fmt"
	"regexp"

	"github.com/egansoft/silly/utils"
)

const (
	numRegexpGroups = 6
	urlPathGroup    = 1
	typeTokenGroup  = 4
	payloadGroup    = 5

	cmdToken = "$"
	fsToken  = ":"
)

var lineRegexp = regexp.MustCompile(`^((/(\w+|\[\w+\]))+)\s+(\$|:)\s+(.+)$`)

func ParseFile(filepath string) (*Router, error) {
	lines, err := utils.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	r := New()
	for _, line := range lines {
		err := parseAndInsertLine(line, r)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func parseAndInsertLine(line string, r *Router) error {
	matches := lineRegexp.FindStringSubmatch(line)
	if len(matches) != numRegexpGroups {
		return syntaxError(line)
	}

	urlPath := matches[urlPathGroup]
	typeToken := matches[typeTokenGroup]
	payload := matches[payloadGroup]

	path := utils.UrlToPath(urlPath)
	if typeToken == cmdToken {
		return r.InsertCmd(path, payload)
	} else if typeToken == fsToken {
		return r.InsertFs(path, payload)
	}

	return syntaxError(line)
}

func syntaxError(line string) error {
	return fmt.Errorf("Invalid route: %s", line)
}
