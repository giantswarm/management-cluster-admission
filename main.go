/*
Copyright 2021 Giant Swarm.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"github.com/giantswarm/microerror"
	"github.com/go-logr/zapr"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	capiv1alpha2 "sigs.k8s.io/cluster-api/api/v1alpha2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	//+kubebuilder:scaffold:imports

	"github.com/giantswarm/management-cluster-admission/controllers"
	"github.com/giantswarm/management-cluster-admission/pkg/project"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(capiv1alpha2.AddToScheme(scheme))

	//+kubebuilder:scaffold:scheme
}

var flags = struct {
	MetricsAddr          string
	EnableLeaderElection bool
	ProbeAddr            string
}{}

func initFlags() (errors []error) {
	// Flag/configuration names.
	const (
		flagLeaderElect            = "leader-elect"
		flagMetricsBindAddres      = "metrics-bind-address"
		flagHealthProbeBindAddress = "health-probe-bind-address"
	)

	// Flag binding.
	flag.StringVar(&flags.MetricsAddr, flagMetricsBindAddres, ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&flags.ProbeAddr, flagHealthProbeBindAddress, ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&flags.EnableLeaderElection, flagLeaderElect, false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	// Parse flags and configuration.
	flag.Parse()
	errors = append(errors, initFlagsFromEnv()...)

	// Validation.

	//if flags.Name == "" {
	//	errors = append(errors, fmt.Errorf("--%s flag must be set", flagName))
	//}

	return
}

func initFlagsFromEnv() (errors []error) {
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		if f.Changed {
			return
		}
		env := project.Name() + "_" + f.Name
		env = strings.ReplaceAll(env, ".", "_")
		env = strings.ReplaceAll(env, "-", "_")
		env = strings.ToUpper(env)
		v, ok := os.LookupEnv(env)
		if !ok {
			return
		}
		fmt.Printf("Setting --%s flag to value of $%s\n", f.Name, env)
		err := f.Value.Set(v)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to set --%s value using %q environment variable", f.Name, env))
		}
	})

	return
}

func main() {
	errs := initFlags()
	if len(errs) > 0 {
		ss := make([]string, len(errs))
		for i := range errs {
			ss[i] = errs[i].Error()
		}
		fmt.Fprintf(os.Stderr, "Error: %s\n", strings.Join(ss, "\nError: "))
		os.Exit(2)
	}

	err := mainE(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", microerror.Pretty(err, true))
		os.Exit(1)
	}
}

func mainE(ctx context.Context) error {
	rootLog, err := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  zapcore.OmitKey,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}.Build()
	if err != nil {
		return microerror.Mask(err)
	}

	ctrl.SetLogger(zapr.NewLogger(rootLog))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     flags.MetricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: flags.ProbeAddr,
		LeaderElection:         flags.EnableLeaderElection,
		LeaderElectionID:       "02c7a966.giantswarm.io",
	})
	if err != nil {
		return microerror.Mask(err)
	}

	controllersLog := rootLog.Named("controllers")
	if err = (&controllers.DeploymentValidator{
		Client: mgr.GetClient(),
		Log:    controllersLog.Named("deployment-validator").Sugar(),
	}).SetupWebhookWithManager(mgr); err != nil {
		return microerror.Mask(err)
	}
	if err = (&controllers.OrganizationValidator{
		Client: mgr.GetClient(),
		Log:    controllersLog.Named("organization-validator").Sugar(),
	}).SetupWebhookWithManager(mgr); err != nil {
		return microerror.Mask(err)
	}

	//+kubebuilder:scaffold:builder

	healthProbe := &controllers.HealthProbe{
		Client: mgr.GetClient(),
	}

	if err := mgr.AddHealthzCheck("healthz", healthProbe.HealthzCheck); err != nil {
		return microerror.Mask(err)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return microerror.Mask(err)
	}

	rootLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		return microerror.Mask(err)
	}

	return nil
}
