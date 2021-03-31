package controllers

import (
	"context"
	"net/http"

	appsv1 "k8s.io/api/apps/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//+kubebuilder:webhook:path=/validate-apps-v1-deployment,mutating=false,failurePolicy=fail,sideEffects=None,groups=apps,resources=deployments,verbs=create;update,versions=v1,name=vdeployment.kb.io,admissionReviewVersions={v1,v1beta1}

// +kubebuilder:rbac:groups=apps/v1,resources=deployments,verbs=get;list;watch;create;update;patch;delete

type DeploymentValidator struct {
	client.Client

	decoder *admission.Decoder
}

// So the decoder is injected.
var _ admission.DecoderInjector = &DeploymentValidator{}

func (v *DeploymentValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	deployment := &appsv1.Deployment{}

	err := v.decoder.Decode(req, deployment)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	return admission.Denied("ALWAYS FAIL")
	//return admission.Allowed("")
}

func (v *DeploymentValidator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	hookServer := mgr.GetWebhookServer()
	// Note it has to be the same path as in kubebuilder:webhook:path
	// marker.
	hookServer.Register("/validate-apps-v1-deployment", &webhook.Admission{Handler: v})
	return nil
}

func (v *DeploymentValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
