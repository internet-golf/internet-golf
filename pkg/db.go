package internetgolf

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/asdine/storm/v3"
)

type StorageSettings struct {
	DataDirectory string
}

// pass in an empty string to use the default data directory
func (s *StorageSettings) Init(nonDefaultDataDirectory string) {
	var dataDirectoryError error
	s.DataDirectory, dataDirectoryError = getDataDirectory(nonDefaultDataDirectory)
	if dataDirectoryError != nil {
		panic("Could not create data directory: " + dataDirectoryError.Error())
	}
	fmt.Printf("Initialized data directory to %s\n", s.DataDirectory)
}

type Storage interface {
	Init() error
	SaveDeployments(d []Deployment) error
	GetDeployments() ([]Deployment, error)
}

// i found the database package "storm" on github and didn't realize until after
// i had created a storage implementation using it that it hasn't been updated
// in 5 years. i guess it's fine???
type StormStorage struct {
	Settings StorageSettings
	dbFile   string
}

func (s *StormStorage) Init() error {
	s.dbFile = path.Join(s.Settings.DataDirectory, "internet.db")

	db, dbOpenErr := storm.Open(s.dbFile)
	if dbOpenErr != nil {
		return dbOpenErr
	}
	defer db.Close()

	deploymentBucketErr := db.Init(&Deployment{})
	if deploymentBucketErr != nil {
		return deploymentBucketErr
	}

	return nil
}

func (s *StormStorage) GetDeployments() ([]Deployment, error) {
	db, dbOpenErr := storm.Open(s.dbFile)
	if dbOpenErr != nil {
		return nil, dbOpenErr
	}
	defer db.Close()

	var existingDeployments []Deployment
	loadingExistingErr := db.All(&existingDeployments)
	if loadingExistingErr != nil {
		return nil, loadingExistingErr
	}
	return existingDeployments, nil
}

func (s *StormStorage) SaveDeployments(d []Deployment) error {
	db, dbOpenErr := storm.Open(s.dbFile)
	if dbOpenErr != nil {
		return dbOpenErr
	}
	defer db.Close()

	for _, d := range d {
		if !d.DontPersist {
			saveErr := db.Save(&d)
			if saveErr != nil {
				fmt.Printf(
					"could not save deployment with name %s: %+v\n",
					d.Name, saveErr,
				)
			}
		}
	}

	return nil
}

// receives a dataDirectoryPath; translates "$HOME" to the user's home
// directory; creates a directory at the path if it doesn't already exist.
// why is this a separate function ???
func getDataDirectory(dataDirectoryPath string) (string, error) {
	if strings.Index(dataDirectoryPath, "$HOME") != -1 {
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
