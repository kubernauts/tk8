package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kubernauts/tk8/pkg/common"
	"github.com/spf13/viper"
)

func (r *Rke) CreateCluster() error {

	provisioner := "rke"
	// validateJSON
	err := r.ValidateConfig()
	if err != nil {
		return err
	}

	// create RKE cluster config file
	err = r.CreateConfig()
	if err != nil {
		return err
	}

	err = getProvisioner(provisioner)
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
	configFileName := "aws-" + r.ClusterName + ".yaml"
	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)
	exists := isExistsClusterConfig(configFileName)
	if !exists {
		return fmt.Errorf("No such cluster exists with name - ", r.ClusterName)
	}
	go func() {
		Provisioners[provisioner].Destroy(nil)
	}()

	// Delete AWS cluster config file
	err := r.DeleteConfig()
	if err != nil {
		return fmt.Errorf("Error deleting cluster config ...")
	}

	return nil
}

func (r *Rke) CreateConfig() error {
	viper.New()
	viper.SetConfigType("yaml")

	configFileName := "rke-" + r.ClusterName + ".yaml"

	viper.SetConfigFile(configFileName)
	viper.AddConfigPath(common.REST_API_STORAGEPATH)

	viper.Set("rke.cluster_name", r.ClusterName)
	viper.Set("rke.node_os", r.NodeOs)
	viper.Set("rke.rke_aws_region", r.ClusterName)

	viper.Set("rke.authorization", r.ClusterName)
	viper.Set("rke.rke_node_instance_type", r.ClusterName)
	viper.Set("rke.node_count", r.ClusterName)
	viper.Set("rke.cloud_provider", r.CloudProvider)

	log.Println(viper.AllKeys())
	log.Println(viper.AllSettings())

	err := viper.WriteConfig()

	if err != nil {
		return err
	}
	return nil
}

func (r *Rke) ValidateConfig() error {
	if r.ClusterName == "" {
		return fmt.Errorf("Cluster name cannot be empty")
	}

	configFileName := "rke-" + r.ClusterName + ".yaml"
	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)

	if isExistsClusterConfig(configFileName) {
		return fmt.Errorf("Cluster name must be unique")
	}

	return nil
}

func (r *Rke) DeleteConfig() error {

	configFileName := "aws-" + r.ClusterName + ".yaml"
	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)

	err := os.Remove(configFileName)
	if err != nil {
		return err
	}
	return nil
}

func (r *Rke) UpdateConfig() error {
	return nil
}
