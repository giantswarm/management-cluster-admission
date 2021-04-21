package controllers

import (
	"context"
	"net/http"

	"github.com/giantswarm/microerror"
	appsv1 "k8s.io/api/apps/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch
type HealthProbe struct {
	client.Client
}

func (v *HealthProbe) HealthzCheck(_ *http.Request) error {
	// Check whether manager's client is able to communicate with Kubernetes API.
	kubeSystemDeployments := appsv1.DeploymentList{}
	err := v.List(context.Background{}, &kubeSystemDeployments, &client.ListOptions{
		Namespace: "kube-system",
		Limit:     1,
	})
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}
