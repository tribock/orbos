// Code generated by "gen-kindstubs -parentpath=github.com/caos/orbiter/internal/kinds/loadbalancers/dynamic/kinds -versions=v1 -kind=orbiter.caos.ch/DynamicKubeAPILoadBalancer from file gen.go"; DO NOT EDIT.

package kubeapi

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/caos/orbiter/internal/core/operator"

	"github.com/caos/orbiter/internal/kinds/loadbalancers/dynamic/kinds/kubeapi/adapter"
	"github.com/caos/orbiter/internal/kinds/loadbalancers/dynamic/kinds/kubeapi/model"
	v1builder "github.com/caos/orbiter/internal/kinds/loadbalancers/dynamic/kinds/kubeapi/model/v1"
)

type Version int

const (
	unknown Version = iota
	v1
)

type assembler struct {
	path      []string
	overwrite func(map[string]interface{})
	builder   adapter.Builder
	built     adapter.Adapter
}

func New(configPath []string, overwrite func(map[string]interface{}), builder adapter.Builder) operator.Assembler {
	return &assembler{configPath, overwrite, builder, nil}
}

func (a *assembler) String() string { return "orbiter.caos.ch/DynamicKubeAPILoadBalancer" }
func (a *assembler) BuildContext() ([]string, func(map[string]interface{})) {
	return a.path, a.overwrite
}
func (a *assembler) Ensure(ctx context.Context, secrets *operator.Secrets, ensuredDependencies map[string]interface{}) (interface{}, error) {
	return a.built.Ensure(ctx, secrets, ensuredDependencies)
}
func (a *assembler) Build(serialized map[string]interface{}, nodeagentupdater operator.NodeAgentUpdater, secrets *operator.Secrets, dependant interface{}) (operator.Kind, interface{}, []operator.Assembler, error) {

	kind := operator.Kind{}
	if err := mapstructure.Decode(serialized, &kind); err != nil {
		return kind, nil, nil, err
	}

	if kind.Kind != "orbiter.caos.ch/DynamicKubeAPILoadBalancer" {
		return kind, nil, nil, fmt.Errorf("Kind %s must be \"orbiter.caos.ch/DynamicKubeAPILoadBalancer\"", kind.Kind)
	}

	var spec model.UserSpec
	var subassemblersBuilder func(model.Config) ([]operator.Assembler, error)
	switch kind.Version {
	case v1.String():
		spec, subassemblersBuilder = v1builder.Build(serialized, secrets, dependant)
	default:
		return kind, nil, nil, fmt.Errorf("Unknown version %s", kind.Version)
	}

	cfg, adapter, err := a.builder.Build(spec, nodeagentupdater)
	if err != nil {
		return kind, nil, nil, err
	}
	a.built = adapter

	if subassemblersBuilder == nil {
		return kind, cfg, nil, nil
	}

	subassemblers, err := subassemblersBuilder(cfg)
	if err != nil {
		return kind, nil, nil, err
	}

	return kind, cfg, subassemblers, nil
}
