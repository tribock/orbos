package static

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/caos/orbos/internal/operator/orbiter/kinds/clusters/core/infra"
	"github.com/caos/orbos/internal/operator/orbiter/kinds/providers/core"
	"github.com/caos/orbos/mntr"
)

type machinesService struct {
	monitor           mntr.Monitor
	desired           *DesiredV0
	bootstrapKey      []byte
	maintenanceKey    []byte
	maintenanceKeyPub []byte
	statusFile        string
	onCreate          func(machine infra.Machine, pool string) error
	cache             map[string]cachedMachines
}

// TODO: Dont accept the whole spec. Accept exactly the values needed (check other constructors too)
func NewMachinesService(
	monitor mntr.Monitor,
	desired *DesiredV0,
	bootstrapKey []byte,
	maintenanceKey []byte,
	maintenanceKeyPub []byte,
	id string,
	onCreate func(machine infra.Machine, pool string) error) core.MachinesService {
	return &machinesService{
		monitor,
		desired,
		bootstrapKey,
		maintenanceKey,
		maintenanceKeyPub,
		filepath.Join("/var/orbiter", id),
		onCreate,
		nil,
	}
}

func (c *machinesService) ListPools() ([]string, error) {

	pools := make([]string, 0)

	for key := range c.desired.Spec.Pools {
		pools = append(pools, key)
	}

	return pools, nil
}

func (c *machinesService) List(poolName string) (infra.Machines, error) {
	pool, err := c.cachedPool(poolName)
	if err != nil {
		return nil, err
	}

	return pool.Machines(), nil
}

func (c *machinesService) Create(poolName string) (infra.Machine, error) {
	pool, err := c.cachedPool(poolName)
	if err != nil {
		return nil, err
	}

	for _, machine := range pool {

		if len(c.maintenanceKeyPub) == 0 {
			panic("no maintenance key")
		}
		if err := machine.WriteFile(c.desired.Spec.RemotePublicKeyPath, bytes.NewReader(c.maintenanceKeyPub), 600); err != nil {
			return nil, err
		}

		if !machine.active {

			if err := machine.WriteFile(c.statusFile, strings.NewReader("active"), 600); err != nil {
				return nil, err
			}

			if err := c.onCreate(machine, poolName); err != nil {
				return nil, err
			}

			machine.active = true
			return machine, nil
		}
	}

	return nil, errors.New("no machines left")
}

func (c *machinesService) cachedPool(poolName string) (cachedMachines, error) {

	specifiedMachines, ok := c.desired.Spec.Pools[poolName]
	if !ok {
		return nil, fmt.Errorf("pool %s does not exist", poolName)
	}

	cache, ok := c.cache[poolName]
	if ok {
		return cache, nil
	}

	newCache := make([]*machine, 0)
	for _, spec := range specifiedMachines {
		machine := newMachine(c.monitor, c.statusFile, c.desired.Spec.RemoteUser, &spec.ID, string(spec.IP))
		if err := machine.UseKey(c.maintenanceKey, c.bootstrapKey); err != nil {
			return nil, err
		}

		buf := new(bytes.Buffer)
		if err := machine.ReadFile(c.statusFile, buf); err != nil {
			// treat as inactive
		}
		machine.active = strings.Contains(buf.String(), "active")
		buf.Reset()
		newCache = append(newCache, machine)
	}

	if c.cache == nil {
		c.cache = make(map[string]cachedMachines)
	}
	c.cache[poolName] = newCache
	return newCache, nil
}

type cachedMachines []*machine

func (c cachedMachines) Machines() infra.Machines {
	machines := make([]infra.Machine, 0)
	for _, machine := range c {
		if machine.active {
			machines = append(machines, machine)
		}
	}
	return machines
}
