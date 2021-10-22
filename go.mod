module github.com/giantswarm/management-cluster-admission

go 1.16

require (
	github.com/giantswarm/apiextensions/v3 v3.35.0
	github.com/giantswarm/microerror v0.3.0
	github.com/go-logr/zapr v0.4.0
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.19.1
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	sigs.k8s.io/cluster-api v1.0.0
	sigs.k8s.io/controller-runtime v0.10.2
)

replace (
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
)
