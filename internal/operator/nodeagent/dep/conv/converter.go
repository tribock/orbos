package conv

import (
	"github.com/caos/orbos/internal/operator/nodeagent/dep/health"
	"github.com/pkg/errors"

	"github.com/caos/orbos/internal/operator/common"
	"github.com/caos/orbos/internal/operator/nodeagent"
	"github.com/caos/orbos/internal/operator/nodeagent/dep"
	"github.com/caos/orbos/internal/operator/nodeagent/dep/cri"
	"github.com/caos/orbos/internal/operator/nodeagent/dep/hostname"
	"github.com/caos/orbos/internal/operator/nodeagent/dep/k8s/kubeadm"
	"github.com/caos/orbos/internal/operator/nodeagent/dep/k8s/kubectl"
	"github.com/caos/orbos/internal/operator/nodeagent/dep/k8s/kubelet"
	"github.com/caos/orbos/internal/operator/nodeagent/dep/keepalived"
	"github.com/caos/orbos/internal/operator/nodeagent/dep/middleware"
	"github.com/caos/orbos/internal/operator/nodeagent/dep/nginx"
	"github.com/caos/orbos/internal/operator/nodeagent/dep/sshd"
	"github.com/caos/orbos/internal/operator/nodeagent/dep/swap"
	"github.com/caos/orbos/internal/operator/nodeagent/dep/sysctl"
	"github.com/caos/orbos/mntr"
)

type Converter interface {
	Init() func() error
	nodeagent.Converter
}

type dependencies struct {
	monitor mntr.Monitor
	os      dep.OperatingSystemMajor
	pm      *dep.PackageManager
	sysd    *dep.SystemD
	cipher  string
}

func New(monitor mntr.Monitor, os dep.OperatingSystemMajor, cipher string) Converter {
	return &dependencies{monitor, os, nil, nil, cipher}
}

func (d *dependencies) Init() func() error {

	d.pm = dep.NewPackageManager(d.monitor, d.os.OperatingSystem)
	d.sysd = dep.NewSystemD(d.monitor)

	return func() error {
		if err := d.pm.Init(); err != nil {
			return err
		}
		return d.pm.RefreshInstalled()
	}
}

func (d *dependencies) ToDependencies(sw common.Software) []*nodeagent.Dependency {

	dependencies := []*nodeagent.Dependency{{
		Desired:   sw.Sysctl,
		Installer: sysctl.New(d.monitor),
	}, {
		Desired:   sw.Health,
		Installer: health.New(d.monitor, d.sysd),
	}, {
		Desired:   sw.Hostname,
		Installer: hostname.New(),
	}, {
		Desired:   sw.Swap,
		Installer: swap.New("/etc/fstab"),
	}, {
		Desired:   sw.KeepaliveD,
		Installer: keepalived.New(d.monitor, d.pm, d.sysd, d.os.OperatingSystem, d.cipher),
	}, {
		Desired:   sw.SSHD,
		Installer: sshd.New(d.sysd),
	}, {
		Desired:   sw.Nginx,
		Installer: nginx.New(d.monitor, d.pm, d.sysd, d.os.OperatingSystem),
	}, {
		Desired:   sw.Containerruntime,
		Installer: cri.New(d.monitor, d.os, d.pm, d.sysd),
	}, {
		Desired:   sw.Kubelet,
		Installer: kubelet.New(d.monitor, d.os.OperatingSystem, d.pm, d.sysd),
	}, {
		Desired:   sw.Kubectl,
		Installer: kubectl.New(d.os.OperatingSystem, d.pm),
	}, {
		Desired:   sw.Kubeadm,
		Installer: kubeadm.New(d.os.OperatingSystem, d.pm),
	},
	}

	for key, dependency := range dependencies {
		dependency.Installer = middleware.AddLogging(d.monitor, dependency.Installer)
		dependencies[key] = dependency
	}

	return dependencies
}

func (d *dependencies) ToSoftware(dependencies []*nodeagent.Dependency, pkg func(nodeagent.Dependency) common.Package) (sw common.Software) {

	for _, dependency := range dependencies {
		switch i := middleware.Unwrap(dependency.Installer).(type) {
		case sysctl.Installer:
			sw.Sysctl = pkg(*dependency)
		case health.Installer:
			sw.Health = pkg(*dependency)
		case hostname.Installer:
			sw.Hostname = pkg(*dependency)
		case swap.Installer:
			sw.Swap = pkg(*dependency)
		case kubelet.Installer:
			sw.Kubelet = pkg(*dependency)
		case kubeadm.Installer:
			sw.Kubeadm = pkg(*dependency)
		case kubectl.Installer:
			sw.Kubectl = pkg(*dependency)
		case cri.Installer:
			sw.Containerruntime = pkg(*dependency)
		case keepalived.Installer:
			sw.KeepaliveD = pkg(*dependency)
		case nginx.Installer:
			sw.Nginx = pkg(*dependency)
		case sshd.Installer:
			sw.SSHD = pkg(*dependency)
		default:
			panic(errors.Errorf("No installer type for dependency %s found", i))
		}
	}

	return sw
}
