package provisioner

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/CrowdSurge/banner"
	"github.com/blang/semver"
	"github.com/kubernauts/tk8/internal/cluster"
	"github.com/kubernauts/tk8/internal/templates"
)

type EKS struct {
}

var Name string

func (p EKS) Init(args []string) {
	Name = cluster.Name
	if len(Name) == 0 {
		Name = "TK8EKS"
	}
}

func (p EKS) Setup(args []string) {
	banner.Print("kubernauts eks cli")

	fmt.Println()

	fmt.Println()

	kube, err := exec.LookPath("kubectl")
	if err != nil {
		log.Fatal("kubectl not found, kindly check")
	}
	fmt.Printf("Found kubectl at %s\n", kube)
	rr, err := exec.Command("kubectl", "version", "--client", "--short").Output()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(rr))

	//Check if kubectl version is greater or equal to 1.10

	parts := strings.Split(string(rr), " ")

	KubeCtlVer := strings.Replace((parts[2]), "v", "", -1)

	v1, err := semver.Make("1.10.0")
	v2, err := semver.Make(strings.TrimSpace(KubeCtlVer))

	if v2.LT(v1) {
		log.Fatalln("kubectl client version on this system is less than the required version 1.10.0")
	}

	os.MkdirAll("./inventory/"+Name+"/provisioner", 0755)
	exec.Command("cp", "-rfp", "./provisioner/eks/", "./inventory/"+Name+"/provisioner").Run()
	cluster.ParseTemplate(templates.Credentials, "./inventory/"+Name+"/provisioner/credentials.tfvars", cluster.GetCredentials())
	cluster.ParseTemplate(templates.VariablesEKS, "./inventory/"+Name+"/provisioner/variables.tfvars", cluster.GetEKSConfig())

	// Check if AWS authenticator binary is present in the working directory
	if _, err := exec.LookPath("provisioner/eks/aws-iam-authenticator"); err != nil {
		log.Fatalln("AWS Authenticator binary not found")
	}

	// Check if terraform binary is present in the working directory
	if _, err := exec.LookPath("terraform"); err != nil {
		log.Fatalln("Terraform binary not found in the installation folder")
	}

	log.Println("Terraform binary exists in the installation folder, terraform version:")

	terr, err := exec.Command("terraform", "version").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(terr))

	// Check if a terraform state file aclready exists
	if _, err := os.Stat("./inventory/" + Name + "/provisioner/terraform.tfstate"); err == nil {
		log.Println("There is an existing cluster, please remove terraform.tfstate file or delete the installation before proceeding")
	} else {

		// Terraform Initialization and create the infrastructure

		log.Println("starting terraform init")

		terrInit := exec.Command("terraform", "init")
		terrInit.Dir = "./inventory/" + Name + "/provisioner/"
		out, _ := terrInit.StdoutPipe()
		terrInit.Start()
		scanInit := bufio.NewScanner(out)
		for scanInit.Scan() {
			m := scanInit.Text()
			fmt.Println(m)
		}

		terrInit.Wait()

	}
	log.Println("starting terraform apply")
	terrSet := exec.Command("terraform", "apply", "-var-file=credentials.tfvars", "-auto-approve")
	terrSet.Dir = "./inventory/" + Name + "/provisioner/"
	stdout, err := terrSet.StdoutPipe()
	terrSet.Stderr = terrSet.Stdout
	terrSet.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	terrSet.Wait()

	// Export KUBECONFIG file to the installation folder
	log.Println("Exporting kubeconfig file to the installation folder")

	kubeconf := exec.Command("terraform", "output", "./inventory/"+Name+"/provisioner/kubeconfig")
	kubeconf.Dir = "./inventory/" + Name + "/provisioner/"

	// open the out file for writing
	outfile, err := os.Create("./inventory/" + Name + "/provisioner/kubeconfig")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	kubeconf.Stdout = outfile

	err = kubeconf.Start()
	if err != nil {
		panic(err)
	}
	kubeconf.Wait()

	log.Println("To use the kubeconfig file, do the following:")

	log.Println("export KUBECONFIG=~/inventory/" + Name + "/provisioner/kubeconfig")

	// Output configmap to create Worker nodes

	log.Println("Exporting Worker nodes config-map to the installation folder")

	confmap := exec.Command("terraform", "output", "./inventory/"+Name+"/provisioner/config-map")
	confmap.Dir = "./inventory/" + Name + "/provisioner/"

	// open the out file for writing
	outconf, err := os.Create("./inventory/" + Name + "/provisioner/config-map-aws-auth.yaml")
	if err != nil {
		panic(err)
	}
	defer outconf.Close()
	confmap.Stdout = outconf

	err = confmap.Start()
	if err != nil {
		panic(err)
	}
	confmap.Wait()

	// Create Worker nodes usign the Configmap created above

	log.Println("Creating Worker Nodes")
	WorkerNodeSet := exec.Command("kubectl", "--kubeconfig", "./inventory/"+Name+"/provisioner/kubeconfig", "apply", "-f", "./inventory/"+Name+"/provisioner/config-map-aws-auth.yaml")
	WorkerNodeSet.Dir = "./"

	workerNodeOut, err := WorkerNodeSet.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(workerNodeOut))

	log.Println("Worker nodes are coming up one by one, it will take some time depending on the number of worker nodes you specified")

	os.Exit(0)
}

func (p EKS) Upgrade(args []string) {
	cluster.NotImplemented()
}

func (p EKS) Destroy(args []string) {
	if len(args) > 0 {
		fmt.Println()
		fmt.Println("Invalid, there is no need to use arguments with this command")
		fmt.Println()
		fmt.Println("Simple use : ekscluster delete")
		fmt.Println()
		os.Exit(0)
	}

	cluster.SetClusteName()
	p.Init(nil)

	banner.Print("kubernauts eks cli")

	fmt.Println()

	fmt.Println()

	// Check if terraform binary is present in the working directory
	if _, err := exec.LookPath("terraform"); err != nil {
		log.Fatalln("Terraform not found on the system")
	}

	log.Println("Terraform is installed, terraform version:")

	terr, err := exec.Command("terraform", "version").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(terr))

	// Check if a terraform state file already exists
	if _, err := os.Stat("./inventory/" + Name + "/provisioner/terraform.tfstate"); err != nil {
		log.Fatalln("Terraform.tfstate file not found, seems there is no existing cluster definition in this directory for the workspace" + Name)
	}

	// Terraform destroy the EKS cluster

	log.Println("starting terraform destroy")

	terrDel := exec.Command("terraform", "destroy", "-var-file=credentials.tfvars", "-force")
	terrDel.Dir = "./inventory/" + Name + "/provisioner"
	out, _ := terrDel.StdoutPipe()
	terrDel.Start()
	scanDel := bufio.NewScanner(out)
	for scanDel.Scan() {
		m := scanDel.Text()
		fmt.Println(m)
	}

	terrDel.Wait()

	// Delete terraform state file

	log.Println("Removing the terraform state file")

	// err = os.Remove("./inventory/" + Name + "/provisioner/terraform.tfstate")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func NewEKS() cluster.Provisioner {
	provisioner := new(EKS)
	return provisioner
}
