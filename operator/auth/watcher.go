// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package auth

import (
	"context"

	"github.com/cilium/workerpool"
	"github.com/sirupsen/logrus"

	"github.com/cilium/cilium/operator/auth/identity"
	"github.com/cilium/cilium/pkg/hive"
	"github.com/cilium/cilium/pkg/hive/cell"
	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	"github.com/cilium/cilium/pkg/k8s/resource"
)

// params contains all the dependencies for the identity-gc.
// They will be provided through dependency injection.
type params struct {
	cell.In

	Logger         logrus.FieldLogger
	Lifecycle      hive.Lifecycle
	IdentityClient identity.Provider
	Identity       resource.Resource[*v2.CiliumIdentity]

	Cfg Config
}

// IdentityWatcher represents the Cilium identities watcher.
type IdentityWatcher struct {
	logger logrus.FieldLogger

	identityClient identity.Provider
	identity       resource.Resource[*v2.CiliumIdentity]
	wg             *workerpool.WorkerPool
	cfg            Config
}

func registerIdentityWatcher(p params) {
	if !p.Cfg.Enabled {
		return
	}
	iw := &IdentityWatcher{
		logger:         p.Logger,
		identityClient: p.IdentityClient,
		identity:       p.Identity,
		wg:             workerpool.New(1),
		cfg:            p.Cfg,
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
		var err error
		switch e.Kind {
		case resource.Upsert:
			err = iw.upsertIdentity(ctx, e.Object)
			iw.logger.WithError(err).WithField("identity", e.Object).Info("Upserted identity")
		case resource.Delete:
			err = iw.deleteIdentity(ctx, e.Object)
			iw.logger.WithError(err).WithField("identity", e.Object).Info("Deleted identity")
		}
		e.Done(err)
	}
	return nil
}

func (iw *IdentityWatcher) upsertIdentity(ctx context.Context, identity *v2.CiliumIdentity) error {
	return iw.identityClient.Upsert(ctx, identity.Name)
}

func (iw *IdentityWatcher) deleteIdentity(ctx context.Context, identity *v2.CiliumIdentity) error {
	return iw.identityClient.Delete(ctx, identity.Name)
}
