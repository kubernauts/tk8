package server

import (
	"fmt"
	"github.com/kubernauts/tk8/api"
	"github.com/kubernauts/tk8/pkg/common"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
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

// CreateCluster creates EKS cluster
func (e *Eks) CreateCluster() error {

	// create AWS cluster config file
	configFileName := "eks-" + e.ClusterName + ".yaml"
	s := NewStore(common.REST_API_STORAGE, configFileName, common.REST_API_STORAGEPATH, common.REST_API_STORAGEREGION)
	provisioner := "eks"
	// validateJSON
	err := s.ValidateConfig()
	if err != nil {
		logrus.Errorf("Error validating config ::: %s", err)
		return err
	}

	err = getProvisioner(provisioner)
	if err != nil {
		logrus.Errorf("Error getting provisioner ::: %s", err)
		return err
	}

	// create EKS cluster config file
	err = s.CreateConfig(e)
	if err != nil {
		logrus.Errorf("Error creating config ::: %s", err)
		return err
	}

	go func() {
		Provisioners[provisioner].Init(nil)
		Provisioners[provisioner].Setup(nil)
	}()

	return nil
}

// DestroyCluster destroys EKS cluster
func (e *Eks) DestroyCluster() error {
	configFileName := "eks-" + e.ClusterName + ".yaml"

	s := NewStore(common.REST_API_STORAGE, configFileName, common.REST_API_STORAGEPATH, common.REST_API_STORAGEREGION)
	exists, _ := s.CheckConfigExists()
	if !exists {
		logrus.Errorf("No such cluster with name  ::: %s", e.ClusterName)
		return fmt.Errorf("No such cluster with name - %s", e.ClusterName)
	}
	go func() {
		Provisioners["eks"].Destroy(nil)
	}()

	// Delete AWS cluster config file
	err := s.DeleteConfig()
	if err != nil {
		logrus.Errorf("Error deleting config ::: %s", err)
		return fmt.Errorf("Error deleting cluster named %s", e.ClusterName)
	}
	return nil
}

// GetCluster gets the details of thge requested EKS cluster
func (e *Eks) GetCluster(name string) (api.Cluster, error) {

	configFileName := "eks-" + name + ".yaml"
	s := NewStore(common.REST_API_STORAGE, configFileName, common.REST_API_STORAGEPATH, common.REST_API_STORAGEREGION)
	exists, _ := s.CheckConfigExists()
	if !exists {
		logrus.Errorf("Error getting config ::: %s", e.ClusterName)
		return nil, fmt.Errorf("No cluster found with the provided name ::: %s", name)
	}

	yamlFile, err := s.GetConfig()
	if err != nil {
		logrus.Errorf("Error getting details of cluster named %s , Error details are %s ", name, err.Error())
		return nil, err
	}
	eksConfig := &EksYaml{}
	err = yaml.Unmarshal(yamlFile, eksConfig)
	if err != nil {
		logrus.Errorf("unable to decode into rke config struct, %v", err)
		return nil, fmt.Errorf("unable to decode into rke config struct, %v", err)
	}

	return eksConfig.Eks, nil
}
