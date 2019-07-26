package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubernauts/tk8/api"
	//	"github.com/kubernauts/tk8/pkg/common"
	"github.com/sirupsen/logrus"

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

	logrus.Println(viper.AllKeys())
	logrus.Println(viper.AllSettings())

	fullPath := filepath.Join(l.FilePath, l.FileName)
	err := viper.WriteConfigAs(fullPath)
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalStore) GetConfigs() (api.AllClusters, error) {
	files, _ := ioutil.ReadDir(l.FilePath)
	clusters := make(api.AllClusters)
	for _, file := range files {
		switch {
		case strings.Contains(file.Name(), "aws-"):
			configFileName := filepath.Join(l.FilePath, file.Name())

			awsConfig := &AwsYaml{}
			yamlFile, err := ioutil.ReadFile(configFileName)
			if err != nil {
				logrus.Printf("yamlFile.Get err   #%v ", err)
			}
			err = yaml.Unmarshal(yamlFile, awsConfig)

			if err != nil {
				fmt.Printf("unable to decode into aws config struct, %v", err)
				continue
			}
			clusters["aws"] = append([]api.Cluster{awsConfig.Aws})

		case strings.Contains(file.Name(), "eks-"):
			configFileName := filepath.Join(l.FilePath, file.Name())
			eksConfig := &EksYaml{}

			yamlFile, err := ioutil.ReadFile(configFileName)
			if err != nil {
				logrus.Printf("yamlFile.Get err   #%v ", err)
			}
			err = yaml.Unmarshal(yamlFile, eksConfig)
			if err != nil {
				fmt.Printf("unable to decode into eks config struct, %v", err)
				continue
			}
			clusters["eks"] = append([]api.Cluster{eksConfig.Eks})

		case strings.Contains(file.Name(), "rke-"):
			configFileName := filepath.Join(l.FilePath, file.Name())
			rkeConfig := &RkeYaml{}
			yamlFile, err := ioutil.ReadFile(configFileName)
			if err != nil {
				logrus.Printf("yamlFile.Get err   #%v ", err)
			}
			err = yaml.Unmarshal(yamlFile, rkeConfig)
			if err != nil {
				fmt.Printf("unable to decode into rke config struct, %v", err)
				continue
			}
			clusters["rke"] = append([]api.Cluster{rkeConfig.Rke})
		}
	}
	return clusters, nil
}

func (l *LocalStore) GetConfig() ([]byte, error) {
	fullPath := filepath.Join(l.FilePath, l.FileName)
	yamlFile, err := ioutil.ReadFile(fullPath)
	if err != nil {
		logrus.Printf("yamlFile.Get err   #%v ", err)
		return nil, err
	}

	return yamlFile, nil
}
