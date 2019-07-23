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
	Region     string
}

func NewStore(storetype, name, path, region string) api.ConfigStore {
	switch storetype {
	case "s3":
		return NewS3Store(name, path, region)
	case "local":
		return NewLocalStore(name, path)
	default:
		return nil
	}
}
