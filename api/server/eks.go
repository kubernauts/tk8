package server

import (
	"fmt"
	"gopkg.in/yaml.v2"

	//"os"

	"github.com/kubernauts/tk8/api"
	"github.com/kubernauts/tk8/pkg/common"
	//"github.com/spf13/viper"
)

type EksYaml struct {
	Eks *Eks `yaml:"eks"`
}
type Eks struct {
	ClusterName         string `yaml:"cluster-name" json:"cluster-name"`
	AwsRegion           string `yaml:"aws_region" json:"aws_region"`
	NodeInstanceType    string `yaml:"node-instance-type" json:"node-instance-type"`
	DesiredCapacity     int    `yaml:"desired-capacity" json:"desired-capacity"`
	AutoscallingMaxSize int    `yaml:"autoscalling-max-size" json:"autoscalling-max-size"`
	AutoscallingMinSize int    `yaml:"autoscalling-min-size" json:"autoscalling-min-size"`
	KeyFilePath         string `yaml:"key-file-path" json:"key-file-path"`
}

func (e *Eks) CreateCluster() error {

	// create AWS cluster config file
	configFileName := "eks-" + e.ClusterName + ".yaml"
	s := NewStore(common.REST_API_STORAGE, configFileName, common.REST_API_STORAGEPATH)
	provisioner := "eks"
	// validateJSON
	err := s.ValidateConfig()
	if err != nil {
		return err
	}

	err = getProvisioner(provisioner)
	if err != nil {
		return err
	}

	// create EKS cluster config file
	//l := NewLocalStore(configFileName, common.REST_API_STORAGEPATH)
	err = s.CreateConfig(e)
	if err != nil {
		return err
	}

	go func() {
		Provisioners[provisioner].Init(nil)
		Provisioners[provisioner].Setup(nil)
	}()

	return nil
}

func (e *Eks) DestroyCluster() error {

	provisioner := "eks"
	configFileName := "eks-" + e.ClusterName + ".yaml"

	s := NewStore(common.REST_API_STORAGE, configFileName, common.REST_API_STORAGEPATH)
	exists, _ := s.CheckConfigExists()
	if !exists {
		return fmt.Errorf("No such cluster exists with name - ", e.ClusterName)
	}
	//	go func() {
	Provisioners[provisioner].Destroy(nil)
	//	}()

	// Delete AWS cluster config file
	err := s.DeleteConfig()
	if err != nil {
		return fmt.Errorf("Error deleting cluster config ...")
	}

	return nil

}

func (e *Eks) GetCluster(name string) (api.Cluster, error) {

	configFileName := "eks-" + name + ".yaml"
	s := NewStore(common.REST_API_STORAGE, configFileName, common.REST_API_STORAGEPATH)
	exists, _ := s.CheckConfigExists()
	if !exists {
		return nil, fmt.Errorf("No cluster found with the provided name ::: ", name)
	}

	yamlFile, err := s.GetConfig()
	eksConfig := &EksYaml{}
	err = yaml.Unmarshal(yamlFile, eksConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into rke config struct, %v", err)
	}

	return eksConfig.Eks, nil
}
