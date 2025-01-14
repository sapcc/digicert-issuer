/*
Copyright 2022 SAP SE.

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
	"flag"
	"fmt"
	"os"
	"time"

	// Load k8s auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	certmanagerv1beta1 "github.com/sapcc/digicert-issuer/apis/certmanager/v1beta1"
	certmanagerv1beta1controller "github.com/sapcc/digicert-issuer/controllers/certmanager"
	"github.com/sapcc/digicert-issuer/pkg/version"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	// +kubebuilder:scaffold:scheme
	for _, addToSchemeFunc := range []func(s *runtime.Scheme) error{
		clientgoscheme.AddToScheme,
		certmanagerv1.AddToScheme,
		certmanagerv1beta1.AddToScheme,
	} {
		utilruntime.Must(addToSchemeFunc(scheme))
	}
}

func main() {
	var (
		metricsAddr                        string
		printVersionAndExit                bool
		backoffDurationProvisionerNotReady time.Duration
		backoffDurationRequestPending      time.Duration
		disableRootCA                      bool
	)

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080",
		"The address the metric endpoint binds to.")

	flag.BoolVar(&printVersionAndExit, "version", false,
		"Print version and exit.")

	flag.DurationVar(&backoffDurationProvisionerNotReady, "backoff-duration-provisioner-not-ready", 10*time.Second,
		"The backoff duration if the provisioner is not ready.")

	flag.DurationVar(&backoffDurationRequestPending, "backoff-duration-request-pending", 15*time.Minute,
		"The backoff duration if certificate request is pending.")

	flag.BoolVar(&disableRootCA, "disable-root-ca", false,
		"Enabling this removes root CA from CertificateRequest")

	flag.Parse()

	if printVersionAndExit {
		fmt.Println(version.Print("digicert-issuer"))
		os.Exit(0)
	}

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Cache:  cache.Options{DefaultTransform: cache.TransformStripManagedFields()},
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: metricsAddr,
		},
		LeaderElection:   true,
		LeaderElectionID: "digicertissuer.cloud.sap",
	})
	handleError(err, "unable to start manager")

	err = (&certmanagerv1beta1controller.DigicertIssuerReconciler{}).SetupWithManager(mgr)
	handleError(err, "unable to initialize controller", "controller", "digicertIssuer")

	err = (&certmanagerv1beta1controller.CertificateRequestReconciler{
		BackoffDurationProvisionerNotReady: backoffDurationProvisionerNotReady,
		BackoffDurationRequestPending:      backoffDurationRequestPending,
		DefaultProviderNamespace:           getValueFromEnvironmentOrDefault("POD_NAMESPACE", "kube-system"),
		DisableRootCA:                      disableRootCA,
	}).SetupWithManager(mgr)
	handleError(err, "unable to initialize controller", "controller", "certificateRequest")

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	err = mgr.Start(ctrl.SetupSignalHandler())
	handleError(err, "problem running manager")
}

func handleError(err error, message string, keysAndVals ...interface{}) {
	if err != nil {
		setupLog.Error(err, message, keysAndVals...)
		os.Exit(1)
	}
}

func getValueFromEnvironmentOrDefault(envKey, defaultValue string) string {
	if val, ok := os.LookupEnv(envKey); ok {
		return val
	}
	return defaultValue
}
