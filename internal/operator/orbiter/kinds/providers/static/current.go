package static

import (
	"github.com/caos/orbos/internal/operator/orbiter/kinds/clusters/core/infra"
	"github.com/caos/orbos/internal/operator/orbiter/kinds/providers/core"
	"github.com/caos/orbos/internal/tree"
)

func addPools(current *Current, desired *DesiredV0, machinesSvc core.MachinesService) error {
	current.Current.pools = make(map[string]infra.Pool)
	for pool := range desired.Spec.Pools {
		current.Current.pools[pool] = newInfraPool(pool, machinesSvc)
	}

	unconfiguredPools, err := machinesSvc.ListPools()
	if err != nil {
		return nil
	}
	for _, unconfiguredPool := range unconfiguredPools {
		if _, ok := current.Current.pools[unconfiguredPool]; !ok {
			current.Current.pools[unconfiguredPool] = newInfraPool(unconfiguredPool, machinesSvc)
		}
	}
	return nil
}

type Current struct {
	Common  *tree.Common `yaml:",inline"`
	Current struct {
		pools      map[string]infra.Pool `yaml:"-"`
		Ingresses  map[string]*infra.Address
		cleanupped <-chan error `yaml:"-"`
	}
}

func (c *Current) Pools() map[string]infra.Pool {
	return c.Current.pools
}
func (c *Current) Ingresses() map[string]*infra.Address {
	return c.Current.Ingresses
}
func (c *Current) Cleanupped() <-chan error {
	return c.Current.cleanupped
}
