package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/giantswarm/microerror"
	"go.uber.org/zap"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"

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
	var err error

	const (
		appNameLabel    = "app.kubernetes.io/name"
		appVersionLabel = "app.kubernetes.io/version"
	)

	if len(deployment.Labels) == 0 {
		return admission.Allowed("No labels"), nil
	}
	appName, ok := deployment.Labels[appNameLabel]
	if !ok {
		return admission.Allowed(fmt.Sprintf("Label %#q not found", appNameLabel)), nil
	}
	appVersion, ok := deployment.Labels[appVersionLabel]
	if !ok {
		return admission.Allowed(fmt.Sprintf("Label %#q not found", appVersionLabel)), nil
	}

	var sameNameSelector labels.Selector
	{
		s := appNameLabel + "=" + appName + "," + appVersionLabel + "=" + appVersion
		sameNameSelector, err = labels.Parse(s)
		if err != nil {
			return admission.Response{}, microerror.Mask(fmt.Errorf("failed to create selector for string = %#q with error: %s", s, err))
		}
	}

	sameNameDeployments := appsv1.DeploymentList{}
	{
		err := v.List(ctx, &sameNameDeployments, &client.ListOptions{
			LabelSelector: sameNameSelector,
		})
		if err != nil {
			return admission.Response{}, microerror.Mask(err)
		}
	}

	switch {
	case req.Operation == admissionv1.Create && len(sameNameDeployments.Items) > 0:
		return admission.Denied(fmt.Sprintf(
			"Found %d deployments for selector = %q and operation = %q, expected at most 0",
			len(sameNameDeployments.Items), sameNameSelector, req.Operation,
		)), nil
	case req.Operation == admissionv1.Create:
		return admission.Allowed(fmt.Sprintf(
			"Found %d deployments for selector = %q and operation = %q",
			len(sameNameDeployments.Items), sameNameSelector, req.Operation,
		)), nil
	case req.Operation == admissionv1.Update && len(sameNameDeployments.Items) > 1:
		return admission.Denied(fmt.Sprintf(
			"Found %d deployments for selector = %q and operation = %q, expected at most 1",
			len(sameNameDeployments.Items), sameNameSelector, req.Operation,
		)), nil
	case req.Operation == admissionv1.Update:
		return admission.Allowed(fmt.Sprintf(
			"Found %d deployments for selector = %q and operation = %q",
			len(sameNameDeployments.Items), sameNameSelector, req.Operation,
		)), nil
	}

	return admission.Response{}, microerror.Mask(fmt.Errorf("unsupported operation %#q", req.Operation))
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
