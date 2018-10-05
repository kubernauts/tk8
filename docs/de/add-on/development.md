# TK8 Add-On Entwicklung

Die Add-On-Implementierung muss eine allgemeine Lösung sein, daher stellen wir ein cmd zur Verfügung, um ein Standard-Beispiel-Add-on zu liefern, das vom Benutzer angepasst werden kann.

## Entwicklungsbefehle

Der Befehl ein Addon zu erstellen ist folgender:

```shell
tk8 addon create my-addon
```

Dieser Befehl ruft die tk8-addon-Entwicklung von GitHub ab und erstellt einen neuen Ordner unter ./addons/my-addon.

Das Beispiel ist ein einfaches Nginx Deployment und ein LoadBalancer-Service, um darauf zuzugreifen. So kann der Benutzer, der dieses Add-on erstellt hat, es direkt verwenden und auf den k8s-Cluster anwenden.

```shell
tk8 addon install my-addon --kubeconfig /path/to/cluster/config
```

und kann mit dem folgenden Befehl wieder entfernt werden

```shell
tk8 addon destroy my-addon --kubeconfig /path/to/cluster/config
```

Das Standard-Entwickler-Add-on enthält keine main.sh-Datei. Aber wir müssen eine Dokumentation dafür erstellen. Unsere eigenen Add-ons könnten es nutzen und brauchen.

## TK8 Add-On-Struktur

Für die allgemeine Nutzung von Add-Ons mit git haben wir einen Standardrahmen definiert. Die die Ordnerstruktur, die yml-Struktur und ein Beispiel enthält.

Die Ordnerstruktur

→ addons

| → my-addon

| →  | → LICENCE

| →  | → Readme.md

| →  | → main.yml

| →  | → main.sh

Die main.yml enthält alle notwendigen Informationen für k8s und erstellt alle benötigten Deployments und Services.

Optional gibt es eine main.sh, mit der man externe Repositories herunterladen oder eine main.yml erstellen kann.