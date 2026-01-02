package db

import (
	"fmt"

	"github.com/asdine/storm/v3"
	"github.com/internet-golf/internet-golf/pkg/resources"
	"github.com/internet-golf/internet-golf/pkg/utils"
)

type Db interface {
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
	config *utils.Config
	dbFile string
}

func NewDb(config *utils.Config, files *resources.FileManager) (Db, error) {
	dbFile := files.DbPath

	storm, stormOpenErr := storm.Open(dbFile)
	if stormOpenErr != nil {
		return nil, stormOpenErr
	}
	defer storm.Close()

	deploymentBucketErr := storm.Init(&Deployment{})
	if deploymentBucketErr != nil {
		return nil, fmt.Errorf("Error creating deployment bucket: %+v", deploymentBucketErr)
	}

	usersBucketErr := storm.Init(&ExternalUser{})
	if usersBucketErr != nil {
		return nil, fmt.Errorf("Error creating users bucket: %+v", usersBucketErr)
	}

	// create and return object that implements Db

	db := &StormDb{
		config: config,
		dbFile: dbFile,
	}

	return db, nil
}

func (s *StormDb) GetStorageDirectory() string {
	return s.config.DataDirectory
}

// i am sort of not sure why i open and close the database in each and every one
// of these methods, instead of keeping it open

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
