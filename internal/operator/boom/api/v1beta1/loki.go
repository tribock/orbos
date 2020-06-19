package v1beta1

import "github.com/caos/orbos/internal/operator/boom/api/v1beta1/storage"

type Loki struct {
	//Flag if tool should be deployed
	//@default: false
	Deploy bool `json:"deploy" yaml:"deploy"`
	//Spec to define which logs will get persisted
	//@default: nil
	Logs *Logs `json:"logs,omitempty" yaml:"logs,omitempty"`
	//Spec to define how the persistence should be handled
	//@default: nil
	Storage *storage.Spec `json:"storage,omitempty" yaml:"storage,omitempty"`
	//Flag if loki-output should be a clusteroutput instead a output crd
	//@default: false
	ClusterOutput bool `json:"clusterOutput,omitempty" yaml:"clusterOutput,omitempty"`
}

// Logs: When the logs spec is nil all logs will get persisted in loki.
type Logs struct {
	//Bool if logs will get persisted for ambassador
	Ambassador bool `json:"ambassador"`
	//Bool if logs will get persisted for grafana
	Grafana bool `json:"grafana"`
	//Bool if logs will get persisted for argo-cd
	Argocd bool `json:"argocd"`
	//Bool if logs will get persisted for kube-state-metrics
	KubeStateMetrics bool `json:"kube-state-metrics" yaml:"kube-state-metrics"`
	//Bool if logs will get persisted for prometheus-node-exporter
	PrometheusNodeExporter bool `json:"prometheus-node-exporter"  yaml:"prometheus-node-exporter"`
	//Bool if logs will get persisted for prometheus-operator
	PrometheusOperator bool `json:"prometheus-operator" yaml:"prometheus-operator"`
	//Bool if logs will get persisted for logging-operator
	LoggingOperator bool `json:"logging-operator" yaml:"logging-operator"`
	//Bool if logs will get persisted for loki
	Loki bool `json:"loki"`
	//Bool if logs will get persisted for prometheus
	Prometheus bool `json:"prometheus"`
}
