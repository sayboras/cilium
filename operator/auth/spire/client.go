// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package spire

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	entryv1 "github.com/spiffe/spire-api-sdk/proto/spire/api/server/entry/v1"
	"github.com/spiffe/spire-api-sdk/proto/spire/api/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/cilium/cilium/operator/auth/identity"
	"github.com/cilium/cilium/pkg/hive/cell"
)

var Cell = cell.Module(
	"spire-client",
	"Spire Server API Client",
	cell.Config(ClientConfig{}),
	cell.Provide(NewClient),
)

// ClientConfig contains the configuration for the SPIRE client.
type ClientConfig struct {
	AuthMTLSEnabled   bool   `mapstructure:"mesh-auth-mtls-enabled"`
	SpiffeTrustDomain string `mapstructure:"mesh-auth-spiffe-trust-domain"`
}

// Flags adds the flags used by ClientConfig.
func (cfg ClientConfig) Flags(flags *pflag.FlagSet) {
	flags.BoolVar(&cfg.AuthMTLSEnabled,
		"mesh-auth-mtls-enabled",
		false,
		"The flag to enable mTLS for the SPIRE server.")

	flags.StringVar(&cfg.SpiffeTrustDomain,
		"mesh-auth-spiffe-trust-domain",
		"spiffe.cilium.io",
		"The trust domain for the SPIFFE identity.")
}

var defaultSelectors = []*types.Selector{
	{
		Type:  "cilium",
		Value: "mtls",
	},
}

const defaultParentID = "/dclient"

type Client struct {
	cfg   ClientConfig
	entry entryv1.EntryClient
}

func NewClient(cfg ClientConfig) (identity.Provider, error) {
	if !cfg.AuthMTLSEnabled {
		return &noopClient{}, nil
	}

	// <hackedcode blame="maartje">
	// note this will hang till the socket and certificate is there, consider this in the future
	source, err := workloadapi.NewX509Source(context.TODO(), workloadapi.WithClientOptions(workloadapi.WithAddr("unix:///run/spire/sockets/agent/agent.sock")))
	if err != nil {
		return &noopClient{}, err
	}

	tlsConfig := tlsconfig.MTLSClientConfig(source, source, tlsconfig.AuthorizeAny())

	conn, err := grpc.Dial("spire-server.spire.svc.cluster.local:8081", grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return nil, err
	}

	// </hackedcode>

	return &Client{
		cfg:   cfg,
		entry: entryv1.NewEntryClient(conn),
	}, nil
}

func (c *Client) Upsert(ctx context.Context, id string) error {
	entry, err := c.entry.GetEntry(ctx, &entryv1.GetEntryRequest{Id: id})
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		return err
	}

	desired := []*types.Entry{
		{
			SpiffeId: &types.SPIFFEID{
				TrustDomain: c.cfg.SpiffeTrustDomain,
				Path:        fmt.Sprintf("/cilium-id/%s", id),
			},
			ParentId: &types.SPIFFEID{
				TrustDomain: c.cfg.SpiffeTrustDomain,
				Path:        defaultParentID,
			},
			Selectors: defaultSelectors,
		},
	}

	if entry == nil {
		_, err = c.entry.BatchCreateEntry(ctx, &entryv1.BatchCreateEntryRequest{Entries: desired})
		return err
	}

	_, err = c.entry.BatchUpdateEntry(ctx, &entryv1.BatchUpdateEntryRequest{
		Entries: desired,
	})
	return err
}

func (c *Client) Delete(ctx context.Context, id string) error {
	entries, err := c.entry.ListEntries(ctx, &entryv1.ListEntriesRequest{
		Filter: &entryv1.ListEntriesRequest_Filter{
			BySpiffeId: &types.SPIFFEID{
				TrustDomain: c.cfg.SpiffeTrustDomain,
				Path:        fmt.Sprintf("/cilium-id/%s", id),
			},
			ByParentId: &types.SPIFFEID{
				TrustDomain: c.cfg.SpiffeTrustDomain,
				Path:        defaultParentID,
			},
			BySelectors: &types.SelectorMatch{
				Selectors: defaultSelectors,
				Match:     types.SelectorMatch_MATCH_EXACT,
			},
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return nil
		}
		return err
	}

	var ids = make([]string, len(entries.Entries))
	for _, e := range entries.Entries {
		ids = append(ids, e.Id)
	}

	_, err = c.entry.BatchDeleteEntry(ctx, &entryv1.BatchDeleteEntryRequest{
		Ids: ids,
	})

	return err
}
