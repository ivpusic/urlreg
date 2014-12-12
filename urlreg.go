package urlregex

import (
	"errors"
	"github.com/ivpusic/golog"
	"regexp"
	"strings"
)

var logger *golog.Logger

type UrlRegex struct {
	Regex *regexp.Regexp
}

func init() {
	logger = golog.GetLogger("github.com/ivpusic/urlreg")
	logger.Level = golog.WARN
}

// Method will collect named parameters values.
// If there are not matches, method will return error.
func (u *UrlRegex) Match(url string) (map[string]string, error) {
	params := u.Regex.SubexpNames()
	matches := u.Regex.FindAllStringSubmatch(url, -1)

	if matches == nil {
		return nil, errors.New("cannot match " + url)
	}

	result := map[string]string{}
	for i, n := range matches[0] {
		if len(params[i]) > 0 {
			result[params[i]] = n
		}
	}

	return result, nil
}

// Converting regex into UrlRegex object
// Function will build regex with named groups to be able to extract parameters later
func Pattern(pattern string) *UrlRegex {
	parts := strings.Split(pattern, "/")
	regex := "^"

	for i, part := range parts {
		if len(part) > 0 {
			if i > 0 {
				regex += "\\/"
			}

			// do we have special character?
			switch part[0] == ':' {
			case true:
				groupName := "(?P<" + part[1:] + ">"
				regex += groupName + ".[^\\/]*)"
			case false:
				regex += regexp.QuoteMeta(part)
			}
		}
	}

	if pattern[len(pattern)-1] == '/' {
		regex += "\\/$"
	} else {
		regex += "$"
	}

	logger.Debug("regex: " + regex)

	return &UrlRegex{
		Regex: regexp.MustCompile(regex),
	}
}
