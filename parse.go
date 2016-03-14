package goshen

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

func parseConfig(configContent []byte) ([]*SshConfigEntry, error) {

	// normalize two or more lines to just two lines
	lineRe := regexp.MustCompile("\r?\n{2,}")
	configContent = lineRe.ReplaceAll(configContent, []byte("\n\n"))

	configScanner := bufio.NewScanner(bytes.NewReader(configContent))

	splitRe := regexp.MustCompile("\\s{1,}")

	entries := make([]*SshConfigEntry, 0)
	var entry *SshConfigEntry

	for configScanner.Scan() {

		line := strings.TrimSpace(configScanner.Text())

		// search for the beginning of a new block
		if strings.HasPrefix(line, "Host ") || strings.HasPrefix(line, "Match ") {
			entry = &SshConfigEntry{}
			entries = append(entries, entry)
		}

		directiveParts := splitRe.Split(line, 2)
		if len(directiveParts) == 2 {
			entry.Set(directiveParts[0], directiveParts[1])
		}
	}

	if err := configScanner.Err(); err != nil {
		return entries, err
	}

	return entries, nil
}
