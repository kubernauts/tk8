package provisioner

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/kubernauts/tk8-provisioner-baremetal/internal/cluster"
	"github.com/spf13/viper"
)

type Baremetal struct {
}

func (p Baremetal) Init(args []string) {
	// Check for the load balancer configmap file in /tk8/baremetal

	if _, err := os.Stat("./baremetal/lb-config.yml"); err != nil {
		fmt.Println("The baremetal configmap does not exist, please check")
		os.Exit(1)
	}

	// Check for the kubeconfig file in /tk8/baremetal

	if _, err := os.Stat("./baremetal/kubeconfig"); err != nil {
		fmt.Println("The baremetal configmap does not exist, please check")
		os.Exit(1)
	}

	// Get kubeconfig file location
	fmt.Println("Please enter the path to your kubeconfig")
	var kubeConfig string
	fmt.Scanln(&kubeConfig)

	if _, err := os.Stat(kubeConfig); err != nil {
		fmt.Println("Kubeconfig not found, kindly check")
		os.Exit(1)
	}

	// check if kubectl is installed
	terr, err := exec.LookPath("kubectl")
	if err != nil {
		log.Fatal("Kubectl command not found, kindly check")
	}
	fmt.Printf("Found kubectl at %s\n", terr)
	rr, err := exec.Command("kubectl", "version", "--client=true").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(rr))

	// Deploy MetallB Artifacts i.e deployment, configmap, etc

	// Deploy the Controller, speaker, RBAC, service account
	kubeSet := exec.Command("kubectl", "--kubeconfig", kubeConfig, "create", "-f", "https://raw.githubusercontent.com/google/metallb/v0.5.0/manifests/metallb.yaml")
	kubeSet.Dir = "./baremetal/"
	stdout, _ := kubeSet.StdoutPipe()
	kubeSet.Stderr = kubeSet.Stdout
	kubeSet.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		//log.Printf(m)
	}

	kubeSet.Wait()

	// Deploy the Configmap
	ConfigSet := exec.Command("kubectl", "--kubeconfig", kubeConfig, "create", "-f", "lb-config.yml")
	ConfigSet.Dir = "./baremetal/"
	Configout, _ := ConfigSet.StdoutPipe()
	ConfigSet.Stderr = ConfigSet.Stdout
	ConfigSet.Start()
	Configscanner := bufio.NewScanner(Configout)
	for Configscanner.Scan() {
		m := Configscanner.Text()
		fmt.Println(m)
		//log.Printf(m)
	}

	ConfigSet.Wait()

	// Print out the objects under metallb namespace
	PrintSet := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", "metallb-system", "get", "all")
	PrintSet.Dir = "./baremetal/"
	Printout, _ := PrintSet.StdoutPipe()
	PrintSet.Stderr = ConfigSet.Stdout
	PrintSet.Start()
	Objectscanner := bufio.NewScanner(Printout)
	for Objectscanner.Scan() {
		m := Objectscanner.Text()
		fmt.Println(m)
		//log.Printf(m)
	}

	PrintSet.Wait()

	os.Exit(0)

}

func (p Baremetal) Setup(args []string) {
	// check if ansible is installed
	terr, err := exec.LookPath("ansible")
	if err != nil {
		log.Fatal("Ansible command not found, kindly check")
	}
	fmt.Printf("Found Ansible at %s\n", terr)
	rr, err := exec.Command("ansible", "--version").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(rr))

	//Start Kubernetes Installation

	//check if kubespray cluster folder for the (baremetal)exists

	// Copy sample-inventory as indicated in the kubespray docs

	if _, err := os.Stat("./kubespray/inventory/barecluster"); err == nil {
		fmt.Println("Inventory folder already exists")
	} else {
		exec.Command("cp", "-LRp", "./kubespray/inventory/sample", "./kubespray/inventory/barecluster").Run()

		//Make a copy of kubeconfig on Ansible host after installation
		f, err := os.OpenFile("./kubespray/inventory/barecluster/group_vars/k8s-cluster.yml", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()
		fmt.Fprintf(f, "kubeconfig_localhost: true\n")

	}

	//Check for the host.ini file if it exists

	if _, err := os.Stat("./baremetal/hosts.ini"); err != nil {
		fmt.Println("Kubespray hosts.ini does not exist in tk8/baremetal folder, please check")
		os.Exit(1)
	}

	// Copy the hosts.ini file from the baremetal folder to the inventory folder
	exec.Command("cp", "-rf", "./baremetal/hosts.ini", "./kubespray/inventory/barecluster/hosts.ini").Run()

	//Get the value of Operating system username and whether to become root user
	viper.SetConfigName("variables")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./baremetal/")
	venv := viper.ReadInConfig() // Find and read the config file
	if venv != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", venv))
	}

	OsUser := viper.GetString("os.username")
	AnsibleUser := "-e ansible_user=" + OsUser
	AnsibleBecome := viper.GetString("os.become")

	fmt.Print(OsUser, AnsibleBecome, AnsibleUser)

	if AnsibleBecome == "yes" {
		kubeSet := exec.Command("ansible-playbook", "-i", "./inventory/barecluster/hosts.ini", "./cluster.yml", AnsibleUser, "-b", "--become-user=root", "--flush-cache")
		kubeSet.Dir = "./kubespray/"
		stdout, _ := kubeSet.StdoutPipe()
		kubeSet.Stderr = kubeSet.Stdout
		kubeSet.Start()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
			//log.Printf(m)
		}

		kubeSet.Wait()

		os.Exit(0)
	}

	if AnsibleBecome == "no" {
		kubeSet := exec.Command("ansible-playbook", "-i", "./inventory/barecluster/hosts.ini", "./cluster.yml", AnsibleUser, "--flush-cache")
		kubeSet.Dir = "./kubespray/"
		stdout, _ := kubeSet.StdoutPipe()
		kubeSet.Stderr = kubeSet.Stdout
		kubeSet.Start()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
			//log.Printf(m)
		}

		kubeSet.Wait()

		os.Exit(0)
	}

}

func (p Baremetal) Scale(args []string) {
	cluster.NotImplemented()

}

func (p Baremetal) Reset(args []string) {
	cluster.NotImplemented()

}

func (p Baremetal) Remove(args []string) {
	cluster.NotImplemented()
}

func (p Baremetal) Upgrade(args []string) {
	cluster.NotImplemented()
}

func (p Baremetal) Destroy(args []string) {
	// check if ansible is installed
	terr, err := exec.LookPath("ansible")
	if err != nil {
		log.Fatal("Ansible command not found, kindly check")
	}
	fmt.Printf("Found Ansible at %s\n", terr)
	rr, err := exec.Command("ansible", "--version").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(rr))

	//Check for the host.ini file if it exists

	if _, err := os.Stat("./baremetal/hosts.ini"); err != nil {
		fmt.Println("Kubespray hosts.ini does not exist in tk8/baremetal folder, please check")
		os.Exit(1)
	}

	//Get the value of Operating system username and whether to become root user
	viper.SetConfigName("variables")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./baremetal/")
	venv := viper.ReadInConfig() // Find and read the config file
	if venv != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", venv))
	}

	OsUser := viper.GetString("os.username")
	AnsibleUser := "-e ansible_user=" + OsUser
	AnsibleBecome := viper.GetString("os.become")

	if AnsibleBecome == "yes" {
		kubeSet := exec.Command("ansible-playbook", "-i", "./inventory/barecluster/hosts.ini", "./reset.yml", AnsibleUser, "-b", "--become-user=root", "--flush-cache")
		kubeSet.Dir = "./kubespray/"
		stdout, _ := kubeSet.StdoutPipe()
		kubeSet.Stderr = kubeSet.Stdout
		kubeSet.Start()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
			//log.Printf(m)
		}

		kubeSet.Wait()

		// Remove the cluster inventory folder
		err = os.RemoveAll("./kubespray/inventory/barecluster")
		if err != nil {
			fmt.Println(err)
		}

		os.Exit(0)

	}

	if AnsibleBecome == "no" {
		kubeSet := exec.Command("ansible-playbook", "-i", "./inventory/barecluster/hosts.ini", "./reset.yml", AnsibleUser, "--flush-cache")
		kubeSet.Dir = "./kubespray/"
		stdout, _ := kubeSet.StdoutPipe()
		kubeSet.Stderr = kubeSet.Stdout
		kubeSet.Start()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
			//log.Printf(m)
		}

		kubeSet.Wait()

		// Remove the cluster inventory folder
		err = os.RemoveAll("./kubespray/inventory/barecluster")
		if err != nil {
			fmt.Println(err)
		}

		os.Exit(0)
	}
}

func NewBaremetal() cluster.Provisioner {
	provisioner := new(Baremetal)
	return provisioner
}
