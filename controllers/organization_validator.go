package controllers

import (
	"context"
	"fmt"
	"net/http"

	securityv1alpha1 "github.com/giantswarm/apiextensions/v3/pkg/apis/security/v1alpha1"
	"github.com/giantswarm/apiextensions/v3/pkg/label"
	"github.com/giantswarm/microerror"
	"go.uber.org/zap"
	admissionv1 "k8s.io/api/admission/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	capi "sigs.k8s.io/cluster-api/api/v1alpha3"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/validate-security-giantswarm-io-v1alpha1-organization,mutating=false,failurePolicy=fail,sideEffects=None,groups=security.giantswarm.io,resources=organizations,verbs=delete,versions=v1alpha1,name=organization.security.giantswarm.io,admissionReviewVersions={v1,v1beta1,v1alpha1}

// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters,verbs=list;watch

type OrganizationValidator struct {
	client.Client
	Log *zap.SugaredLogger

	decoder *admission.Decoder
}

// So the decoder is injected.
var _ admission.DecoderInjector = &OrganizationValidator{}

func (v *OrganizationValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	v.Log = v.Log.Desugar().With(
		zap.String("kind", req.Kind.String()),
		zap.String("object", client.ObjectKey{Namespace: req.Namespace, Name: req.Name}.String()),
		zap.String("op", string(req.Operation)),
		zap.String("UID", string(req.UID)),
	).Sugar()

	organization := &securityv1alpha1.Organization{}

	err := v.decoder.DecodeRaw(req.OldObject, organization)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	resp, err := v.handle(ctx, req, organization)
	if err != nil {
		v.Log.With("stack", microerror.JSON(err)).Error("failed with error")
		return admission.Errored(500, err)
	}

	return resp
}

func (v *OrganizationValidator) handle(ctx context.Context, req admission.Request, organization *securityv1alpha1.Organization) (admission.Response, error) {
	var err error

	orgName := organization.Name
	orgLabel := label.Organization

	var orgClustersSelector labels.Selector
	{
		s := orgLabel + "=" + orgName
		orgClustersSelector, err = labels.Parse(s)
		if err != nil {
			return admission.Response{}, microerror.Mask(fmt.Errorf("failed to create selector for string = %#q with error: %s", s, err))
		}
	}

	orgClusters := capi.ClusterList{}
	{
		err := v.List(ctx, &orgClusters, &client.ListOptions{
			LabelSelector: orgClustersSelector,
		})
		if apimeta.IsNoMatchError(err) {
			// If the CRD is not in place, just ignore the error.
		}
		if err != nil {
			return admission.Response{}, microerror.Mask(err)
		}
	}

	if req.Operation == admissionv1.Delete && len(orgClusters.Items) > 0 {
		return admission.Denied(fmt.Sprintf(
			"Found %d clusters for selector = %q and operation = %q, expected at most 0",
			len(orgClusters.Items), orgClustersSelector, req.Operation,
		)), nil
	}

	return admission.Allowed(fmt.Sprintf(
		"No clusters found for selector = %q and operation = %q",
		orgClustersSelector, req.Operation,
	)), nil
}

func (v *OrganizationValidator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	hookServer := mgr.GetWebhookServer()
	// Note it has to be the same path as in kubebuilder:webhook:path
	// marker.
	hookServer.Register("/validate-security-giantswarm-io-v1alpha1-organization", &webhook.Admission{Handler: v})
	return nil
}

func (v *OrganizationValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
