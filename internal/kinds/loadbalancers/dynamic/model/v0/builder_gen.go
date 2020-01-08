// Code generated by "gen-kindstubs -parentpath=github.com/caos/orbiter/internal/kinds/loadbalancers -versions=v0 -kind=orbiter.caos.ch/DynamicLoadBalancer from file gen.go"; DO NOT EDIT.

package v0

import (
	"errors"

	"github.com/caos/orbiter/internal/core/operator"
	"github.com/caos/orbiter/internal/kinds/loadbalancers/dynamic/model"
)

var build func(map[string]interface{}, *operator.Secrets, interface{}) (model.UserSpec, func(model.Config) ([]operator.Assembler, error))

func Build(spec map[string]interface{}, secrets *operator.Secrets, dependant interface{}) (model.UserSpec, func(cfg model.Config) ([]operator.Assembler, error)) {
	if build != nil {
		return build(spec, secrets, dependant)
	}
	return model.UserSpec{}, func(_ model.Config) ([]operator.Assembler, error) {
		return nil, errors.New("Version v0 for kind orbiter.caos.ch/DynamicLoadBalancer is not yet supported")
	}
}
