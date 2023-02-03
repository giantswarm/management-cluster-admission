package project

var (
	description = "Webhooks for Management Clusters."
	gitSHA      = "n/a"
	name        = "management-cluster-admission"
	source      = "https://github.com/giantswarm/management-cluster-admission"
	version     = "0.8.1"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
