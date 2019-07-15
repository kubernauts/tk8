package server

type AllClusters map[string][]Cluster

type Cluster interface {
	CreateCluster() error
	DestroyCluster() error
}

type Config interface {
	CreateConfig() error
	DeleteConfig() error
	UpdateConfig() error
	ValidateConfig() error
}

type AwsYaml struct {
	Aws *Aws `yaml:"aws"`
}
type RkeYaml struct {
	Rke *Rke `yaml:"rke"`
}
type EksYaml struct {
	Eks *Eks `yaml:"eks"`
}
type Rke struct {
	ClusterName         string `yaml:"cluster_name"`
	NodeOs              string `yaml:"node_os"`
	RkeAwsRegion        string `yaml:"rke_aws_region"`
	Authorization       string `yaml:"authorization"`
	RkeNodeInstanceType string `yaml:"rke_node_instance_type"`
	NodeCount           int    `yaml:"node_count"`
	CloudProvider       string `yaml:"cloud_provider"`
}

type Eks struct {
	ClusterName         string `yaml:"cluster-name"`
	AwsRegion           string `yaml:"aws_region"`
	NodeInstanceType    string `yaml:"node-instance-type"`
	DesiredCapacity     int    `yaml:"desired-capacity"`
	AutoscallingMaxSize int    `yaml:"autoscalling-max-size"`
	AutoscallingMinSize int    `yaml:"autoscalling-min-size"`
	KeyFilePath         string `yaml:"key-file-path"`
}

type Aws struct {
	Clustername                  string `yaml:"clustername"`
	Os                           string `yaml:"os"`
	AwsAccessKeyID               string `yaml:"aws_access_key_id"`
	AwsSecretAccessKey           string `yaml:"aws_secret_access_key"`
	AwsSSHKeypair                string `yaml:"aws_ssh_keypair"`
	AwsDefaultRegion             string `yaml:"aws_default_region"`
	AwsVpcCidrBlock              string `yaml:"aws_vpc_cidr_block"`
	AwsCidrSubnetsPrivate        string `yaml:"aws_cidr_subnets_private"`
	AwsCidrSubnetsPublic         string `yaml:"aws_cidr_subnets_public"`
	AwsBastionSize               string `yaml:"aws_bastion_size"`
	AwsKubeMasterNum             int    `yaml:"aws_kube_master_num"`
	AwsKubeMasterSize            string `yaml:"aws_kube_master_size"`
	AwsEtcdNum                   int    `yaml:"aws_etcd_num"`
	AwsEtcdSize                  string `yaml:"aws_etcd_size"`
	AwsKubeWorkerNum             int    `yaml:"aws_kube_worker_num"`
	AwsKubeWorkerSize            string `yaml:"aws_kube_worker_size"`
	AwsElbAPIPort                int    `yaml:"aws_elb_api_port"`
	K8SSecureAPIPort             int    `yaml:"k8s_secure_api_port"`
	KubeInsecureApiserverAddress string `yaml:"kube_insecure_apiserver_address"`
	KubeadmEnabled               bool   `yaml:"kubeadm_enabled"`
	KubeNetworkPlugin            string `yaml:"kube_network_plugin"`
}
