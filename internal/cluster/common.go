package cluster

import (
	"fmt"
	"os"

	"github.com/alecthomas/template"
	"github.com/spf13/viper"
)

// AwsCredentials defines the structure to hold AWS auth credentials.
type AwsCredentials struct {
	AwsAccessKeyID   string
	AwsSecretKey     string
	AwsAccessSSHKey  string
	AwsDefaultRegion string
}

// DistOS defines the structure to hold the dist OS informations.
// It is possible to easily extend the list of OS.
// Append new entry to cluster.DistOSMap and use the key(string) in the config.
type DistOS struct {
	User     string
	AmiOwner string
	OS       string
}

// DistOSMap holds the main OS distrubution mapping informations.
var DistOSMap = map[string]DistOS{
	"centos": DistOS{
		User:     "centos",
		AmiOwner: "688023202711",
		OS:       "dcos-centos7-*",
	},
	"ubuntu": DistOS{
		User:     "ubuntu",
		AmiOwner: "099720109477",
		OS:       "ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-*",
	},
	"coreos": DistOS{
		User:     "core",
		AmiOwner: "595879546273",
		OS:       "CoreOS-stable-*",
	},
	//"debian": DistOS{
	//	User:     "admin",
	//	AmiOwner: "379101102735",
	//	OS:       "debian-jessie-amd64-hvm-*",
	//},
	//"opensuse": DistOS{
	//	User:     "ec2-user",
	//	AmiOwner: "056126556840",
	//	OS:       "opensuse/openSUSE-42.3-x86_64-*",
	//},
}

type ClusterConfig struct {
	AwsClusterName               string
	AwsVpcCidrBlock              string
	AwsCidrSubnetsPrivate        string
	AwsCidrSubnetsPublic         string
	AwsBastionSize               string
	AwsKubeMasterNum             string
	AwsKubeMasterSize            string
	AwsEtcdNum                   string
	AwsEtcdSize                  string
	AwsKubeWorkerNum             string
	AwsKubeWorkerSize            string
	AwsElbAPIPort                string
	K8sSecureAPIPort             string
	KubeInsecureApiserverAddress string
}

// ReadViperConfigFile is define the config paths and read the configuration file.
func ReadViperConfigFile(configName string) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/tk8")
	viper.AddConfigPath("./../..")
	verr := viper.ReadInConfig() // Find and read the config file.
	if verr != nil {             // Handle errors reading the config file.
		CreateConfig()
	}
}

// GetDistConfig is used to get config details specific to a particular distribution.
// Used to determine various details such as the SSH user about the distribution.
func GetDistConfig() (string, string, string) {
	ReadViperConfigFile("config")
	awsAmiID := viper.GetString("aws.ami_id")
	awsInstanceOS := viper.GetString("aws.os")
	sshUser := viper.GetString("aws.ssh_user")
	return awsAmiID, awsInstanceOS, sshUser
}

// GetCredentials get the aws credentials from the config file.
func GetCredentials() AwsCredentials {
	ReadViperConfigFile("config")
	return AwsCredentials{
		AwsAccessKeyID:   viper.GetString("aws.aws_access_key_id"),
		AwsSecretKey:     viper.GetString("aws.aws_secret_access_key"),
		AwsAccessSSHKey:  viper.GetString("aws.aws_ssh_keypair"),
		AwsDefaultRegion: viper.GetString("aws.aws_default_region"),
	}
}

// GetClusterConfig get the configuration from the config file.
func GetClusterConfig() ClusterConfig {
	ReadViperConfigFile("config")
	return ClusterConfig{
		AwsClusterName:               viper.GetString("aws.clustername"),
		AwsVpcCidrBlock:              viper.GetString("aws.aws_vpc_cidr_block"),
		AwsCidrSubnetsPrivate:        viper.GetString("aws.aws_cidr_subnets_private"),
		AwsCidrSubnetsPublic:         viper.GetString("aws.aws_cidr_subnets_public"),
		AwsBastionSize:               viper.GetString("aws.aws_bastion_size"),
		AwsKubeMasterNum:             viper.GetString("aws.aws_kube_master_num"),
		AwsKubeMasterSize:            viper.GetString("aws.aws_kube_master_size"),
		AwsEtcdNum:                   viper.GetString("aws.aws_etcd_num"),
		AwsEtcdSize:                  viper.GetString("aws.aws_etcd_size"),
		AwsKubeWorkerNum:             viper.GetString("aws.aws_kube_worker_num"),
		AwsKubeWorkerSize:            viper.GetString("aws.aws_kube_worker_size"),
		AwsElbAPIPort:                viper.GetString("aws.aws_elb_api_port"),
		K8sSecureAPIPort:             viper.GetString("aws.k8s_secure_api_port"),
		KubeInsecureApiserverAddress: viper.GetString("aws."),
	}
}

func parseTemplate(templateString string, outputFileName string, data interface{}) {
	// open template
	template := template.New("template")
	template, _ = template.Parse(templateString)
	// open output file
	outputFile, err := os.Create(GetFilePath(outputFileName))
	defer outputFile.Close()
	if err != nil {
		ExitErrorf("Error creating file %s: %v", outputFile, err)
	}
	err = template.Execute(outputFile, data)
	ErrorCheck("Error executing template: %v", err)

}

func ErrorCheck(msg string, err error) {
	if err != nil {
		ExitErrorf(msg, err)
	}
}

func ExitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
