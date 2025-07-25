// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package k8s

import (
	"testing"

	"github.com/cilium/hive/hivetest"
	"github.com/stretchr/testify/require"

	cmtypes "github.com/cilium/cilium/pkg/clustermesh/types"
	slim_corev1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/core/v1"
	slim_discovery_v1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/discovery/v1"
	slim_metav1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
	"github.com/cilium/cilium/pkg/loadbalancer"
)

func TestEndpoints_DeepEqual(t *testing.T) {
	type fields struct {
		svcEP *Endpoints
	}
	type args struct {
		o *Endpoints
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{

		{
			name: "both equal",
			fields: fields{
				svcEP: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.1"): {
							Ports: map[loadbalancer.L4Addr][]string{
								{
									Protocol: loadbalancer.NONE,
									Port:     1,
								}: {"foo"},
							},
							NodeName: "k8s1",
						},
					},
				},
			},
			args: args{
				o: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.1"): {
							Ports: map[loadbalancer.L4Addr][]string{
								{
									Protocol: loadbalancer.NONE,
									Port:     1,
								}: {"foo"},
							},
							NodeName: "k8s1",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "different BE IPs",
			fields: fields{
				svcEP: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.1"): {
							Ports: map[loadbalancer.L4Addr][]string{
								{
									Protocol: loadbalancer.NONE,
									Port:     1,
								}: {"foo"},
							},
						},
					},
				},
			},
			args: args{
				o: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.2"): {
							Ports: map[loadbalancer.L4Addr][]string{
								{
									Protocol: loadbalancer.NONE,
									Port:     1,
								}: {"foo"},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "ports different name",
			fields: fields{
				svcEP: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.1"): {
							Ports: map[loadbalancer.L4Addr][]string{
								{
									Protocol: loadbalancer.NONE,
									Port:     1,
								}: {"foo"},
							},
						},
					},
				},
			},
			args: args{
				o: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.1"): {
							Ports: map[loadbalancer.L4Addr][]string{
								{
									Protocol: loadbalancer.NONE,
									Port:     1,
								}: {"foz"},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "ports different content",
			fields: fields{
				svcEP: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.1"): {
							Ports: map[loadbalancer.L4Addr][]string{
								{
									Protocol: loadbalancer.NONE,
									Port:     1,
								}: {"foo"},
							},
						},
					},
				},
			},
			args: args{
				o: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.1"): {
							Ports: map[loadbalancer.L4Addr][]string{
								{
									Protocol: loadbalancer.NONE,
									Port:     2,
								}: {"foo"},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "ports different one is bigger",
			fields: fields{
				svcEP: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.1"): {
							Ports: map[loadbalancer.L4Addr][]string{
								{
									Protocol: loadbalancer.NONE,
									Port:     1,
								}: {"foo"},
							},
						},
					},
				},
			},
			args: args{
				o: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.1"): {
							Ports: map[loadbalancer.L4Addr][]string{
								{
									Protocol: loadbalancer.NONE,
									Port:     1,
								}: {"foo"},
								{
									Protocol: loadbalancer.NONE,
									Port:     2,
								}: {"baz"},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name:   "backend different one is nil",
			fields: fields{},
			args: args{
				o: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.1"): {
							Ports: map[loadbalancer.L4Addr][]string{
								{
									Protocol: loadbalancer.NONE,
									Port:     1,
								}: {"foo"},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "node name different",
			fields: fields{
				svcEP: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.1"): {
							NodeName: "k8s2",
						},
					},
				},
			}, args: args{
				o: &Endpoints{
					Backends: map[cmtypes.AddrCluster]*Backend{
						cmtypes.MustParseAddrCluster("172.20.0.1"): {
							NodeName: "k8s1",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "both nil",
			args: args{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.svcEP.DeepEqual(tt.args.o); got != tt.want {
				t.Errorf("Endpoints.DeepEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseK8sEPSlicev1(t *testing.T) {
	nodeName := "k8s1"
	hostname := "pod-1"

	meta := slim_metav1.ObjectMeta{
		Name:      "foo",
		Namespace: "bar",
		Labels:    map[string]string{slim_discovery_v1.LabelServiceName: "quux"},
	}
	sliceID := EndpointSliceID{
		ServiceName:       loadbalancer.NewServiceName("bar", "quux"),
		EndpointSliceName: "foo",
	}

	newEmptyEndpoints := func() *Endpoints {
		eps := newEndpoints()
		eps.ObjectMeta = meta
		eps.EndpointSliceID = sliceID
		return eps
	}

	type args struct {
		eps            *slim_discovery_v1.EndpointSlice
		overrideConfig func()
	}
	tests := []struct {
		name        string
		setupArgs   func() args
		setupWanted func() *Endpoints
	}{
		{
			name: "empty endpoint",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv4,
						ObjectMeta:  meta,
					},
				}
			},
			setupWanted: newEmptyEndpoints,
		},
		{
			name: "endpoint with an address and port",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv4,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Addresses: []string{
									"172.0.0.1",
								},
								DeprecatedTopology: map[string]string{
									"kubernetes.io/hostname": nodeName,
								},
								Hostname: func() *string { return &hostname }(),
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				svcEP := newEmptyEndpoints()
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.1")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
					},
					NodeName: nodeName,
					Hostname: hostname,
				}
				return svcEP
			},
		},
		{
			name: "endpoint with an address and 2 ports",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv4,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Addresses: []string{
									"172.0.0.1",
								},
								DeprecatedTopology: map[string]string{
									"kubernetes.io/hostname": nodeName,
								},
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
							{
								Name:     func() *string { a := "http-test-svc-2"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8081); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				svcEP := newEmptyEndpoints()
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.1")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
					NodeName: nodeName,
				}
				return svcEP
			},
		},
		{
			name: "endpoint with 2 addresses and 2 ports",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv4,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Addresses: []string{
									"172.0.0.1",
								},
								DeprecatedTopology: map[string]string{
									"kubernetes.io/hostname": nodeName,
								},
							},
							{
								Addresses: []string{
									"172.0.0.2",
								},
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
							{
								Name:     func() *string { a := "http-test-svc-2"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8081); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				svcEP := newEmptyEndpoints()
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.1")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
					NodeName: nodeName,
				}
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.2")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
				}
				return svcEP
			},
		},
		{
			name: "endpoint with 2 addresses, 1 address not ready and 2 ports",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv4,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Addresses: []string{
									"172.0.0.1",
								},
								DeprecatedTopology: map[string]string{
									"kubernetes.io/hostname": nodeName,
								},
							},
							{
								Addresses: []string{
									"172.0.0.2",
								},
							},
							{
								Conditions: slim_discovery_v1.EndpointConditions{
									Ready: func() *bool { a := false; return &a }(),
								},
								Addresses: []string{
									"172.0.0.3",
								},
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
							{
								Name:     func() *string { a := "http-test-svc-2"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8081); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				svcEP := newEmptyEndpoints()
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.1")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
					NodeName: nodeName,
				}
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.2")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
				}
				return svcEP
			},
		}, {
			name: "endpoint with 2 addresses, 1 address not ready and 2 ports",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv4,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Addresses: []string{
									"172.0.0.1",
								},
								NodeName: func() *string { return &nodeName }(),
							},
							{
								Addresses: []string{
									"172.0.0.2",
								},
							},
							{
								Conditions: slim_discovery_v1.EndpointConditions{
									Ready: func() *bool { a := false; return &a }(),
								},
								Addresses: []string{
									"172.0.0.3",
								},
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
							{
								Name:     func() *string { a := "http-test-svc-2"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8081); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				svcEP := newEmptyEndpoints()
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.1")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
					NodeName: nodeName,
				}
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.2")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
				}
				return svcEP
			},
		},
		{
			name: "endpoints with some addresses not ready and not serving and terminating",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv4,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Addresses: []string{
									"172.0.0.1",
								},
							},
							{
								Conditions: slim_discovery_v1.EndpointConditions{
									Ready:       func() *bool { a := false; return &a }(),
									Serving:     func() *bool { a := false; return &a }(),
									Terminating: func() *bool { a := true; return &a }(),
								},
								Addresses: []string{
									"172.0.0.2",
								},
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
							{
								Name:     func() *string { a := "http-test-svc-2"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8081); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				svcEP := newEmptyEndpoints()
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.1")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
				}
				return svcEP
			},
		},
		{
			name: "endpoints with some addresses not ready and serving and terminating",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv4,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Addresses: []string{
									"172.0.0.1",
								},
							},
							{
								Conditions: slim_discovery_v1.EndpointConditions{
									Ready:       func() *bool { a := false; return &a }(),
									Serving:     func() *bool { a := true; return &a }(),
									Terminating: func() *bool { a := true; return &a }(),
								},
								Addresses: []string{
									"172.0.0.2",
								},
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
							{
								Name:     func() *string { a := "http-test-svc-2"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8081); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				svcEP := newEmptyEndpoints()
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.1")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
				}
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.2")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
					Terminating: true,
				}
				return svcEP
			},
		},

		{
			name: "endpoints with all addresses ready and serving and terminating",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv4,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Conditions: slim_discovery_v1.EndpointConditions{
									Ready:       func() *bool { a := true; return &a }(),
									Serving:     func() *bool { a := true; return &a }(),
									Terminating: func() *bool { a := true; return &a }(),
								},
								Addresses: []string{
									"172.0.0.1",
								},
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
							{
								Name:     func() *string { a := "http-test-svc-2"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8081); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				svcEP := newEmptyEndpoints()
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.1")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
				}
				return svcEP
			},
		},
		{
			name: "endpoints with all addresses not ready and not serving and terminating",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv4,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Conditions: slim_discovery_v1.EndpointConditions{
									Ready:       func() *bool { a := false; return &a }(),
									Terminating: func() *bool { a := true; return &a }(),
								},
								Addresses: []string{
									"172.0.0.1",
								},
							},
							{
								Conditions: slim_discovery_v1.EndpointConditions{
									Ready:       func() *bool { a := false; return &a }(),
									Terminating: func() *bool { a := true; return &a }(),
								},
								Addresses: []string{
									"172.0.0.2",
								},
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
							{
								Name:     func() *string { a := "http-test-svc-2"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8081); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				svcEP := newEmptyEndpoints()
				return svcEP
			},
		},
		{
			name: "endpoints with all addresses not ready and serving and terminating",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv4,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Conditions: slim_discovery_v1.EndpointConditions{
									Ready:       func() *bool { a := false; return &a }(),
									Serving:     func() *bool { a := true; return &a }(),
									Terminating: func() *bool { a := true; return &a }(),
								},
								Addresses: []string{
									"172.0.0.1",
								},
							},
							{
								Conditions: slim_discovery_v1.EndpointConditions{
									Ready:       func() *bool { a := false; return &a }(),
									Serving:     func() *bool { a := true; return &a }(),
									Terminating: func() *bool { a := true; return &a }(),
								},
								Addresses: []string{
									"172.0.0.2",
								},
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
							{
								Name:     func() *string { a := "http-test-svc-2"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8081); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				svcEP := newEmptyEndpoints()
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.1")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
					Terminating: true,
				}
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.2")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8081): {"http-test-svc-2"},
					},
					Terminating: true,
				}
				return svcEP
			},
		},
		{
			name: "endpoints have zone hints",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv4,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Addresses: []string{"172.0.0.1"},
								Hints: &slim_discovery_v1.EndpointHints{
									ForZones: []slim_discovery_v1.ForZone{{Name: "testing"}},
								},
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				svcEP := newEmptyEndpoints()
				svcEP.Backends[cmtypes.MustParseAddrCluster("172.0.0.1")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
					},
					HintsForZones: []string{"testing"},
				}
				return svcEP
			},
		},
		{
			name: "endpoint with IPv6 address type",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeIPv6,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Addresses: []string{
									"fd00::1",
								},
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				svcEP := newEmptyEndpoints()
				svcEP.Backends[cmtypes.MustParseAddrCluster("fd00::1")] = &Backend{
					Ports: map[loadbalancer.L4Addr][]string{
						loadbalancer.NewL4Addr(loadbalancer.TCP, 8080): {"http-test-svc"},
					},
				}
				return svcEP
			},
		},
		{
			name: "endpoint with FQDN address type",
			setupArgs: func() args {
				return args{
					eps: &slim_discovery_v1.EndpointSlice{
						AddressType: slim_discovery_v1.AddressTypeFQDN,
						ObjectMeta:  meta,
						Endpoints: []slim_discovery_v1.Endpoint{
							{
								Addresses: []string{
									"foo.example.com",
								},
							},
						},
						Ports: []slim_discovery_v1.EndpointPort{
							{
								Name:     func() *string { a := "http-test-svc"; return &a }(),
								Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
								Port:     func() *int32 { a := int32(8080); return &a }(),
							},
						},
					},
				}
			},
			setupWanted: func() *Endpoints {
				// We don't support FQDN address types. Should be empty.
				return newEmptyEndpoints()
			},
		},
	}
	for _, tt := range tests {
		args := tt.setupArgs()
		want := tt.setupWanted()
		if args.overrideConfig != nil {
			args.overrideConfig()
		}
		got := ParseEndpointSliceV1(hivetest.Logger(t), args.eps)
		require.Equal(t, want, got, "Test name: %q", tt.name)
	}
}

func Test_parseEndpointPortV1(t *testing.T) {
	type args struct {
		port slim_discovery_v1.EndpointPort
	}
	tests := []struct {
		name     string
		args     args
		portName string
		l4Addr   loadbalancer.L4Addr
		fail     bool
	}{
		{
			name: "tcp-port",
			args: args{
				port: slim_discovery_v1.EndpointPort{
					Name:     func() *string { a := "http-test-svc"; return &a }(),
					Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolTCP; return &a }(),
					Port:     func() *int32 { a := int32(8080); return &a }(),
				},
			},
			portName: "http-test-svc",
			l4Addr: loadbalancer.L4Addr{
				Protocol: loadbalancer.TCP,
				Port:     8080,
			},
		},
		{
			name: "udp-port",
			args: args{
				port: slim_discovery_v1.EndpointPort{
					Name:     func() *string { a := "http-test-svc"; return &a }(),
					Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolUDP; return &a }(),
					Port:     func() *int32 { a := int32(8080); return &a }(),
				},
			},
			portName: "http-test-svc",
			l4Addr: loadbalancer.L4Addr{
				Protocol: loadbalancer.UDP,
				Port:     8080,
			},
		},
		{
			name: "sctp-port",
			args: args{
				port: slim_discovery_v1.EndpointPort{
					Name:     func() *string { a := "sctp-test-svc"; return &a }(),
					Protocol: func() *slim_corev1.Protocol { a := slim_corev1.ProtocolSCTP; return &a }(),
					Port:     func() *int32 { a := int32(5555); return &a }(),
				},
			},
			portName: "sctp-test-svc",
			l4Addr: loadbalancer.L4Addr{
				Protocol: loadbalancer.SCTP,
				Port:     5555,
			},
		},
		{
			name: "unset-protocol-should-have-tcp-port",
			args: args{
				port: slim_discovery_v1.EndpointPort{
					Name: func() *string { a := "http-test-svc"; return &a }(),
					Port: func() *int32 { a := int32(8080); return &a }(),
				},
			},
			portName: "http-test-svc",
			l4Addr: loadbalancer.L4Addr{
				Protocol: loadbalancer.TCP,
				Port:     8080,
			},
		},
		{
			name: "unset-port-number-should-fail",
			args: args{
				port: slim_discovery_v1.EndpointPort{
					Name: func() *string { a := "http-test-svc"; return &a }(),
				},
			},
			fail: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPortName, gotL4Addr, ok := parseEndpointPortV1(tt.args.port)
			if !ok && !tt.fail {
				t.Errorf("parseEndpointPortV1() failed with %v", tt.args.port)
			}
			if gotPortName != tt.portName {
				t.Errorf("parseEndpointPortV1() got = %v, want %v", gotPortName, tt.portName)
			}
			if !gotL4Addr.Equals(tt.l4Addr) {
				t.Errorf("parseEndpointPortV1() got1 = %v, want %v", gotL4Addr, tt.l4Addr)
			}
		})
	}
}
