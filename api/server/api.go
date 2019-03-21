package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kubernauts/tk8/internal/addon"

	"github.com/kubernauts/tk8/pkg/common"
	"github.com/spf13/viper"

	aws "github.com/kubernauts/tk8-provisioner-aws"
	azure "github.com/kubernauts/tk8-provisioner-azure"
	baremetal "github.com/kubernauts/tk8-provisioner-baremetal"
	eks "github.com/kubernauts/tk8-provisioner-eks"
	nutanix "github.com/kubernauts/tk8-provisioner-nutanix"
	openstack "github.com/kubernauts/tk8-provisioner-openstack"
	rke "github.com/kubernauts/tk8-provisioner-rke"
)

const (
	// APIVersion for cluster APIs
	APIVersion = "v1"
)

type tk8Api struct {
	restBase
}

type Provisioner interface {
	Init(args []string)
	Setup(args []string)
	Scale(args []string)
	Remove(args []string)
	Reset(args []string)
	Upgrade(args []string)
	Destroy(args []string)
}

var Provisioners = map[string]Provisioner{
	"aws":       aws.NewAWS(),
	"azure":     azure.NewAzure(),
	"baremetal": baremetal.NewBaremetal(),
	"eks":       eks.NewEKS(),
	"nutanix":   nutanix.NewNutanix(),
	"openstack": openstack.NewOpenstack(),
	"rke":       rke.NewRKE(),
}

type Master struct {
	Count string `json:"3"`
	Size  string `json:"t3.large"`
}

type NodeConfig struct {
	Count int    `json:"count"`
	Size  string `json:"size"`
}
type Access struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type Config struct {
	Name             string     `json:"name"`
	Os               string     `json:"os"`
	Provisioner      string     `json:"provisioner"`
	Installer        string     `json:"installer"`
	Region           string     `json:"region"`
	Master           NodeConfig `json:"master"`
	Etcd             NodeConfig `json:"etcd"`
	Node             NodeConfig `json:"node"`
	Access           Access     `json:"access"`
	GenerateKeyPair  bool       `json:"generateKeyPair"`
	Cidr             string     `json:"cidr"`
	SubnetPrivate    string     `json:"subnetPrivate"`
	SubnetPublic     string     `json:"subnetPublic"`
	KeyPair          string     `json:"keyPair"`
	SSHPublicKeyPath string     `json:"sshPublicKeyPath"`
}

func newTk8Api() restServer {
	return &tk8Api{
		restBase: restBase{
			version: APIVersion,
			name:    "TK8 API",
		},
	}
}

func responseStatus(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func (c *tk8Api) installAddon(w http.ResponseWriter, r *http.Request) {
	method := "install"
	var addonReq AddonRequest
	var addon addon.Addon
	var err error

	if err = json.NewDecoder(r.Body).Decode(&addonReq); err != nil {

		fmt.Println("returning error here")
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}

	err = addon.Install(addonReq.Name, addonReq.Scope)
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &AddonResponse{Error: responseStatus(err)}
	if err == nil {
		resp.Status = "Successfully installed"
	}
	json.NewEncoder(w).Encode(&resp)
	//json.NewEncoder(w).Encode(cluster)
}

func (c *tk8Api) destroyAddon(w http.ResponseWriter, r *http.Request) {
	method := "destroy"
	var addonReq AddonRequest
	var addon addon.Addon
	var err error

	if err = json.NewDecoder(r.Body).Decode(&addonReq); err != nil {

		fmt.Println("returning error here")
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}

	err, _ = addon.Destroy(addonReq.Name, addonReq.Scope)
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &AddonResponse{Error: responseStatus(err)}
	if err == nil {
		resp.Status = "Successfully destroyed"
	}
	json.NewEncoder(w).Encode(&resp)
}

func (c *tk8Api) sendNotImplemented(w http.ResponseWriter, method string) {

	c.sendError(c.name, method, w, "Not implemented.", http.StatusNotImplemented)
}

func (c *tk8Api) createCluster(w http.ResponseWriter, r *http.Request) {
	method := "createHandler"
	enableCors(&w)
	var config Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		fmt.Println("returning error here")
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}
	go func() {
		createConfig(config)
		GetProvisioner(config.Provisioner)
		time.Sleep(10 * time.Second)
		Provisioners[config.Provisioner].Init(nil)
		Provisioners[config.Provisioner].Setup(nil)
	}()

	json.NewEncoder(w).Encode(config)
}

func (c *tk8Api) destroyCluster(w http.ResponseWriter, r *http.Request) {
	method := "destroy"
	var err error
	enableCors(&w)
	var config Config
	if err = json.NewDecoder(r.Body).Decode(&config); err != nil {
		fmt.Println("returning error here")
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}
	go func() {
		createConfig(config)
		GetProvisioner(config.Provisioner)
		time.Sleep(10 * time.Second)
		Provisioners[config.Provisioner].Destroy(nil)
	}()
	json.NewEncoder(w).Encode("Destroying cluster ... ")
}

func (c *tk8Api) createInfraOnly(w http.ResponseWriter, req *http.Request) {
	config := req.ParseForm()
	Provisioners["aws"].Init(nil)
	json.NewEncoder(w).Encode(config)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func createConfig(config Config) {
	switch prov := config.Provisioner; prov {
	case "eks":
		log.Println("prov eks")
		switch inst := config.Installer; inst {
		default:
			createConfigEKSTK8(config)
			log.Println("inst tk8")
		}
	default:
		log.Println("prov aws")
		switch inst := config.Installer; inst {
		default:
			log.Println("inst tk8")
			createConfigAWSTK8(config)
		}
	}
}

func createConfigEKSTK8(config Config) {
	viper.New()
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")

	viper.Set("eks.cluster-name", config.Name)
	viper.Set("eks.aws_region", config.Region)
	viper.Set("eks.node-instance-type", config.Node.Size)
	viper.Set("eks.desired-capacity", config.Node.Count)
	viper.Set("eks.autoscalling-max-size", config.Node.Count)
	viper.Set("eks.autoscalling-min-size", config.Node.Count)
	viper.Set("eks.key-file-path", "~/.ssh/id_rsa.pub")

	log.Println(viper.AllKeys())
	log.Println(viper.AllSettings())

	viper.WriteConfig()
}

func createConfigAWSRKE(config Config) {
	viper.New()
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")
}

func createConfigAWSTK8(config Config) {
	viper.New()
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")

	viper.Set("aws.clustername", config.Name)
	viper.Set("aws.os", config.Os)
	viper.Set("aws.aws_access_key_id", config.Access.Key)
	viper.Set("aws.aws_secret_access_key", config.Access.Secret)
	viper.Set("aws.aws_ssh_keypair", config.KeyPair)
	viper.Set("aws.aws_default_region", config.Region)
	viper.Set("aws.aws_vpc_cidr_block", config.Cidr)
	viper.Set("aws.aws_cidr_subnets_private", config.SubnetPrivate)
	viper.Set("aws.aws_cidr_subnets_public", config.SubnetPublic)
	viper.Set("aws.aws_bastion_size", config.Master.Size)
	viper.Set("aws.aws_kube_master_num", config.Master.Count)
	viper.Set("aws.aws_kube_master_size", config.Master.Size)
	viper.Set("aws.aws_etcd_num", config.Etcd.Count)
	viper.Set("aws.aws_etcd_size", config.Etcd.Size)
	viper.Set("aws.aws_kube_worker_num", config.Node.Count)
	viper.Set("aws.aws_kube_worker_size", config.Node.Size)
	viper.Set("aws.aws_elb_api_port", 6443)
	viper.Set("aws.k8s_secure_api_port", 6443)
	viper.Set("aws.kube_insecure_apiserver_address", "0.0.0.0")
	viper.Set("aws.kubeadm_enabled", false)
	viper.Set("aws.kube_network_plugin", "flannel")

	log.Println(viper.AllKeys())
	log.Println(viper.AllSettings())

	viper.WriteConfig()
}

func GetProvisioner(provisioner string) error {
	if _, ok := Provisioners[provisioner]; ok {
		if _, err := os.Stat("./provisioner/" + provisioner); err == nil {
			return nil
		}
		log.Println("get provisioner " + provisioner)
		os.Mkdir("./provisioner", 0755)
		common.CloneGit("./provisioner", "https://github.com/kubernauts/tk8-provisioner-"+provisioner, provisioner)
		common.ReplaceGit("./provisioner/" + provisioner)
		return nil

	}
	return errors.New("provisioner not supported")

}
