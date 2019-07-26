package server

import (
	"fmt"
	"github.com/kubernauts/tk8/api"
	"github.com/kubernauts/tk8/pkg/common"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type RkeYaml struct {
	Rke *Rke `yaml:"rke"`
}

type Rke struct {
	ClusterName         string `yaml:"cluster_name" json:"cluster_name"`
	NodeOs              string `yaml:"node_os" json:"node_os"`
	RkeAwsRegion        string `yaml:"rke_aws_region" json:"rke_aws_region"`
	Authorization       string `yaml:"authorization" json:"authorization"`
	RkeNodeInstanceType string `yaml:"rke_node_instance_type" json:"rke_node_instance_type"`
	NodeCount           int    `yaml:"node_count" json:"node_count"`
	CloudProvider       string `yaml:"cloud_provider" json:"cloud_provider"`
}

// CreateCluster creates RKE cluster
func (r *Rke) CreateCluster() error {

	configFileName := "rke-" + r.ClusterName + ".yaml"
	s := NewStore(common.REST_API_STORAGE, configFileName, common.REST_API_STORAGEPATH, common.REST_API_STORAGEREGION)
	provisioner := "rke"
	// validateJSON
	err := s.ValidateConfig()
	if err != nil {
		logrus.Errorf("Error validating config --  %s", err.Error())
		return err
	}

	err = getProvisioner(provisioner)
	if err != nil {
		return err
	}

	// create RKE cluster config file
	err = s.CreateConfig(r)
	if err != nil {
		logrus.Errorf("Error creating config --  %s", err.Error())
		return err
	}

	go func() {
		Provisioners[provisioner].Init(nil)
		Provisioners[provisioner].Setup(nil)
	}()

	return nil
}

// DestroyCluster destroys given RKE cluster
func (r *Rke) DestroyCluster() error {
	configFileName := "rke-" + r.ClusterName + ".yaml"
	s := NewStore(common.REST_API_STORAGE, configFileName, common.REST_API_STORAGEPATH, common.REST_API_STORAGEREGION)

	exists, _ := s.CheckConfigExists()
	if !exists {
		logrus.Errorf("Error no such cluster existswith name --  %s", r.ClusterName)
		return fmt.Errorf("No such cluster exists with name - %s", r.ClusterName)
	}
	go func() {
		Provisioners["rke"].Destroy(nil)
	}()

	// Delete AWS cluster config file
	err := s.DeleteConfig()
	if err != nil {
		logrus.Errorf("Error deleting cluster --  %s", err.Error())
		return fmt.Errorf("Error deleting cluster named  - %s", r.ClusterName)
	}

	return nil
}

// GetCluster for RKE
func (r *Rke) GetCluster(name string) (api.Cluster, error) {

	configFileName := "rke-" + name + ".yaml"
	s := NewStore(common.REST_API_STORAGE, configFileName, common.REST_API_STORAGEPATH, common.REST_API_STORAGEREGION)

	exists, _ := s.CheckConfigExists()
	if !exists {
		logrus.Errorf("Error no such cluster existswith name --  %s", name)
		return nil, fmt.Errorf("No cluster found with the provided name ::: %s", name)
	}

	rkeConfig := &RkeYaml{}
	yamlFile, err := s.GetConfig()
	if err != nil {
		logrus.Errorf("Error unable to get cluster details , %v", err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, rkeConfig)
	if err != nil {
		return nil, fmt.Errorf("Error unable to get cluster, %v", err)
	}

	return rkeConfig.Rke, nil
}
