package settings

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DataDirectory string
	LocalOnly     bool
	Verbose       bool
	AdminApiPort  string
}

// creates a new config object with the data that you pass in.
//
// note that `dataDirectory` is given special treatment; the string "$HOME" is
// replaced with the current OS user's home directory, and if no directory at
// the given path exists, it's immediately created.
//
// this should only need to be called once for every time the server is started.
func NewConfig(dataDirectory string, localOnly bool, verbose bool, adminApiPort string) *Config {
	dataDirectory, dataDirectoryError := setupDataDirectory(dataDirectory)
	if dataDirectoryError != nil {
		panic("Could not create data directory: " + dataDirectoryError.Error())
	}
	fmt.Printf("Initialized data directory to %s\n", dataDirectory)

	return &Config{
		DataDirectory: dataDirectory,
		LocalOnly:     localOnly,
		Verbose:       verbose,
		AdminApiPort:  adminApiPort,
	}
}

// receives a dataDirectoryPath; translates "$HOME" to the user's home
// directory; creates a directory at the path if it doesn't already exist.
func setupDataDirectory(dataDirectoryPath string) (string, error) {
	if strings.Contains(dataDirectoryPath, "$HOME") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New(
				"could not obtain home directory, so data directory could not be " +
					"created using it. please manually configure the data directory",
			)
		}
		// hopefully this replaceAll doesn't have weird consequences -
		// everything still seems to work here on windows
		homeDir = strings.ReplaceAll(homeDir, "\\", "/")
		dataDirectoryPath = strings.ReplaceAll(dataDirectoryPath, "$HOME", homeDir)
	}

	if _, err := os.Lstat(dataDirectoryPath); err != nil {
		fmt.Printf("Creating data directory at %v\n", dataDirectoryPath)
		// TODO: hope that 0750 works for permissions. can Caddy access the result???
		// will this work recursively?
		if os.Mkdir(dataDirectoryPath, 0750) != nil {
			return "", errors.New("could not create data directory at " + dataDirectoryPath)
		}
	}

	return dataDirectoryPath, nil
}
