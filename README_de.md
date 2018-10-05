![Logo](docs/images/tk8.png)

# TK8: Ein Multi-Cloud, Multi-Cluster Kubernetes Plattform-Installations- und Integrationstool

TK8 ist ein in Go geschriebenes Kommandozeilenprogramm und installiert Kubernet vollautomatisch in jeder Umgebung. Mit TK8 können Sie verschiedene Kubernet-Cluster mit unterschiedlichen Konfigurationen zentral verwalten. Darüber hinaus bietet TK8 mit seiner einfachen Add-On-Integration die Möglichkeit, Erweiterungen schnell, sauber und einfach auf die verschiedenen Kubernetes-Cluster zu verteilen.

Dazu gehören ein Jmeter-Cluster für Lasttests, Prometheus für Monitoring, Jaeger, Linkerd oder Zippkin für Tracing, Ambassador API Gateway mit Envoy für Ingress und Load Balancing, Istio als Mesh-Support-Lösung, Jenkins-X für CI/CD-Integration. Darüber hinaus unterstützt das Add-On-System auch die Verwaltung von Helmpaketen.

## Inhaltsverzeichnis

Die Dokumentation sowie ein detailliertes Inhaltsverzeichnis finden Sie hier.

* [Inhaltsverzeichnis](docs/de/SUMMARY.md)

## Installation

Der TK8 CLI benötigt einige Abhängigkeiten, um seine Aufgaben zu erfüllen.
Im Moment brauchen wir hier noch Ihre Hilfe, aber wir arbeiten bereits an einem Setup-Skript, das diese Aufgaben für Sie erledigt.

### Terraform

Terraform ist erforderlich, um die Infrastruktur automatisch in der gewünschten Umgebung einzurichten.
[Terraform-Installation](https://www.terraform.io/intro/getting-started/install.html)

### Ansible

Es ist erforderlich, dass Ansible die automatisierten Installationsroutinen in der gewünschten und automatisch erstellten Umgebung ausführt.
[Ansible Installation](https://docs.ansible.com/ansible/2.5/installation_guide/intro_installation.html#Installation- der Steuerungsmaschine)

### Kubectl

Kubectl wird vom CLI für das Rollout der Add-ons und von Ihnen für den Zugriff auf Ihre Cluster benötigt.
[kubectl Installation](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

### Python und Pipeline

In den automatisierten Routinen werden Python-Skripte verwendet, dazu werden zusätzlich mit Pip Abhängigkeiten geladen.
[Python-Installation](https://www.python.org/downloads/)
[Rohrinstallation](https://pip.pypa.io/en/stable/installing/)

### AWS IAM Authenticator

Wenn Sie einen EKS-Cluster mit TK8 installieren möchten, muss der [AWS IAM Authenticator](https://github.com/kubernetes-sigs/aws-iam-authenticator) ausführbar sein _(/usr/local/bin)_. Diese ist im Provisioner-Paket EKS des TK8 CLI enthalten oder befindet sich unter dem angegebenen Link.

## Verwendung

Da es mit dem TK8 CLI verschiedene Zielplattformen gibt und wir diese in der Dokumentation separat ausführlich beschrieben haben, möchten wir Ihnen nur ein Beispiel mit AWS geben.

Laden Sie die ausführbare Datei für Ihr Betriebssystem aus dem Abschnitt Release herunter oder erstellen Sie Ihre eigene Version mit dem Befehl `go build`.

Erstellen Sie einen separaten Ordner und speichern Sie die ausführbare Datei dort, eine Konfigurationsdatei ist ebenfalls erforderlich. Diese Datei befindet sich unter dem Namen config.yaml.example. Geben Sie hier die notwendigen Parameter für Ihren Cluster sowie den AWS CLI Key und das Secret ein. Zusätzlich sollten Sie Ihre AWS-Anmeldeinformationen in die Umgebungsvariablen eintragen, da Teile des CLI (EKS-Cluster) sie dort benötigen.

`export AWS_SECRET_ACCESS_KEY=xxx`
`export AWS_ACCESS_KEY_ID=xxx`

Sie führen dann die CLI mit dem Befehl aus:
`tk8 cluster install aws`

Mit diesem Befehl erstellt der TK8 CLI alle benötigten Ressourcen in AWS und installiert dafür einen Kubernet-Cluster.

Wenn Sie den Cluster nicht mehr benötigen, können Sie den Befehl verwenden:
`tk8 cluster destroy aws`
um automatisch alle Ressourcen zu entfernen.

## Contributing

Für die Bereitstellung von Add-ons haben wir einen separaten Dokumentationsbereich und Beispiele, wie Sie Ihre Erweiterungen erstellen und in das TK8-Projekt integrieren können. Sie können uns auch unter Slack erreichen.

Als Plattformanbieter haben wir hier einen separaten Dokumentationsbereich, der sich nur mit der Integration einer Plattform in TK8 beschäftigt. Hier finden Sie detaillierte Anweisungen und Beispiele, wie TK8 Ihre Integration durchführen wird. Sie können uns auch im Nachhinein erreichen.

Um am Kern teilzunehmen, erstellen Sie bitte ein Problem oder kontaktieren Sie uns in Slack.

Nehmen Sie Kontakt auf
[Schließt euch uns auf Kubernauts Slack Channel an](https://kubernauts-slack-join.herokuapp.com/)

## Credits

Gründer und Initiator dieses Projekts ist[Arash Kaffamanesh](https://github.com/arashkaffamanesh) Gründer und CEO der [cloudssky GmbH](https://cloudssky.com/de/) und[Kubernauts GmbH](https://kubernauts.de/en/home/)

Unterstützt wird das Projekt von Cloud Computing-Experten der cloudssky GmbH und der Kubernauts GmbH.
[Christopher Adigun](https://github.com/infinitydon)
[Arush Salil](https://github.com/arush-sal)
[Manuel Müller](https://github.com/MuellerMH)
[Niki](https://github.com/niki-1905)
[Anoop](https://github.com/anoopl)

Ein großer Dank geht an die Mitwirkenden von [Kubespray](https://github.com/kubernetes-incubator/kubespray), deren großartige Arbeit wir als Grundlage für den Aufbau und die Installation von Kubernetes in der AWS Cloud verwenden.

Darüber hinaus möchten wir uns bei den Mitwirkenden von [kubeadm](https://github.com/kubernetes/kubernetes/tree/master/cmd/kubeadm) bedanken, das derzeit nicht nur Teil des Kubespray-Projekts, sondern auch des TK8 ist.

Ebenfalls ein großes Dankeschön an [Wesley Charles Blake](https://github.com/WesleyCharlesBlake), auf dessen Grundlage wir unsere EKS-Integration anbieten konnten.

