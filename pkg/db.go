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

type Db interface {
	Init(settings StorageSettings) error
	GetStorageDirectory() string
	SaveDeployments(d []Deployment) error
	GetDeployments() ([]Deployment, error)
	SaveExternalUser(u ExternalUser) error
	GetExternalUser(externalId string) (ExternalUser, error)
	SaveBearerToken(b BearerToken) error
	GetBearerToken(string) (BearerToken, error)
}

// i found the database package "storm" on github and didn't realize until after
// i had created a storage implementation using it that it hasn't been updated
// in 5 years. i guess it's fine??? implements the `Db` interface.
type StormDb struct {
	settings StorageSettings
	dbFile   string
}

func (s *StormDb) Init(settings StorageSettings) error {
	s.settings = settings
	s.dbFile = path.Join(settings.DataDirectory, "internet.db")

	db, dbOpenErr := storm.Open(s.dbFile)
	if dbOpenErr != nil {
		return dbOpenErr
	}
	defer db.Close()

	deploymentBucketErr := db.Init(&Deployment{})
	if deploymentBucketErr != nil {
		return fmt.Errorf("Error creating deployment bucket: %+v", deploymentBucketErr)
	}

	usersBucketErr := db.Init(&ExternalUser{})
	if usersBucketErr != nil {
		return fmt.Errorf("Error creating users bucket: %+v", usersBucketErr)
	}

	return nil
}

func (s *StormDb) GetStorageDirectory() string {
	return s.settings.DataDirectory
}

func (s *StormDb) GetDeployments() ([]Deployment, error) {
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

func (s *StormDb) SaveDeployments(d []Deployment) error {
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
					"could not save deployment %s: %+v\n",
					d.Url, saveErr,
				)
			}
		}
	}

	return nil
}

func (s *StormDb) SaveExternalUser(u ExternalUser) error {
	db, dbOpenErr := storm.Open(s.dbFile)
	if dbOpenErr != nil {
		return dbOpenErr
	}
	defer db.Close()

	return db.Save(&u)
}

func (s *StormDb) GetExternalUser(externalId string) (ExternalUser, error) {
	db, dbOpenErr := storm.Open(s.dbFile)
	if dbOpenErr != nil {
		return ExternalUser{}, dbOpenErr
	}
	defer db.Close()

	var result ExternalUser
	err := db.Get("ExternalUser", externalId, &result)
	if err != nil {
		return ExternalUser{}, err
	}

	return result, nil
}

func (s *StormDb) SaveBearerToken(token BearerToken) error {
	db, dbOpenErr := storm.Open(s.dbFile)
	if dbOpenErr != nil {
		return dbOpenErr
	}
	defer db.Close()

	return db.Save(&token)
}

func (s *StormDb) GetBearerToken(id string) (BearerToken, error) {
	db, dbOpenErr := storm.Open(s.dbFile)
	if dbOpenErr != nil {
		return BearerToken{}, dbOpenErr
	}
	defer db.Close()

	var result BearerToken
	err := db.Get("BearerToken", id, &result)
	if err != nil {
		return BearerToken{}, err
	}

	return result, nil
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
