// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package auth

import (
	"context"

	"github.com/cilium/workerpool"
	"github.com/sirupsen/logrus"

	"github.com/cilium/cilium/pkg/hive"
	"github.com/cilium/cilium/pkg/hive/cell"
	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	k8sClient "github.com/cilium/cilium/pkg/k8s/client"
	"github.com/cilium/cilium/pkg/k8s/resource"
)

// params contains all the dependencies for the identity-gc.
// They will be provided through dependency injection.
type params struct {
	cell.In

	Logger    logrus.FieldLogger
	Lifecycle hive.Lifecycle
	Clientset k8sClient.Clientset
	Identity  resource.Resource[*v2.CiliumIdentity]

	Cfg Config
}

// IdentityWatcher represents the Cilium identities watcher.
type IdentityWatcher struct {
	logger logrus.FieldLogger

	clientset k8sClient.Clientset
	identity  resource.Resource[*v2.CiliumIdentity]
	wg        *workerpool.WorkerPool
	cfg       Config
}

func registerIdentityWatcher(p params) {
	if !p.Clientset.IsEnabled() || !p.Cfg.Enabled {
		return
	}
	iw := &IdentityWatcher{
		logger:    p.Logger,
		clientset: p.Clientset,
		identity:  p.Identity,
		wg:        workerpool.New(1),
		cfg:       p.Cfg,
	}
	p.Lifecycle.Append(hive.Hook{
		OnStart: func(ctx hive.HookContext) error {
			return iw.wg.Submit("identity-watcher", iw.run)
		},
		OnStop: func(ctx hive.HookContext) error {
			return iw.wg.Close()
		},
	})
}

func (iw *IdentityWatcher) run(ctx context.Context) error {
	for e := range iw.identity.Events(ctx) {
		// Doing nothing right now
		e.Done(nil)
	}
	return nil
}
