package grafana

import (
	"github.com/caos/orbos/internal/operator/boom/api/v1beta1/grafana/admin"
	"github.com/caos/orbos/internal/operator/boom/api/v1beta1/grafana/auth"
	"github.com/caos/orbos/internal/operator/boom/api/v1beta1/network"
	"github.com/caos/orbos/internal/operator/boom/api/v1beta1/storage"
	"reflect"
)

type Grafana struct {
	//Flag if tool should be deployed
	//@default: false
	Deploy bool `json:"deploy" yaml:"deploy"`
	//Spec for the definition of the admin account
	Admin *admin.Admin `json:"admin,omitempty" yaml:"admin,omitempty"`
	//Spec for additional datasources
	Datasources []*Datasource `json:"datasources,omitempty" yaml:"datasources,omitempty"`
	//Spec for additional Dashboardproviders
	DashboardProviders []*Provider `json:"dashboardproviders,omitempty" yaml:"dashboardproviders,omitempty"`
	//Spec to define how the persistence should be handled
	Storage *storage.Spec `json:"storage,omitempty" yaml:"storage,omitempty"`
	//Network configuration, used for SSO and external access
	Network *network.Network `json:"network,omitempty" yaml:"network,omitempty"`
	//Authorization and Authentication configuration for SSO
	Auth *auth.Auth `json:"auth,omitempty" yaml:"auth,omitempty"`
	//List of plugins which get added to the grafana instance
	Plugins []string `json:"plugins,omitempty" yaml:"plugins,omitempty"`
}

func (x *Grafana) MarshalYAML() (interface{}, error) {
	type Alias Grafana
	return &Alias{
		Deploy:             x.Deploy,
		Admin:              admin.ClearEmpty(x.Admin),
		Datasources:        x.Datasources,
		DashboardProviders: x.DashboardProviders,
		Storage:            x.Storage,
		Network:            x.Network,
		Auth:               auth.ClearEmpty(x.Auth),
	}, nil
}

func ClearEmpty(x *Grafana) *Grafana {
	if x == nil {
		return nil
	}

	marshaled := Grafana{
		Deploy:             x.Deploy,
		Admin:              admin.ClearEmpty(x.Admin),
		Datasources:        x.Datasources,
		DashboardProviders: x.DashboardProviders,
		Storage:            x.Storage,
		Network:            x.Network,
		Auth:               auth.ClearEmpty(x.Auth),
	}
	if reflect.DeepEqual(marshaled, Grafana{}) {
		return &Grafana{}
	}
	return &marshaled
}

type Datasource struct {
	//Name of the datasource
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	//Type of the datasource (for example prometheus)
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	//URL to the datasource
	Url string `json:"url,omitempty" yaml:"url,omitempty"`
	//Access defintion of the datasource
	Access string `json:"access,omitempty" yaml:"access,omitempty"`
	//Boolean if datasource should be used as default
	IsDefault bool `json:"isDefault,omitempty" yaml:"isDefault,omitempty"`
}

type Provider struct {
	//ConfigMaps in which the dashboards are stored
	ConfigMaps []string `json:"configMaps,omitempty" yaml:"configMaps,omitempty"`
	//Local folder in which the dashboards are mounted
	Folder string `json:"folder,omitempty" yaml:"folder,omitempty"`
}
