package roadmap

import "regexp"

var (
	outcomeNamePattern = regexp.MustCompile(`^(O[0-9]{2})-.+`)
	taskNamePattern    = regexp.MustCompile(`^(T[0-9]{3})-.+\.md$`)
)

func OutcomeID(name string) (string, bool) {
	matches := outcomeNamePattern.FindStringSubmatch(name)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}

func TaskID(name string) (string, bool) {
	matches := taskNamePattern.FindStringSubmatch(name)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}
