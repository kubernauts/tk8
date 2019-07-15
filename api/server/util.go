package server

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/kubernauts/tk8/pkg/common"
)

func responseStatus(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// func createConfig(config Config) {
// 	switch prov := config.Provisioner; prov {
// 	case "eks":
// 		log.Println("prov eks")
// 		switch inst := config.Installer; inst {
// 		default:
// 			createConfigEKSTK8(config)
// 			log.Println("inst tk8")
// 		}
// 	default:
// 		log.Println("prov aws")
// 		switch inst := config.Installer; inst {
// 		default:
// 			log.Println("inst tk8")
// 			createConfigAWSTK8(config)
// 		}
// 	}
// }

// Creating config files logic has to change
// file name should be based on the name of the cluster
// check if the file already exists
// collated under a common place such as .tk8 under home dir

func getProvisioner(provisioner string) error {
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

func isExistsClusterConfig(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
