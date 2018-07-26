package parser

import (
	"fmt"
	"regexp"

	"github.com/egansoft/silly/tree"
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

func ParseFile(filepath string) (*tree.Tree, error) {
	lines, err := utils.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	t := tree.New()
	for _, line := range lines {
		err := parseAndInsertLine(line, t)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func parseAndInsertLine(line string, t *tree.Tree) error {
	matches := lineRegexp.FindStringSubmatch(line)
	if len(matches) != numRegexpGroups {
		return syntaxError(line)
	}

	urlPath := matches[urlPathGroup]
	typeToken := matches[typeTokenGroup]
	payload := matches[payloadGroup]

	if typeToken == cmdToken {
		path, cmd, err := parseCmd(urlPath, payload)
		if err == nil {
			return t.InsertCmd(path, cmd)
		}
	} else if typeToken == fsToken {
		path, fs, err := parseFs(urlPath, payload)
		if err == nil {
			return t.InsertFs(path, fs)
		}
	}

	return syntaxError(line)
}

func parseCmd(urlPathString string, payload string) ([]string, string, error) {
	urlPath := utils.UrlToPath(urlPathString)
	return urlPath, payload, nil
}

func parseFs(urlPathString string, payload string) ([]string, string, error) {
	urlPath := utils.UrlToPath(urlPathString)
	return urlPath, payload, nil
}

func syntaxError(line string) error {
	return fmt.Errorf("Invalid route: %s", line)
}
