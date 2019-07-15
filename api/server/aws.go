package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kubernauts/tk8/pkg/common"
	"github.com/spf13/viper"
)

func (a *Aws) CreateCluster() error {

	provisioner := "aws"
	// validateJSON
	err := a.ValidateConfig()
	if err != nil {
		return err
	}

	// create AWS cluster config file
	err = a.CreateConfig()
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

func (a *Aws) DestroyCluster() error {

	provisioner := "aws"
	configFileName := "aws-" + a.Clustername + ".yaml"
	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)
	exists := isExistsClusterConfig(configFileName)
	if !exists {
		return fmt.Errorf("No such cluster exists with name - ", a.Clustername)
	}

	go func() {
		Provisioners[provisioner].Destroy(nil)
	}()

	// Delete AWS cluster config file
	err := a.DeleteConfig()
	if err != nil {
		return fmt.Errorf("Error deleting cluster config ...")
	}

	return nil
}

func (a *Aws) CreateConfig() error {
	viper.New()
	viper.SetConfigType("yaml")

	configFileName := "aws-" + a.Clustername + ".yaml"

	viper.SetConfigFile(configFileName)
	viper.AddConfigPath(common.REST_API_STORAGEPATH)

	viper.Set("aws.clustername", a.Clustername)
	viper.Set("aws.os", a.Os)
	viper.Set("aws.aws_access_key_id", a.AwsAccessKeyID)
	viper.Set("aws.aws_secret_access_key", a.AwsSecretAccessKey)
	viper.Set("aws.aws_ssh_keypair", a.AwsSSHKeypair)
	viper.Set("aws.aws_default_region", a.AwsDefaultRegion)
	viper.Set("aws.aws_vpc_cidr_block", a.AwsVpcCidrBlock)
	viper.Set("aws.aws_cidr_subnets_private", a.AwsCidrSubnetsPrivate)
	viper.Set("aws.aws_cidr_subnets_public", a.AwsCidrSubnetsPublic)
	viper.Set("aws.aws_bastion_size", a.AwsBastionSize)
	viper.Set("aws.aws_kube_master_num", a.AwsKubeMasterNum)
	viper.Set("aws.aws_kube_master_size", a.AwsKubeMasterSize)
	viper.Set("aws.aws_etcd_num", a.AwsEtcdNum)
	viper.Set("aws.aws_etcd_size", a.AwsEtcdSize)
	viper.Set("aws.aws_kube_worker_num", a.AwsKubeWorkerNum)
	viper.Set("aws.aws_kube_worker_size", a.AwsKubeWorkerSize)
	viper.Set("aws.aws_elb_api_port", a.AwsElbAPIPort)
	viper.Set("aws.k8s_secure_api_port", a.K8SSecureAPIPort)
	viper.Set("aws.kube_insecure_apiserver_address", a.KubeInsecureApiserverAddress)
	viper.Set("aws.kubeadm_enabled", a.KubeadmEnabled)
	viper.Set("aws.kube_network_plugin", a.KubeNetworkPlugin)

	log.Println(viper.AllKeys())
	log.Println(viper.AllSettings())

	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func (a *Aws) ValidateConfig() error {
	if a.Clustername == "" {
		return fmt.Errorf("Cluster name cannot be empty")
	}

	configFileName := "aws-" + a.Clustername + ".yaml"
	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)

	if isExistsClusterConfig(configFileName) {
		return fmt.Errorf("Cluster name must be unique")
	}

	return nil
}

func (a *Aws) DeleteConfig() error {

	configFileName := "aws-" + a.Clustername + ".yaml"
	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)

	err := os.Remove(configFileName)
	if err != nil {
		return err
	}
	return nil
}

func (a *Aws) UpdateConfig() error {
	return nil
}
