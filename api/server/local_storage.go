package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubernauts/tk8/api"
	"github.com/kubernauts/tk8/pkg/common"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func NewLocalStore(name, path string) *LocalStore {
	return &LocalStore{
		FileName: name,
		FilePath: path,
	}
}

func (l *LocalStore) DeleteConfig() error {
	fullpath := filepath.Join(l.FilePath, l.FileName)
	err := os.Remove(fullpath)
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalStore) ValidateConfig() error {

	return nil
}

func (l *LocalStore) UpdateConfig() error {

	return nil
}

func (l *LocalStore) CheckConfigExists() (bool, error) {
	fullpath := filepath.Join(l.FilePath, l.FileName)
	if _, err := os.Stat(fullpath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
	}
	return true, nil
}

func (l *LocalStore) CreateConfig(t api.Cluster) error {
	viper.New()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(l.FileName)
	viper.AddConfigPath(l.FilePath)

	switch a := t.(type) {
	case *Aws:
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
	case *Eks:
		viper.Set("eks.cluster-name", a.ClusterName)
		viper.Set("eks.aws_region", a.AwsRegion)
		viper.Set("eks.node-instance-type", a.NodeInstanceType)
		viper.Set("eks.desired-capacity", a.DesiredCapacity)
		viper.Set("eks.autoscalling-max-size", a.AutoscallingMaxSize)
		viper.Set("eks.autoscalling-min-size", a.AutoscallingMinSize)
		viper.Set("eks.key-file-path", "~/.ssh/id_rsa.pub")
	case *Rke:
		viper.Set("rke.cluster_name", a.ClusterName)
		viper.Set("rke.node_os", a.NodeOs)
		viper.Set("rke.rke_aws_region", a.ClusterName)

		viper.Set("rke.authorization", a.ClusterName)
		viper.Set("rke.rke_node_instance_type", a.ClusterName)
		viper.Set("rke.node_count", a.ClusterName)
		viper.Set("rke.cloud_provider", a.CloudProvider)
	}

	log.Println(viper.AllKeys())
	log.Println(viper.AllSettings())

	fullPath := filepath.Join(l.FilePath, l.FileName)
	err := viper.WriteConfigAs(fullPath)
	if err != nil {
		return err
	}
	return nil
}
func ReadClusterConfigs() api.AllClusters {

	files, _ := ioutil.ReadDir(common.REST_API_STORAGEPATH)
	clusters := make(api.AllClusters)
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
			clusters["aws"] = append([]api.Cluster{awsConfig.Aws})

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
			clusters["eks"] = append([]api.Cluster{eksConfig.Eks})

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
			clusters["rke"] = append([]api.Cluster{rkeConfig.Rke})
		}
	}
	return clusters
}

func (l *LocalStore) GetConfig() ([]byte, error) {
	fullPath := filepath.Join(l.FilePath, l.FileName)
	yamlFile, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return nil, err
	}

	return yamlFile, nil
}

// func decodeAwsClusterConfigToStruct(name string) (*Aws, error) {
// 	configFileName := "aws-" + name + ".yaml"
// 	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)

// 	exists := isExistsClusterConfig(configFileName)

// 	if !exists {
// 		return nil, fmt.Errorf("No cluster found with the provided name ::: ", name)
// 	}

// 	awsConfig := &AwsYaml{}
// 	yamlFile, err := ioutil.ReadFile(configFileName)
// 	if err != nil {
// 		log.Printf("yamlFile.Get err   #%v ", err)
// 	}
// 	err = yaml.Unmarshal(yamlFile, awsConfig)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to decode into rke config struct, %v", err)
// 	}

// 	return awsConfig.Aws, nil
// }

// func decodeEksClusterConfigToStruct(name string) (*Eks, error) {
// 	configFileName := "eks-" + name + ".yaml"
// 	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)

// 	exists := isExistsClusterConfig(configFileName)

// 	if !exists {
// 		return nil, fmt.Errorf("No cluster found with the provided name ::: ", name)
// 	}
// 	eksConfig := &EksYaml{}
// 	yamlFile, err := ioutil.ReadFile(configFileName)
// 	if err != nil {
// 		log.Printf("yamlFile.Get err   #%v ", err)
// 	}
// 	err = yaml.Unmarshal(yamlFile, eksConfig)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to decode into rke config struct, %v", err)
// 	}

// 	return eksConfig.Eks, nil
// }

// func decodeRkeClusterConfigToStruct(name string) (*Rke, error) {

// 	configFileName := "rke-" + name + ".yaml"
// 	configFileName = filepath.Join(common.REST_API_STORAGEPATH, configFileName)

// 	exists := isExistsClusterConfig(configFileName)

// 	if !exists {
// 		return nil, fmt.Errorf("No cluster found with the provided name ::: ", name)
// 	}

// 	rkeConfig := &RkeYaml{}
// 	yamlFile, err := ioutil.ReadFile(configFileName)
// 	if err != nil {
// 		log.Printf("yamlFile.Get err   #%v ", err)
// 	}
// 	err = yaml.Unmarshal(yamlFile, rkeConfig)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to decode into rke config struct, %v", err)
// 	}

// 	return rkeConfig.Rke, nil
// }
