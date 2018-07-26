package utils

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var varTokenRegexp = regexp.MustCompile(`\[\w+\]`)

func ReadFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func UrlToPath(url string) []string {
	url = strings.Trim(url, "/")
	return strings.Split(url, "/")
}

func TokenIsVar(token string) bool {
	return varTokenRegexp.MatchString(token)
}

func VarsInCmd(cmd string) []string {
	return varTokenRegexp.FindAllString(cmd, -1)
}
