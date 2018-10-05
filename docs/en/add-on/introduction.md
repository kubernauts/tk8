# TK8 Add-On

## Verwendung

Wie man ein Add-on mit tk8 installiert
Wir müssen jedes Add-On in einem einzigen Repository auf Github für öffentliche Add-Ons und auf gitlab für interne Add-Ons speichern.

Wir machen einen Switch mit installierten tk8 Add-ons und stellen eine Verknüpfung zur Verfügung. Eine Verknüpfung könnte auch ein lokales Add-on sein, daher müssen wir zuerst prüfen, ob es im Ordner eines gibt. Wenn nicht, überprüfen Sie, ob es ein tk8-addon- auf GitHub gibt.

#### Den kompletten Pfad verwenden

```shell
tk8 addon install https://github.com/kubernauts/tk8-addon-rancher
tk8 addon install https://github.com/kubernauts/tk8-addon-prometheus
tk8 addon install https://github.com/kubernauts/tk8-addon-grafana
tk8 addon install https://github.com/kubernauts/tk8-addon-monitoring-stack
tk8 addon install https://github.com/kubernauts/tk8-addon-elk
tk8 addon install https://github.com/kubernauts/tk8-addon-...
tk8 addon install https://github.com/USERNAME/ADDON-REPO
```

#### Verwenden Sie die Shortcuts

```shell
tk8 addon install rancher
tk8 addon install prometheus
tk8 addon install grafana
tk8 addon install monitoring-stack
tk8 addon install elk
```

### Add-on entfernen

```shell
tk8 addon destroy rancher
tk8 addon destroy prometheus
tk8 addon destroye grafana
tk8 addon destroy monitoring-stack
tk8 addon destroy elk
```

## Entwicklung

### Erstellen eines Add-ons

Die create-Methode von tk8 erstellt ein neues Add-on im lokalen Ordner. Dieses Add-on ist ein einfaches Beispiel und bietet alles, was wir brauchen, um mit diesem Add-on zu arbeiten.

[Weitere Informationen hier](development.md)

```shell
tk8 addon erstellen my-addon erstellen
```