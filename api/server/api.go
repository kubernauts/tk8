package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kubernauts/tk8/internal/addon"

	aws "github.com/kubernauts/tk8-provisioner-aws"
	azure "github.com/kubernauts/tk8-provisioner-azure"
	baremetal "github.com/kubernauts/tk8-provisioner-baremetal"
	eks "github.com/kubernauts/tk8-provisioner-eks"
	nutanix "github.com/kubernauts/tk8-provisioner-nutanix"
	openstack "github.com/kubernauts/tk8-provisioner-openstack"
	rke "github.com/kubernauts/tk8-provisioner-rke"
	provisioner "github.com/kubernauts/tk8/pkg/provisioner"
)

const (
	// APIVersion for cluster APIs
	APIVersion = "v1"
)

type tk8Api struct {
	restBase
}

var Provisioners = map[string]provisioner.Provisioner{
	"aws":       aws.NewAWS(),
	"azure":     azure.NewAzure(),
	"baremetal": baremetal.NewBaremetal(),
	"eks":       eks.NewEKS(),
	"nutanix":   nutanix.NewNutanix(),
	"openstack": openstack.NewOpenstack(),
	"rke":       rke.NewRKE(),
}

func newTk8Api() restServer {
	return &tk8Api{
		restBase: restBase{
			version: APIVersion,
			name:    "TK8 API",
		},
	}
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

func (c *tk8Api) createAWSClusterHandler(w http.ResponseWriter, r *http.Request) {
	method := "createHandler"
	enableCors(&w)
	var aws Aws
	var err error
	if err = json.NewDecoder(r.Body).Decode(&aws); err != nil {
		fmt.Println("returning error here")
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}
	err = aws.CreateCluster()
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(aws)

}

func (c *tk8Api) createRKEClusterHandler(w http.ResponseWriter, r *http.Request) {
	method := "createHandler"
	enableCors(&w)
	var rke Rke
	var err error
	if err = json.NewDecoder(r.Body).Decode(&rke); err != nil {
		fmt.Println("returning error here")
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}
	err = rke.CreateCluster()
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(rke)

}

func (c *tk8Api) createEKSClusterHandler(w http.ResponseWriter, r *http.Request) {
	method := "createHandler"
	enableCors(&w)
	var eks Eks
	var err error
	if err = json.NewDecoder(r.Body).Decode(&eks); err != nil {
		fmt.Println("returning error here")
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}
	err = eks.CreateCluster()
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(eks)

}
func (c *tk8Api) getClusters(w http.ResponseWriter, r *http.Request) {
	//method := "getClusters"
	enableCors(&w)

	clusters := ReadClusterConfigs()
	if len(clusters) > 0 {
		json.NewEncoder(w).Encode(clusters)
		return
	}

	json.NewEncoder(w).Encode("No clusters found")

	// lists all the clusters that are created
	// by listing all the files created in the config directory

}

func (c *tk8Api) getCluster(w http.ResponseWriter, r *http.Request) {
	method := "getCluster"
	vars := mux.Vars(r)
	clusterID, ok := vars["id"]

	if !ok || clusterID == "" {
		c.sendError(c.name, method, w, "Missing id param", http.StatusBadRequest)
		return
	}
	// now that you have the name of the cluster that you want to get
	// checkout
}

func (c *tk8Api) getAWSClusterHandler(w http.ResponseWriter, r *http.Request) {
	method := "getAWSClusterHandler"
	vars := mux.Vars(r)
	clusterName, ok := vars["id"]

	if !ok || clusterName == "" {
		c.sendError(c.name, method, w, "Missing id param", http.StatusBadRequest)
		return
	}

	// check if cluster exists
	aws, err := decodeAwsClusterConfigToStruct(clusterName)
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(aws)
}

func (c *tk8Api) getEKSlusterHandler(w http.ResponseWriter, r *http.Request) {
	method := "getEKSClusterHandler"
	vars := mux.Vars(r)
	clusterName, ok := vars["id"]

	if !ok || clusterName == "" {
		c.sendError(c.name, method, w, "Missing id param", http.StatusBadRequest)
		return
	}

	// check if cluster exists
	eks, err := decodeEksClusterConfigToStruct(clusterName)
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(eks)
}

func (c *tk8Api) getRKEClusterHandler(w http.ResponseWriter, r *http.Request) {
	method := "getRKEClusterHandler"
	vars := mux.Vars(r)
	clusterName, ok := vars["id"]

	if !ok || clusterName == "" {
		c.sendError(c.name, method, w, "Missing id param", http.StatusBadRequest)
		return
	}

	// check if cluster exists
	rke, err := decodeRkeClusterConfigToStruct(clusterName)
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(rke)
}

func (c *tk8Api) destroyAWSClusterHandler(w http.ResponseWriter, r *http.Request) {
	method := "destroyAWSClusterHandler"
	vars := mux.Vars(r)
	clusterName, ok := vars["id"]

	if !ok {
		c.sendError(c.name, method, w, "Cluster name not passed", http.StatusBadRequest)
		return
	}
	var err error

	if clusterName == "" {
		//	c.sendError(c.name, method, w, "Cluster name cannot be empty", http.StatusBadRequest)
		//	return
	}

	// check if cluster exists
	aws, err := decodeAwsClusterConfigToStruct(clusterName)
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Coming here ..... >>>%+v", aws)

	err = aws.DestroyCluster()
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode("Cluster deletion started ...")
}

func (c *tk8Api) destroyEKSClusterHandler(w http.ResponseWriter, r *http.Request) {
	method := "destroyEKSClusterHandler"
	vars := mux.Vars(r)
	clusterName, ok := vars["id"]
	if !ok {
		c.sendError(c.name, method, w, "Cluster name cannot be empty", http.StatusBadRequest)
		return
	}
	var err error

	if clusterName == "" {
		c.sendError(c.name, method, w, "Cluster name cannot be empty", http.StatusBadRequest)
		return
	}

	// check if cluster exists
	eks, err := decodeEksClusterConfigToStruct(clusterName)
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}

	err = eks.DestroyCluster()
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode("Cluster deletion started ...")
}

func (c *tk8Api) destroyRKEClusterHandler(w http.ResponseWriter, r *http.Request) {
	method := "destroyRkeClusterHandler"
	vars := mux.Vars(r)
	clusterName, ok := vars["id"]
	if !ok {
		c.sendError(c.name, method, w, "Cluster name cannot be empty", http.StatusBadRequest)
		return
	}
	var err error

	if clusterName == "" {
		c.sendError(c.name, method, w, "Cluster name cannot be empty", http.StatusBadRequest)
		return
	}

	// check if cluster exists
	rke, err := decodeRkeClusterConfigToStruct(clusterName)
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rke.DestroyCluster()
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode("Cluster deletion started ...")
}

func (c *tk8Api) createInfraOnly(w http.ResponseWriter, req *http.Request) {
	config := req.ParseForm()
	Provisioners["aws"].Init(nil)
	json.NewEncoder(w).Encode(config)
}
