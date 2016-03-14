package goshen

import (
  "reflect"
  "strings"
  "os/user"
  "path"
)

type SshConfigEntry struct {
  Host string
  Match string
  HostName string
  User string
  Port string
  UserKnownHostsFile string
  StrictHostKeyChecking string
  IdentityFile string
  LogLevel string
  additional map[string]string
}

// swap out the shortcut (~) to the users home directory
// with the full path to the home directory
func expandHomeDir(filePath string) (string, error) {

  if strings.HasPrefix(filePath, "~") {

    usr, err := user.Current()

    if err != nil {
      return filePath, err
    }

    filePath = path.Join(usr.HomeDir, strings.Replace(filePath, "~", "", 1))
  }

  return filePath, nil
}

// Set a key on the config entry. If the directive name doesn't exist on the struct,
// it's added to the 'additional' key and can be accessed using the Get method
func (entry *SshConfigEntry) Set(directive string, value string) (err error) {

  if directive == "IdentityFile" {
    value, err = expandHomeDir(value)

    if err != nil {
      return err
    }
  }

  if entry.additional == nil {
    entry.additional = make(map[string]string)
  }

  structValue := reflect.ValueOf(entry).Elem()
  structFieldValue := structValue.FieldByName(directive)

  if !structFieldValue.IsValid() {
    entry.additional[directive] = value
  }

  if structFieldValue.IsValid() {
    structFieldValue.Set(reflect.ValueOf(value))
  }

  return nil
}

// Provide access to the config entry directives. If the given directive
// doesn't exist in the struct, it checks the 'additional' key
func (entry *SshConfigEntry) Get(directive string) (string) {

  if entry.additional == nil {
    entry.additional = make(map[string]string)
  }

  structValue := reflect.ValueOf(entry).Elem()
  structFieldValue := structValue.FieldByName(directive)

  if structFieldValue.IsValid() {
    return structFieldValue.String()
  }

  if !structFieldValue.IsValid() {
    if val, ok := entry.additional[directive]; ok {
      return val
    }
  }

  return ""
}