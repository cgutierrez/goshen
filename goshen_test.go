package goshen_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/cgutierrez/goshen"
)

func TestReadFile(t *testing.T) {

	configPath := "./ssh_config"

	config := goshen.NewSshConfig(configPath)

	if config.Path != configPath {
		t.Error("the config file path was not set")
	}

	if len(config.Entries) != 3 {
		t.Error("entry count doesn't match the amount in the sample config file")
	}
}

func TestGetSetArbitraryKey(t *testing.T) {
	config := goshen.NewSshConfig("./ssh_config")

	config.Entries[0].Set("SomethingThatDoesNotExist", "test-val")
	if config.Entries[0].Get("SomethingThatDoesNotExist") != "test-val" {
		t.Error("The arbitrary key was not set")
	}
}

func TestCreateConfig(t *testing.T) {

	config := goshen.NewSshConfig("")

	configFileName := "./ssh_config_saved"
	testHostName := "test-host"

	config.Entries = append(config.Entries, &goshen.SshConfigEntry{Host: testHostName})
	config.Save(configFileName)

	savedConfig, err := ioutil.ReadFile(configFileName)

	if err != nil {
		t.Errorf("an error occurred while reading the saved config. %s", err.Error())
	}

	if !strings.Contains(string(savedConfig), testHostName) {
		t.Error("The saved config does not contain the test host")
	}

	os.Remove(configFileName)
}
