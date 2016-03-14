package goshen

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

type SshConfig struct {
	Entries []*SshConfigEntry
	Path    string
}

func writeConfigFile(configFile string, lines []string) error {
	file, err := os.Create(configFile)

	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, line := range lines {
		_, err := writer.WriteString(line)

		if err != nil {
			return err
		}
	}

	return nil
}

func (config *SshConfig) Save(outputPath string) error {

	lines := make([]string, 0)

	for _, entry := range config.Entries {

		// skip creating the entry if it doesn't have a Host or Match directive
		if entry.Match == "" && entry.Host == "" {
			continue
		}

		if entry.Match != "" {
			lines = append(lines, fmt.Sprintf("%s %s\n", "Match", entry.Match))
		}

		if entry.Host != "" {
			lines = append(lines, fmt.Sprintf("%s %s\n", "Host", entry.Host))
		}

		structValue := reflect.ValueOf(entry).Elem()
		entryType := structValue.Type()

		for i := 0; i < structValue.NumField(); i++ {
			if entryType.Field(i).Name == "Host" || entryType.Field(i).Name == "Match" || entryType.Field(i).Name == "additional" {
				continue
			}

			structFieldValue := structValue.FieldByName(entryType.Field(i).Name)

			// keep the field from written if it doesn't have a value
			if structFieldValue.String() == "" {
				continue
			}

			lines = append(lines, fmt.Sprintf("%s %s\n", entryType.Field(i).Name, structFieldValue.String()))
		}

		// write additional keys
		for key, val := range entry.additional {
			lines = append(lines, fmt.Sprintf("%s %s\n", key, val))
		}

		// separate each entry with double new lines
		lines = append(lines, "\n")
	}

	if outputPath == "" {
		outputPath = config.Path
	}

	writeConfigFile(outputPath, lines)

	return nil
}

func NewSshConfig(configPath string) *SshConfig {

	if len(configPath) == 0 {
		return &SshConfig{}
	}

	configFile, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Fatal(err)
	}

	entries, err := parseConfig(configFile)

	if err != nil {
		log.Fatal(err)
	}

	config := &SshConfig{Path: configPath, Entries: entries}
	return config
}
