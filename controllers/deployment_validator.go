package controllers

import (
	"context"
	"net/http"

	"github.com/giantswarm/microerror"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//+kubebuilder:webhook:path=/validate-apps-v1-deployment,mutating=false,failurePolicy=fail,sideEffects=None,groups=apps,resources=deployments,verbs=create;update,versions=v1,name=vdeployment.kb.io,admissionReviewVersions={v1,v1beta1}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

type DeploymentValidator struct {
	client.Client
	Log *zap.SugaredLogger

	decoder *admission.Decoder
}

// So the decoder is injected.
var _ admission.DecoderInjector = &DeploymentValidator{}

func (v *DeploymentValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	v.Log = v.Log.Desugar().With(
		zap.String("kind", req.Kind.String()),
		zap.String("object", client.ObjectKey{Namespace: req.Namespace, Name: req.Name}.String()),
		zap.String("op", string(req.Operation)),
		zap.String("UID", string(req.UID)),
	).Sugar()

	deployment := &appsv1.Deployment{}

	err := v.decoder.Decode(req, deployment)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	resp, err := v.handle(ctx, req, deployment)
	if err != nil {
		v.Log.With("stack", microerror.JSON(err)).Error("failed with error")
		return admission.Errored(500, err)
	}

	return resp
}

func (v *DeploymentValidator) handle(ctx context.Context, req admission.Request, deployment *appsv1.Deployment) (admission.Response, error) {
	return admission.Denied("__NOT_IMPLEMENTED__"), nil
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
