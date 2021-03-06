package prometheussystemdexporter

import (
	"github.com/caos/orbos/internal/operator/boom/api/v1beta2"
	"github.com/caos/orbos/internal/operator/boom/api/v1beta2/k8s"
	"github.com/caos/orbos/internal/operator/boom/application/applications/prometheussystemdexporter/yaml"
	"github.com/caos/orbos/mntr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// var _ application.YAMLApplication = (*prometheusSystemdExporter)(nil)

func (*prometheusSystemdExporter) GetYaml(_ mntr.Monitor, toolset *v1beta2.ToolsetSpec) interface{} {

	resources := &k8s.Resources{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("256Mi"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("50m"),
			corev1.ResourceMemory: resource.MustParse("64Mi"),
		},
	}

	spec := toolset.SystemdMetricsExporter
	if spec.Resources != nil {
		resources = spec.Resources
	}

	return yaml.Build(resources)
}
