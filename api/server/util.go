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
