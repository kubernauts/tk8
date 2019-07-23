package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kubernauts/tk8/api"
)

func NewS3Store(name, bucketName string) *S3Store {
	return &S3Store{
		FileName:   name,
		BucketName: bucketName,
	}
}

func (s *S3Store) DeleteConfig() error {
	//select Region to use.
	conf := aws.Config{Region: aws.String("eu-west-1")}
	sess := session.New(&conf)
	svc := s3.New(sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(s.FileName),
	})

	if err != nil {
		return fmt.Errorf("Erorr deleting cluster ::")
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(s.FileName),
	})

	if err != nil {
		return fmt.Errorf("Erorr deleting cluster ::")
	}

	return nil
}

func (s *S3Store) ValidateConfig() error {

	return nil
}

func (s *S3Store) UpdateConfig() error {

	return nil
}

func (s *S3Store) CheckConfigExists() (bool, error) {
	var err error
	conf := aws.Config{Region: aws.String("eu-west-1")}

	svc := s3.New(session.New(&conf))
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(s.BucketName),
		MaxKeys: aws.Int64(100),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return false, fmt.Errorf(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				return false, fmt.Errorf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return false, fmt.Errorf(err.Error())
		}
		return false, err
	}

	for _, item := range result.Contents {
		if *item.Key == s.FileName {
			return true, nil
		}
	}
	return false, nil
}

func (s *S3Store) GetConfig() ([]byte, error) {

	tempPath := "/tmp"
	path := filepath.Join(tempPath, s.FileName)
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	conf := aws.Config{Region: aws.String("eu-west-1")}
	sess := session.New(&conf)
	downloader := s3manager.NewDownloader(sess)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(s.BucketName),
			Key:    aws.String(s.FileName),
		})
	if err != nil {
		return nil, err
	}

	l := NewLocalStore(s.FileName, tempPath)
	yamlFile, err := l.GetConfig()
	if err != nil {
		return nil, err
	}

	deleteFileLocally(path)
	return yamlFile, nil
}

func (s *S3Store) CreateConfig(t api.Cluster) error {

	tempPath := "/tmp"
	l := NewLocalStore(s.FileName, tempPath)

	err := l.CreateConfig(t)
	if err != nil {
		return err
	}
	// Check for existence of s3 bucket
	err = checkIfBucketExists((s.BucketName))
	if err != nil {
		return err
	}

	// check if there exists a file with same cluster name
	exists, err := s.CheckConfigExists()
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("cluster already exists with the same name, please provide a unique name")
	}

	fullpath := filepath.Join(tempPath, s.FileName)
	// upload the locally created config to the s3 bucket
	err = uploadFileToS3(fullpath, s.BucketName)
	if err != nil {
		return err
	}

	//	upon successful upload delete the locally created file
	filename := filepath.Join(tempPath, s.FileName)
	err = deleteFileLocally(filename)
	if err != nil {
		return err
	}

	return nil
}

func deleteFileLocally(path string) error {
	// delete file
	var err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func uploadFileToS3(filename, bucketname string) error {

	//select Region to use.
	conf := aws.Config{Region: aws.String("eu-west-1")}
	sess := session.New(&conf)
	svc := s3manager.NewUploader(sess)

	// open file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open file", filename, err)
		return err
	}
	defer file.Close()

	_, err = svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(filepath.Base(filename)),
		Body:   file,
	})

	if err != nil {
		fmt.Println("error", err)
		return err
	}
	fmt.Println("left uploadS3 ...")

	return nil
}

func checkIfBucketExists(bucketName string) error {
	conf := aws.Config{Region: aws.String("eu-west-1")}
	fmt.Println("Entered checkbucket ...")
	svc := s3.New(session.New(&conf))
	input := &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}

	_, err := svc.HeadBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return fmt.Errorf(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				return fmt.Errorf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return fmt.Errorf(err.Error())
		}
		return err
	}
	fmt.Println("Left checkbucket ... ")
	return nil
}
