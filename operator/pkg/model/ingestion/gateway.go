// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package ingestion

import (
	"strings"

	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	corev1 "k8s.io/api/core/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/cilium/cilium/operator/pkg/model"
)

const (
	allHosts = "*"
)

// AdditionalParameter is parametersRef configmap in GatewayClass.
type AdditionalParameter struct {
	ServiceType         string `mapstructure:"service-type" default:"LoadBalancer"`
	ServiceSecurePort   int    `mapstructure:"secure-port" default:"443"`
	ServiceInsecurePort int    `mapstructure:"insecure-port" default:"80"`
}

// Input is the input for GatewayAPI.
type Input struct {
	Gateway          gatewayv1beta1.Gateway
	HTTPRoutes       []gatewayv1beta1.HTTPRoute
	ReferenceGrants  []gatewayv1alpha2.ReferenceGrant
	Services         []corev1.Service
	AdditionalParams map[string]string
}

// GatewayAPI translates Gateway API resources into a model.
// The current implementation only supports HTTPRoute.
func GatewayAPI(input Input) []model.HTTPListener {
	var res []model.HTTPListener

	for _, l := range input.Gateway.Spec.Listeners {
		if l.Protocol != gatewayv1beta1.HTTPProtocolType &&
			l.Protocol != gatewayv1beta1.HTTPSProtocolType &&
			l.Protocol != gatewayv1beta1.TLSProtocolType {
			continue
		}

		var routes []model.HTTPRoute

		for _, r := range filterRoute(input.Gateway, l, input.HTTPRoutes) {
			computedHost := model.ComputeHosts(toStringSlice(r.Spec.Hostnames), (*string)(l.Hostname))
			// No matching host, skip this route
			if len(computedHost) == 0 {
				continue
			}

			if len(computedHost) == 1 && computedHost[0] == allHosts {
				computedHost = nil
			}

			for _, rule := range r.Spec.Rules {
				bes := make([]model.Backend, 0, len(rule.BackendRefs))
				for _, be := range rule.BackendRefs {
					if !isReferenceAllowed(r.GetNamespace(), be, input.ReferenceGrants) {
						continue
					}
					if (be.Kind != nil && *be.Kind != "Service") || (be.Group != nil && *be.Group != corev1.GroupName) {
						continue
					}
					if serviceExists(string(be.Name), namespaceDerefOr(be.Namespace, r.Namespace), input.Services) {
						bes = append(bes, toBackend(be, r.Namespace))
					}
				}

				var dr *model.DirectResponse
				if len(bes) == 0 {
					dr = &model.DirectResponse{
						StatusCode: 500,
					}
				}

				var requestHeaderFilter *model.HTTPHeaderFilter
				var responseHeaderFilter *model.HTTPHeaderFilter
				var requestRedirectFilter *model.HTTPRequestRedirectFilter
				if len(rule.Filters) > 0 {
					for _, f := range rule.Filters {
						if f.Type == gatewayv1beta1.HTTPRouteFilterRequestHeaderModifier {
							requestHeaderFilter = &model.HTTPHeaderFilter{
								HeadersToAdd:    toHTTPHeaders(f.RequestHeaderModifier.Add),
								HeadersToSet:    toHTTPHeaders(f.RequestHeaderModifier.Set),
								HeadersToRemove: f.RequestHeaderModifier.Remove,
							}
						}

						if f.Type == gatewayv1beta1.HTTPRouteFilterResponseHeaderModifier {
							responseHeaderFilter = &model.HTTPHeaderFilter{
								HeadersToAdd:    toHTTPHeaders(f.ResponseHeaderModifier.Add),
								HeadersToSet:    toHTTPHeaders(f.ResponseHeaderModifier.Set),
								HeadersToRemove: f.ResponseHeaderModifier.Remove,
							}
						}

						if f.Type == gatewayv1beta1.HTTPRouteFilterRequestRedirect {
							requestRedirectFilter = toHTTPRequestRedirectFilter(f.RequestRedirect)
						}
					}
				}

				if len(rule.Matches) == 0 {
					routes = append(routes, model.HTTPRoute{
						Hostnames:              computedHost,
						Backends:               bes,
						DirectResponse:         dr,
						RequestHeaderFilter:    requestHeaderFilter,
						ResponseHeaderModifier: responseHeaderFilter,
						RequestRedirect:        requestRedirectFilter,
					})
				}

				for _, match := range rule.Matches {
					routes = append(routes, model.HTTPRoute{
						Hostnames:              computedHost,
						PathMatch:              toPathMatch(match),
						HeadersMatch:           toHeaderMatch(match),
						QueryParamsMatch:       toQueryMatch(match),
						Method:                 (*string)(match.Method),
						Backends:               bes,
						DirectResponse:         dr,
						RequestHeaderFilter:    requestHeaderFilter,
						ResponseHeaderModifier: responseHeaderFilter,
						RequestRedirect:        requestRedirectFilter,
					})
				}
			}
		}

		params := toParameter(input.AdditionalParams)
		res = append(res, model.HTTPListener{
			Name: string(l.Name),
			Sources: []model.FullyQualifiedResource{
				{
					Name:      input.Gateway.GetName(),
					Namespace: input.Gateway.GetNamespace(),
					Group:     input.Gateway.GroupVersionKind().Group,
					Version:   input.Gateway.GroupVersionKind().Version,
					Kind:      input.Gateway.GroupVersionKind().Kind,
					UID:       string(input.Gateway.GetUID()),
				},
			},
			Port:     uint32(l.Port),
			Hostname: toHostname(l.Hostname),
			TLS:      toTLS(l.TLS, input.Gateway.GetNamespace()),
			Routes:   routes,
			Service: &model.Service{
				Type:             params.ServiceType,
				InsecureNodePort: model.AddressOf[uint32](uint32(params.ServiceInsecurePort)),
				SecureNodePort:   model.AddressOf[uint32](uint32(params.ServiceSecurePort)),
			},
		})
	}

	return res
}

func toHTTPRequestRedirectFilter(redirect *gatewayv1beta1.HTTPRequestRedirectFilter) *model.HTTPRequestRedirectFilter {
	if redirect == nil {
		return nil
	}
	var pathModifier *model.StringMatch
	if redirect.Path != nil {
		pathModifier = &model.StringMatch{}

		switch redirect.Path.Type {
		case gatewayv1beta1.FullPathHTTPPathModifier:
			pathModifier.Exact = *redirect.Path.ReplaceFullPath
		case gatewayv1beta1.PrefixMatchHTTPPathModifier:
			pathModifier.Prefix = *redirect.Path.ReplacePrefixMatch
		}
	}
	return &model.HTTPRequestRedirectFilter{
		Scheme:     redirect.Scheme,
		Hostname:   (*string)(redirect.Hostname),
		Path:       pathModifier,
		Port:       (*int32)(redirect.Port),
		StatusCode: redirect.StatusCode,
	}
}

func filterRoute(gw gatewayv1beta1.Gateway, listener gatewayv1beta1.Listener, routes []gatewayv1beta1.HTTPRoute) []gatewayv1beta1.HTTPRoute {
	var res []gatewayv1beta1.HTTPRoute

	for _, r := range routes {
		for _, parent := range r.Spec.ParentRefs {
			if gw.GetName() != string(parent.Name) ||
				gw.GetNamespace() != namespaceDerefOr(parent.Namespace, r.Namespace) {
				continue
			}
			if parent.SectionName != nil && *parent.SectionName != listener.Name {
				continue
			}
			res = append(res, r)
		}
	}

	return res
}

// isReferenceAllowed returns true if the reference is allowed by the reference grant.
// TODO(tam): only HTTPRoute with Service is supported right now.
// We need to support other routes (e.g. grpc, tls, etc.) later.
func isReferenceAllowed(originatingNamespace string, be gatewayv1beta1.HTTPBackendRef, grants []gatewayv1alpha2.ReferenceGrant) bool {
	if be.Namespace == nil || string(*be.Namespace) == originatingNamespace {
		return true
	}
	for _, g := range grants {
		for _, from := range g.Spec.From {
			if from.Group == gatewayv1beta1.GroupName &&
				from.Kind == "HTTPRoute" && (string)(from.Namespace) == originatingNamespace {
				for _, to := range g.Spec.To {
					if to.Group == corev1.GroupName && to.Kind == "Service" &&
						(to.Name == nil || string(*to.Name) == string(be.Name)) {
						return true
					}
				}
			}
		}
	}
	return false
}

func toHostname(hostname *gatewayv1beta1.Hostname) string {
	if hostname != nil {
		return (string)(*hostname)
	}
	return allHosts
}

func serviceExists(svcName, svcNamespace string, services []corev1.Service) bool {
	for _, svc := range services {
		if svc.GetName() == svcName && svc.GetNamespace() == svcNamespace {
			return true
		}
	}
	return false
}

func toBackend(be gatewayv1beta1.HTTPBackendRef, defaultNamespace string) model.Backend {
	ns := namespaceDerefOr(be.Namespace, defaultNamespace)
	var port *model.BackendPort

	if be.Port != nil {
		port = &model.BackendPort{
			Port: uint32(*be.Port),
		}
	}
	return model.Backend{
		Name:      string(be.Name),
		Namespace: ns,
		Port:      port,
		Weight:    be.Weight,
	}
}

func namespaceDerefOr(namespace *gatewayv1beta1.Namespace, defaultNamespace string) string {
	if namespace != nil && *namespace != "" {
		return string(*namespace)
	}
	return defaultNamespace
}

func toPathMatch(match gatewayv1beta1.HTTPRouteMatch) model.StringMatch {
	if match.Path == nil {
		return model.StringMatch{}
	}

	switch *match.Path.Type {
	case gatewayv1beta1.PathMatchExact:
		return model.StringMatch{
			Exact: *match.Path.Value,
		}
	case gatewayv1beta1.PathMatchPathPrefix:
		return model.StringMatch{
			Prefix: *match.Path.Value,
		}
	case gatewayv1beta1.PathMatchRegularExpression:
		return model.StringMatch{
			Regex: *match.Path.Value,
		}
	}
	return model.StringMatch{}
}

func toHeaderMatch(match gatewayv1beta1.HTTPRouteMatch) []model.KeyValueMatch {
	if len(match.Headers) == 0 {
		return nil
	}
	res := make([]model.KeyValueMatch, 0, len(match.Headers))
	for _, h := range match.Headers {
		t := gatewayv1beta1.HeaderMatchExact
		if h.Type != nil {
			t = *h.Type
		}
		switch t {
		case gatewayv1beta1.HeaderMatchExact:
			res = append(res, model.KeyValueMatch{
				Key: string(h.Name),
				Match: model.StringMatch{
					Exact: h.Value,
				},
			})
		case gatewayv1beta1.HeaderMatchRegularExpression:
			res = append(res, model.KeyValueMatch{
				Key: string(h.Name),
				Match: model.StringMatch{
					Regex: h.Value,
				},
			})
		}
	}
	return res
}

func toQueryMatch(match gatewayv1beta1.HTTPRouteMatch) []model.KeyValueMatch {
	if len(match.QueryParams) == 0 {
		return nil
	}
	res := make([]model.KeyValueMatch, 0, len(match.QueryParams))
	for _, h := range match.QueryParams {
		t := gatewayv1beta1.QueryParamMatchExact
		if h.Type != nil {
			t = *h.Type
		}
		switch t {
		case gatewayv1beta1.QueryParamMatchExact:
			res = append(res, model.KeyValueMatch{
				Key: h.Name,
				Match: model.StringMatch{
					Exact: h.Value,
				},
			})
		case gatewayv1beta1.QueryParamMatchRegularExpression:
			res = append(res, model.KeyValueMatch{
				Key: h.Name,
				Match: model.StringMatch{
					Regex: h.Value,
				},
			})
		}
	}
	return res
}

func toTLS(tls *gatewayv1beta1.GatewayTLSConfig, defaultNamespace string) []model.TLSSecret {
	if tls == nil {
		return nil
	}

	res := make([]model.TLSSecret, 0, len(tls.CertificateRefs))
	for _, cert := range tls.CertificateRefs {
		res = append(res, model.TLSSecret{
			Name:      string(cert.Name),
			Namespace: namespaceDerefOr(cert.Namespace, defaultNamespace),
		})
	}
	return res
}

func toHTTPHeaders(headers []gatewayv1beta1.HTTPHeader) []model.Header {
	if len(headers) == 0 {
		return nil
	}
	res := make([]model.Header, 0, len(headers))
	for _, h := range headers {
		res = append(res, model.Header{
			Name:  string(h.Name),
			Value: h.Value,
		})
	}
	return res
}

func toStringSlice(s []gatewayv1beta1.Hostname) []string {
	res := make([]string, 0, len(s))
	for _, h := range s {
		res = append(res, string(h))
	}
	return res
}

func toParameter(params map[string]string) AdditionalParameter {
	p := AdditionalParameter{}
	defaults.SetDefaults(&p)
	decoder, err := mapstructure.NewDecoder(decoderConfig(&p))
	if err != nil {
		return p
	}
	_ = decoder.Decode(params)
	return p
}

func decoderConfig(target any) *mapstructure.DecoderConfig {
	return &mapstructure.DecoderConfig{
		Result:           target,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
		ZeroFields:  true,
		ErrorUnset:  false,
		ErrorUnused: false,
		// Match field FooBarBaz with "foo-bar-baz" by removing
		// the dashes from the flag.
		MatchName: func(mapKey, fieldName string) bool {
			return strings.EqualFold(
				strings.ReplaceAll(mapKey, "-", ""),
				fieldName)
		},
	}
}
