package server

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	//"os"
	"github.com/kubernauts/tk8/api"
	"github.com/kubernauts/tk8/pkg/common"
	//"github.com/spf13/viper"
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

func (r *Rke) CreateCluster() error {

	configFileName := "rke-" + r.ClusterName + ".yaml"
	s := NewStore(common.REST_API_STORAGE, configFileName, common.REST_API_STORAGEPATH)
	provisioner := "rke"
	// validateJSON
	err := s.ValidateConfig()
	if err != nil {
		return err
	}

	err = getProvisioner(provisioner)
	if err != nil {
		return err
	}

	// create RKE cluster config file
	//l := NewLocalStore(configFileName, common.REST_API_STORAGEPATH)
	err = s.CreateConfig(r)
	if err != nil {
		return err
	}

	go func() {
		Provisioners[provisioner].Init(nil)
		Provisioners[provisioner].Setup(nil)
	}()

	return nil
}
func (r *Rke) DestroyCluster() error {
	provisioner := "rke"
	configFileName := "rke-" + r.ClusterName + ".yaml"
	s := NewStore(common.REST_API_STORAGE, configFileName, common.REST_API_STORAGEPATH)

	exists, _ := s.CheckConfigExists()
	if !exists {
		return fmt.Errorf("No such cluster exists with name - ", r.ClusterName)
	}
	go func() {
		Provisioners[provisioner].Destroy(nil)
	}()

	// Delete AWS cluster config file
	err := s.DeleteConfig()
	if err != nil {
		return fmt.Errorf("Error deleting cluster config ...")
	}

	return nil
}

func (r *Rke) GetCluster(name string) (api.Cluster, error) {

	configFileName := "rke-" + name + ".yaml"
	s := NewStore(common.REST_API_STORAGE, configFileName, common.REST_API_STORAGEPATH)

	exists, _ := s.CheckConfigExists()
	if !exists {
		return nil, fmt.Errorf("No cluster found with the provided name ::: ", name)
	}

	rkeConfig := &RkeYaml{}
	yamlFile, err := s.GetConfig()
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, rkeConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into rke config struct, %v", err)
	}

	return rkeConfig.Rke, nil
}
