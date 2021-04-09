module github.com/giantswarm/management-cluster-admission

go 1.16

require (
	github.com/giantswarm/microerror v0.3.0
	github.com/go-logr/zapr v0.4.0
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.16.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.7.2
)

replace (
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
)
