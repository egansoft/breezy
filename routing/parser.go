package routing

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/egansoft/breezy/utils"
)

const (
	numRegexpGroups = 6
	reqUrlGroup     = 1
	typeTokenGroup  = 4
	payloadGroup    = 5

	cmdToken = "$"
	fsToken  = ":"
)

var lineRegexp = regexp.MustCompile(`^((\/(\w+|\[\w+\]))+|\/)\/?\s+(\$|:)\s+(.+)$`)

func Parse(lines []string) (*Router, error) {
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
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil
	}

	matches := lineRegexp.FindStringSubmatch(line)
	if len(matches) != numRegexpGroups {
		return syntaxError(line)
	}

	reqUrl := matches[reqUrlGroup]
	typeToken := matches[typeTokenGroup]
	payload := matches[payloadGroup]

	path := utils.UrlToPath(reqUrl)
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
