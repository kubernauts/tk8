package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kubernauts/tk8/pkg/common"
	"github.com/spf13/viper"
)

func (e *Eks) CreateCluster() error {

	provisioner := "eks"
	// validateJSON
	err := e.ValidateConfig()
	if err != nil {
		return err
	}

	// create EKS cluster config file
	err = e.CreateConfig()
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

func (e *Eks) DestroyCluster() error {

	provisioner := "eks"
	configFileName := "eks-" + e.ClusterName + ".yaml"
	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)
	exists := isExistsClusterConfig(configFileName)
	if !exists {
		return fmt.Errorf("No such cluster exists with name - ", e.ClusterName)
	}
	go func() {
		Provisioners[provisioner].Destroy(nil)
	}()

	// Delete AWS cluster config file
	err := e.DeleteConfig()
	if err != nil {
		return fmt.Errorf("Error deleting cluster config ...")
	}

	return nil

}

func (e *Eks) CreateConfig() error {

	viper.New()
	viper.SetConfigType("yaml")

	configFileName := "eks-" + e.ClusterName + ".yaml"

	viper.SetConfigFile(configFileName)
	viper.AddConfigPath(common.REST_API_STORAGEPATH)

	viper.Set("eks.cluster-name", e.ClusterName)
	viper.Set("eks.aws_region", e.AwsRegion)
	viper.Set("eks.node-instance-type", e.NodeInstanceType)
	viper.Set("eks.desired-capacity", e.DesiredCapacity)
	viper.Set("eks.autoscalling-max-size", e.AutoscallingMaxSize)
	viper.Set("eks.autoscalling-min-size", e.AutoscallingMinSize)
	viper.Set("eks.key-file-path", "~/.ssh/id_rsa.pub")

	log.Println(viper.AllKeys())
	log.Println(viper.AllSettings())

	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func (e *Eks) ValidateConfig() error {
	if e.ClusterName == "" {
		return fmt.Errorf("Cluster name cannot be empty")
	}

	configFileName := "eks-" + e.ClusterName + ".yaml"
	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)

	if isExistsClusterConfig(configFileName) {
		return fmt.Errorf("Cluster name must be unique")
	}

	return nil
}

func (e *Eks) DeleteConfig() error {

	configFileName := "eks-" + e.ClusterName + ".yaml"
	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)

	err := os.Remove(configFileName)
	if err != nil {
		return err
	}
	return nil
}

func (e *Eks) UpdateConfig() error {
	return nil
}
