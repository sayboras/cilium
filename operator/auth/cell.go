// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package auth

import (
	"github.com/cilium/hive/cell"
	"github.com/spf13/pflag"

	"github.com/cilium/cilium/operator/auth/spire"
)

const (
	mutualAuthEnabled = "mesh-auth-mutual-enabled"
)

var Cell = cell.Module(
	"auth-identity",
	"Cilium Mutual Authentication Identity management",
	spire.Cell,
	cell.Config(defaultConfig),
	cell.Invoke(registerIdentityWatcher),
)

// Config contains the configuration for the identity-gc.
type Config struct {
	Enabled bool `mapstructure:"mesh-auth-mutual-enabled"`
}

// Flags implements cell.Flagger interface.
func (cfg Config) Flags(flags *pflag.FlagSet) {
	flags.Bool(mutualAuthEnabled, cfg.Enabled, "Enable mutual authentication in Cilium")
}

var defaultConfig = Config{
	Enabled: false,
}
