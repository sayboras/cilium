// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package gateway_api

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/cilium/cilium/pkg/logging/logfields"
)

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *gatewayClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	scopedLog := log.WithContext(ctx).WithFields(logrus.Fields{
		logfields.Controller: "gatewayclass",
		logfields.Resource:   req.NamespacedName,
	})

	scopedLog.Info("Reconciling GatewayClass")
	gwc := &gatewayv1beta1.GatewayClass{}
	if err := r.Client.Get(ctx, req.NamespacedName, gwc); err != nil {
		if k8serrors.IsNotFound(err) {
			return success()
		}
		return fail(err)
	}

	// Ignore deleted GatewayClass, this can happen when foregroundDeletion is enabled
	// The reconciliation loop will automatically kick off for related Gateway resources.
	if gwc.GetDeletionTimestamp() != nil {
		return success()
	}

	for _, fn := range []gatewayClassChecker{
		validateParamRef,
	} {
		if res, err := fn(ctx, r.Client, gwc); err != nil {
			return res, err
		}
	}
	setGatewayClassAccepted(gwc, true, gatewayClassAcceptedMessage)

	if err := r.Client.Status().Update(ctx, gwc); err != nil {
		scopedLog.WithError(err).Error("Failed to update GatewayClass status")
		return fail(err)
	}
	scopedLog.Info("Successfully reconciled GatewayClass")
	return success()
}

func validateParamRef(ctx context.Context, c client.Client, gwc *gatewayv1beta1.GatewayClass) (ctrl.Result, error) {
	if gwc.Spec.ParametersRef == nil {
		return success()
	}

	param := gwc.Spec.ParametersRef
	if param.Group != "v1" || param.Kind != "ConfigMap" {
		errorMsg := "only ConfigMap is supported for parametersRef"
		setGatewayClassAccepted(gwc, false, errorMsg)
		return fail(fmt.Errorf(errorMsg))
	}

	cm, err := getParameterRef(ctx, c, gwc)
	if err != nil {
		setGatewayClassAccepted(gwc, false, err.Error())
		return fail(err)
	}

	if cm == nil {
		errorMsg := "parameterRef ConfigMap not found"
		setGatewayClassAccepted(gwc, false, errorMsg)
		return fail(fmt.Errorf(errorMsg))
	}
	return success()
}

func getParameterRef(ctx context.Context, c client.Client, gwc *gatewayv1beta1.GatewayClass) (*corev1.ConfigMap, error) {
	if gwc.Spec.ParametersRef == nil {
		return nil, nil
	}

	param := gwc.Spec.ParametersRef
	if param.Group != "v1" || param.Kind != "ConfigMap" {
		return nil, fmt.Errorf("only ConfigMap is supported for parametersRef")
	}

	cm := &corev1.ConfigMap{}
	if err := c.Get(ctx, client.ObjectKey{
		Namespace: namespaceDerefOr(param.Namespace, gwc.GetNamespace()),
		Name:      param.Name,
	}, cm); err != nil {
		if !k8serrors.IsNotFound(err) {
			return nil, err
		}
	}
	return cm, nil
}
