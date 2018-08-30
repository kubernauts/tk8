// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var ltaas, prom, heapster, rancher bool

// addonCmd represents the addon command
var addonCmd = &cobra.Command{
	Use:   "addon",
	Short: "Install kubernetes addon packages",
	Long: `
Install additional packages on top of your kubernetes deployment. Examples: Prometheus,
Zipkin, Kibana, Load Testing As A Service`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if ltaas {

			// Get kubeconfig file location
			fmt.Println("Please enter the path to your kubeconfig")
			var kubeConfig string
			fmt.Scanln(&kubeConfig)

			if _, err := os.Stat(kubeConfig); err != nil {
				fmt.Println("Kubeconfig not found, kindly check")
				os.Exit(1)
			}

			// check if kubectl is installed
			kerr, err := exec.LookPath("kubectl")
			if err != nil {
				log.Fatal("kubectl command not found, kindly check")
			}
			fmt.Printf("Found kubectl at %s\n", kerr)
			rr, err := exec.Command("kubectl", "--kubeconfig", kubeConfig, "version", "--short").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(string(rr))

			if _, err := os.Stat("./jmeter-kubernetes"); err == nil {
				fmt.Println("LTaaS repo already exist on this system")
			} else {

				err := exec.Command("git", "clone", "https://github.com/kubernauts/jmeter-kubernetes.git").Run()
				if err != nil {
					log.Fatalf("Seems there is a problem cloning the LTaaS repo, %v", err)
					os.Exit(1)
				}
			}

			var replicas, namespace string

			fmt.Println("How many slave replicas do you need?")
			fmt.Scanln(&replicas)

			sedValue := "s/replicas: [0-9]/replicas: " + replicas + "/g"
			setReplica := exec.Command("sed", "-i", sedValue, "./jmeter-kubernetes/jmeter_slaves_deploy.yaml")
			out1, err := setReplica.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(out1))
				os.Exit(1)
			} else {
				fmt.Println(string(out1))
			}
			// Create New Space
			fmt.Println("Please eneter the name of the namspace that should be created for the LTaaS")
			fmt.Scanln(&namespace)

			namespaceCreate := exec.Command("kubectl", "--kubeconfig", kubeConfig, "create", "namespace", namespace)

			out0, err := namespaceCreate.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(out0))
				os.Exit(1)
			} else {
				fmt.Println(string(out0))
			}

			// Create Slave deployments
			slaves := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", namespace, "create", "-f", "./jmeter-kubernetes/jmeter_slaves_deploy.yaml")
			out2, err := slaves.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(out2))
			} else {
				fmt.Println(string(out2))
			}

			// Slave service
			slaveSvc := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", namespace, "create", "-f", "./jmeter-kubernetes/jmeter_slaves_svc.yaml")
			out3, err := slaveSvc.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(out3))
			} else {
				fmt.Println(string(out3))
			}

			// Jmeter Master deployment

			masterCfgmap := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", namespace, "create", "-f", "./jmeter-kubernetes/jmeter_master_configmap.yaml")
			out4, err := masterCfgmap.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(out4))
			} else {
				fmt.Println(string(out4))
			}

			masterDeploy := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", namespace, "create", "-f", "./jmeter-kubernetes/jmeter_master_deploy.yaml")
			out5, err := masterDeploy.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(out5))
			} else {
				fmt.Println(string(out5))
			}

			// Jmeter Influxdb Deployment
			influxDBCfgmap := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", namespace, "create", "-f", "./jmeter-kubernetes/jmeter_influxdb_configmap.yaml")
			out6, err := influxDBCfgmap.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(out6))
			} else {
				fmt.Println(string(out4))
			}

			influxDBdeploy := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", namespace, "create", "-f", "./jmeter-kubernetes/jmeter_influxdb_deploy.yaml")
			out7, err := influxDBdeploy.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(out7))
			} else {
				fmt.Println(string(out5))
			}

			influxDBsvc := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", namespace, "create", "-f", "./jmeter-kubernetes/jmeter_influxdb_svc.yaml")
			out8, err := influxDBsvc.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(out8))
			} else {
				fmt.Println(string(out5))
			}

			// Jmeter Grafana Deployment

			grafanDeploy := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", namespace, "create", "-f", "./jmeter-kubernetes/jmeter_grafana_deploy.yaml")
			out9, err := grafanDeploy.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(out9))
			} else {
				fmt.Println(string(out9))
			}

			grafanaSvc := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", namespace, "create", "-f", "./jmeter-kubernetes/jmeter_grafana_svc.yaml")
			out10, err := grafanaSvc.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(out10))
			} else {
				fmt.Println(string(out10))
			}

			// Print all the configured objects
			objectPrintout := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", namespace, "get", "all")
			stdout, _ := objectPrintout.StdoutPipe()
			objectPrintout.Start()
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(m)
				//log.Printf(m)
			}

			objectPrintout.Wait()

			os.Exit(0)

		}

		if prom {

			// Get kubeconfig file location
			fmt.Println("Please enter the path to your kubeconfig")
			var kubeConfig string
			fmt.Scanln(&kubeConfig)

			if _, err := os.Stat(kubeConfig); err != nil {
				fmt.Println("Kubeconfig not found, kindly check")
				os.Exit(1)
			}

			// check if kubectl is installed
			kerr, err := exec.LookPath("kubectl")
			if err != nil {
				log.Fatal("kubectl command not found, kindly check")
			}
			fmt.Printf("Found kubectl at %s\n", kerr)
			rr, err := exec.Command("kubectl", "--kubeconfig", kubeConfig, "version", "--short").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(string(rr))
			_, err1 := http.Get("https://raw.githubusercontent.com/kubernauts/tk8/master/config-map.yaml")
			if err1 != nil {
				print(err1.Error())
			} else {
				fmt.Println("Deploying prometheus configmap")
				PromConfigMap := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", "kube-system", "create", "-f", "https://raw.githubusercontent.com/kubernauts/tk8/master/config-map.yaml")
				stdout, _ := PromConfigMap.StdoutPipe()
				PromConfigMap.Start()
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Println(m)
					//log.Printf(m)
				}

				PromConfigMap.Wait()
			}

			_, err2 := http.Get("https://raw.githubusercontent.com/kubernauts/tk8/master/prometheus-deployment.yaml")
			if err2 != nil {
				print(err2.Error())
			} else {
				fmt.Println("Creating prometheus deployment")
				PromDeploy := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", "kube-system", "create", "-f", "https://raw.githubusercontent.com/kubernauts/tk8/master/prometheus-deployment.yaml")
				stdout, _ := PromDeploy.StdoutPipe()
				PromDeploy.Start()
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Println(m)
					//log.Printf(m)
				}

				PromDeploy.Wait()
			}

			_, err3 := http.Get("https://raw.githubusercontent.com/kubernauts/tk8/master/prometheus-service.yaml")
			if err3 != nil {
				print(err3.Error())
			} else {
				fmt.Println("Deploying prometheus service")
				PromSvc := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", "kube-system", "create", "-f", "https://raw.githubusercontent.com/kubernauts/tk8/master/prometheus-service.yaml")
				stdout, _ := PromSvc.StdoutPipe()
				PromSvc.Start()
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Println(m)
					//log.Printf(m)
				}

				PromSvc.Wait()

				_, err4 := http.Get("https://raw.githubusercontent.com/kubernauts/tk8/master/grafana-deploy.yaml")
				if err4 != nil {
					print(err4.Error())
				} else {
					fmt.Println("Creating Grafana Deployment")
					GrafanaDeploy := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", "kube-system", "create", "-f", "https://raw.githubusercontent.com/kubernauts/tk8/master/grafana-deploy.yaml")
					stdout, _ := GrafanaDeploy.StdoutPipe()
					GrafanaDeploy.Start()
					scanner := bufio.NewScanner(stdout)
					for scanner.Scan() {
						m := scanner.Text()
						fmt.Println(m)
						//log.Printf(m)
					}

					GrafanaDeploy.Wait()
				}
			}

			_, err5 := http.Get("https://github.com/kubernauts/tk8/blob/master/grafana-svc.yaml")
			if err5 != nil {
				print(err5.Error())
			} else {
				fmt.Println("Deploying Grafana Service")
				GrafanaSvc := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", "kube-system", "create", "-f", "https://github.com/kubernauts/tk8/blob/master/grafana-svc.yaml")
				stdout, _ := GrafanaSvc.StdoutPipe()
				GrafanaSvc.Start()
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Println(m)
					//log.Printf(m)
				}

				GrafanaSvc.Wait()

				os.Exit(0)
			}
		}

		if heapster {

			// Get kubeconfig file location
			fmt.Println("Please enter the path to your kubeconfig")
			var kubeConfig string
			fmt.Scanln(&kubeConfig)

			if _, err := os.Stat(kubeConfig); err != nil {
				fmt.Println("Kubeconfig not found, kindly check")
				os.Exit(1)
			}

			// check if kubectl is installed
			kerr, err := exec.LookPath("kubectl")
			if err != nil {
				log.Fatal("kubectl command not found, kindly check")
			}
			fmt.Printf("Found kubectl at %s\n", kerr)
			rr, err := exec.Command("kubectl", "--kubeconfig", kubeConfig, "version", "--short").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(string(rr))

			if _, err := os.Stat("./prometheus-operator"); err == nil {
				fmt.Println("kube-prometheus clone on this system already exists")
			} else {
				fmt.Println("Cloning the kube-prometheus repo")
				err := exec.Command("git", "clone", "https://github.com/coreos/prometheus-operator.git").Run()
				if err != nil {
					log.Fatalf("Seems there is a problem cloning the kube-prometheus repo, %v", err)
					os.Exit(1)
				}

			}

			setConfig := "./" + kubeConfig
			os.Setenv("KUBECONFIG", setConfig)
			fmt.Println(os.Getenv("KUBECONFIG"))
			//exec.Command("export", setConfig)
			promSet := exec.Command("bash", "-c", "hack/cluster-monitoring/deploy")
			promSet.Dir = "./prometheus-operator/contrib/kube-prometheus/"
			_, er := promSet.Output()
			if er != nil {
				log.Fatal(err)
			}

			os.Exit(0)
		}
		if rancher {
			/* This is to install the Rancher addon where all the k8s objects
            for this are provided with main.yml
            This is applied with the kubectl create -f command
            */
            //get the kubeconfig file full path
			var kubeConfig = getKubeConfig()
            //Check if kubectl cmd is installed
            checkKubectl(kubeConfig)
            //Check if yaml for Rancher is present relative to current directory
            pwd, _ := os.Getwd()
            _, err1 := os.Stat(pwd + "/addons/rancher/master.yaml")
			if err1 != nil {
				print(err1.Error())
			} else {
				fmt.Println("Deploying Rancher")
				RancherDeploy := exec.Command("kubectl", "--kubeconfig", kubeConfig, "create", "-f", pwd + "/addons/rancher/master.yaml")
				stdout, _ := RancherDeploy.StdoutPipe()
				RancherDeploy.Start()
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Println(m)
				}
				RancherDeploy.Wait()
			}
		}

		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

	},
}

func getKubeConfig() string {
    /* This function gets the path to the kubeconfig, cluster details and auth
    for using with the kubectl.
    Then use this to install the addon on this cluster
    */
    fmt.Println("Please enter the path to your kubeconfig:")
    var kubeConfig string
    fmt.Scanln(&kubeConfig)
    fmt.Printf("path: %s\n", kubeConfig)
    if _, err := os.Stat(kubeConfig); err != nil {
        fmt.Println("Kubeconfig file not found, kindly check")
        os.Exit(1)
    }
    return kubeConfig
}

func checkKubectl(kubeConfig string) {
    /*This function is used to check the whether kubectl command is installed &
    it works with the kubeConfig provided
    */
    kerr, err := exec.LookPath("kubectl")
    if err != nil {
        log.Fatal("kubectl command not found, kindly check")
    }
    fmt.Printf("Found kubectl at %s\n", kerr)
    rr, err := exec.Command("kubectl", "--kubeconfig", kubeConfig, "version", "--short").Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf(string(rr))
}

func init() {
	rootCmd.AddCommand(addonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addonCmd.Flags().BoolVarP(&ltaas, "ltaas", "l", false, "Deploy Load Testing As A Service")
	addonCmd.Flags().BoolVarP(&prom, "prom", "p", false, "Deploy prometheus")
	addonCmd.Flags().BoolVarP(&heapster, "heapster", "m", false, "Deploy Heapster")
	addonCmd.Flags().BoolVarP(&rancher, "rancher", "r", false, "Deploy Rancher")
}
