package static

import (
	"github.com/caos/orbos/internal/operator/orbiter/kinds/clusters/core/infra"
	"github.com/caos/orbos/internal/operator/orbiter/kinds/providers/core"
)

var _ infra.Pool = (*infraPool)(nil)

type infraPool struct {
	pool        string
	machinesSvc core.MachinesService
}

func newInfraPool(pool string, machinesSvc core.MachinesService) *infraPool {
	return &infraPool{
		pool:        pool,
		machinesSvc: machinesSvc,
	}
}

func (i *infraPool) EnsureMember(infra.Machine) error {
	// Keepalived health checks should work
	return nil
}

func (i *infraPool) EnsureMembers() error {
	// Keepalived health checks should work
	return nil
}

func (i *infraPool) GetMachines() (infra.Machines, error) {
	return i.machinesSvc.List(i.pool)
}

func (i *infraPool) AddMachine() (infra.Machine, error) {
	return i.machinesSvc.Create(i.pool)
}
