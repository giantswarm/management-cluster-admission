module github.com/giantswarm/management-cluster-admission

go 1.16

require (
	github.com/giantswarm/microerror v0.3.0
	github.com/go-logr/zapr v0.2.0
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/zap v1.16.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.7.2
)
