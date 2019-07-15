package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/kubernauts/tk8/pkg/common"
	"gopkg.in/yaml.v2"
)

func ReadClusterConfigs() AllClusters {

	files, _ := ioutil.ReadDir(common.REST_API_STORAGEPATH)
	clusters := make(AllClusters)
	for _, file := range files {
		switch {
		case strings.Contains(file.Name(), "aws-"):
			configFileName := filepath.Join(common.REST_API_STORAGEPATH, file.Name())

			awsConfig := &AwsYaml{}
			yamlFile, err := ioutil.ReadFile(configFileName)
			if err != nil {
				log.Printf("yamlFile.Get err   #%v ", err)
			}
			err = yaml.Unmarshal(yamlFile, awsConfig)

			if err != nil {
				fmt.Printf("unable to decode into aws config struct, %v", err)
				continue
			}
			clusters["aws"] = append([]Cluster{awsConfig.Aws})

		case strings.Contains(file.Name(), "eks-"):
			configFileName := filepath.Join(common.REST_API_STORAGEPATH, file.Name())
			eksConfig := &EksYaml{}

			yamlFile, err := ioutil.ReadFile(configFileName)
			if err != nil {
				log.Printf("yamlFile.Get err   #%v ", err)
			}
			err = yaml.Unmarshal(yamlFile, eksConfig)
			if err != nil {
				fmt.Printf("unable to decode into eks config struct, %v", err)
				continue
			}
			clusters["eks"] = append([]Cluster{eksConfig.Eks})

		case strings.Contains(file.Name(), "rke-"):
			configFileName := filepath.Join(common.REST_API_STORAGEPATH, file.Name())
			rkeConfig := &RkeYaml{}
			yamlFile, err := ioutil.ReadFile(configFileName)
			if err != nil {
				log.Printf("yamlFile.Get err   #%v ", err)
			}
			err = yaml.Unmarshal(yamlFile, rkeConfig)
			if err != nil {
				fmt.Printf("unable to decode into rke config struct, %v", err)
				continue
			}
			clusters["rke"] = append([]Cluster{rkeConfig.Rke})
		}
	}
	return clusters
}

func decodeAwsClusterConfigToStruct(name string) (*Aws, error) {
	configFileName := "aws-" + name + ".yaml"
	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)

	exists := isExistsClusterConfig(configFileName)

	if !exists {
		return nil, fmt.Errorf("No cluster found with the provided name ::: ", name)
	}

	awsConfig := &AwsYaml{}
	yamlFile, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, awsConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into rke config struct, %v", err)
	}

	return awsConfig.Aws, nil
}

func decodeEksClusterConfigToStruct(name string) (*Eks, error) {
	configFileName := "eks-" + name + ".yaml"
	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)

	exists := isExistsClusterConfig(configFileName)

	if !exists {
		return nil, fmt.Errorf("No cluster found with the provided name ::: ", name)
	}
	eksConfig := &EksYaml{}
	yamlFile, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, eksConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into rke config struct, %v", err)
	}

	return eksConfig.Eks, nil
}

func decodeRkeClusterConfigToStruct(name string) (*Rke, error) {

	configFileName := "rke-" + name + ".yaml"
	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)

	exists := isExistsClusterConfig(configFileName)

	if !exists {
		return nil, fmt.Errorf("No cluster found with the provided name ::: ", name)
	}

	rkeConfig := &RkeYaml{}
	yamlFile, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, rkeConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into rke config struct, %v", err)
	}

	return rkeConfig.Rke, nil
}
