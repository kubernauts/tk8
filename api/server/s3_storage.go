package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kubernauts/tk8/api"
	"github.com/kubernauts/tk8/pkg/common"
)

func NewS3Store(name, bucketName, region string) *S3Store {
	return &S3Store{
		FileName:   name,
		BucketName: bucketName,
		Region:     region,
	}
}

func (s *S3Store) DeleteConfig() error {
	//select Region to use.
	conf := aws.Config{Region: aws.String(s.Region)}
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

func (s *S3Store) GetConfigs() (api.AllClusters, error) {
	var err error
	tempPath := "/tmp"
	conf := aws.Config{Region: aws.String(s.Region)}

	sess := session.New(&conf)
	downloader := s3manager.NewDownloader(sess)

	svc := s3.New(sess)
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(s.BucketName),
		MaxKeys: aws.Int64(100),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, fmt.Errorf(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				return nil, fmt.Errorf(aerr.Error())
			}
		} else {
			// Message from an error.
			return nil, fmt.Errorf(err.Error())
		}
	}

	for _, item := range result.Contents {
		path := filepath.Join(tempPath, *item.Key)
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		_, err = downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(s.BucketName),
				Key:    aws.String(*item.Key),
			})
		if err != nil {
			return nil, err
		}
	}

	l := NewLocalStore("", tempPath)
	allClusters, err := l.GetConfigs()
	if err != nil {
		return nil, err
	}
	return allClusters, nil
}

func (s *S3Store) CheckConfigExists() (bool, error) {
	var err error
	conf := aws.Config{Region: aws.String(s.Region)}

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
			// Message from an error.
			return false, fmt.Errorf(err.Error())
		}
	}

	for _, item := range result.Contents {
		fmt.Println("Item --- ", *item.Key, " --- and --- filename -- ", s.FileName)
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

	conf := aws.Config{Region: aws.String(s.Region)}
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
	conf := aws.Config{Region: aws.String(common.REST_API_STORAGEREGION)}
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
	return nil
}

func checkIfBucketExists(bucketName string) error {
	conf := aws.Config{Region: aws.String("eu-west-1")}
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
	return nil
}

func getBucketObjects(sess *session.Session) {
	query := &s3.ListObjectsV2Input{
		Bucket: aws.String(os.Args[2]),
	}
	svc := s3.New(sess)

	// Flag used to check if we need to go further
	truncatedListing := true

	for truncatedListing {
		resp, err := svc.ListObjectsV2(query)

		if err != nil {
			// Print the error.
			fmt.Println(err.Error())
			return
		}
		// Get all files
		getObjectsAll(resp, svc)
		// Set continuation token
		query.ContinuationToken = resp.NextContinuationToken
		truncatedListing = *resp.IsTruncated
	}
}

func getObjectsAll(bucketObjectsList *s3.ListObjectsV2Output, s3Client *s3.S3) {
	// Iterate through the files inside the bucket
	for _, key := range bucketObjectsList.Contents {
		fmt.Println(*key.Key)
		destFilename := *key.Key
		if strings.HasSuffix(*key.Key, "/") {
			fmt.Println("Got a directory")
			continue
		}
		if strings.Contains(*key.Key, "/") {
			var dirTree string
			// split
			s3FileFullPathList := strings.Split(*key.Key, "/")
			fmt.Println(s3FileFullPathList)
			fmt.Println("destFilename " + destFilename)
			for _, dir := range s3FileFullPathList[:len(s3FileFullPathList)-1] {
				dirTree += "/" + dir
			}
			os.MkdirAll(os.Args[3]+"/"+dirTree, 0775)
		}
		out, err := s3Client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(os.Args[2]),
			Key:    key.Key,
		})
		if err != nil {
			logrus.Fatal(err)
		}
		destFilePath := os.Args[3] + destFilename
		destFile, err := os.Create(destFilePath)
		if err != nil {
			logrus.Fatal(err)
		}
		bytes, err := io.Copy(destFile, out.Body)
		if err != nil {
			logrus.Fatal(err)
		}
		fmt.Printf("File %s contanin %d bytes\n", destFilePath, bytes)
		out.Body.Close()
		destFile.Close()
	}
}
