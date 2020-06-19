# Logs 
 

 When the logs spec is nil all logs will get persisted in loki.


## Structure 
 

| Attribute                | Description                                                   | Default | Collection  |
| ------------------------ | ------------------------------------------------------------- | ------- | ----------  |
| ambassador               | Bool if logs will get persisted for ambassador                |         |             |
| grafana                  | Bool if logs will get persisted for grafana                   |         |             |
| argocd                   | Bool if logs will get persisted for argo-cd                   |         |             |
| kube-state-metrics       | Bool if logs will get persisted for kube-state-metrics        |         |             |
| prometheus-node-exporter | Bool if logs will get persisted for prometheus-node-exporter  |         |             |
| prometheus-operator      | Bool if logs will get persisted for prometheus-operator       |         |             |
| logging-operator         | Bool if logs will get persisted for logging-operator          |         |             |
| loki                     | Bool if logs will get persisted for loki                      |         |             |
| prometheus               | Bool if logs will get persisted for prometheus                |         |             |