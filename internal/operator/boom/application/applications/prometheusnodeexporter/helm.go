package prometheusnodeexporter

import (
	toolsetslatest "github.com/caos/orbos/internal/operator/boom/api/latest"
	"github.com/caos/orbos/internal/operator/boom/application/applications/prometheusnodeexporter/helm"
	"github.com/caos/orbos/internal/operator/boom/templator/helm/chart"
	"github.com/caos/orbos/mntr"
)

func (p *PrometheusNodeExporter) SpecToHelmValues(monitor mntr.Monitor, toolset *toolsetslatest.ToolsetSpec) interface{} {
	// spec := toolset.PrometheusNodeExporter
	values := helm.DefaultValues(p.GetImageTags())

	// if spec.ReplicaCount != 0 {
	// 	values.ReplicaCount = spec.ReplicaCount
	// }

	spec := toolset.NodeMetricsExporter

	if spec == nil {
		return values
	}

	if spec.Resources != nil {
		values.Resources = spec.Resources
	}

	return values
}

func (p *PrometheusNodeExporter) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (c *PrometheusNodeExporter) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
