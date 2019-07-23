package server

import (
	"github.com/kubernauts/tk8/api"
)

type Storage struct {
	StoreType string
}
type LocalStore struct {
	*Storage
	FileName string
	FilePath string
}

type S3Store struct {
	*Storage
	FileName   string
	BucketName string
}

func NewStore(storetype, name, location string) api.ConfigStore {
	switch storetype {
	case "s3":
		return NewS3Store(name, location)
	case "local":
		return NewLocalStore(name, location)
	default:
		return nil
	}
}
