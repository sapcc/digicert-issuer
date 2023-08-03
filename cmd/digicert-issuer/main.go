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
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = certmanagerv1.AddToScheme(scheme)
	_ = certmanagerv1beta1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var (
		namespace                string
		defaultProviderNamespace string
		metricsAddr              string
		enableLeaderElection     bool
		printVersionAndExit      bool
		syncPeriod,
		backoffDurationProvisionerNotReady time.Duration
		backoffDurationRequestPending time.Duration
		disableRootCA                 bool
	)

	flag.StringVar(&namespace, "namespace", "",
		"Confine operator to the given namespace.")

	flag.StringVar(&defaultProviderNamespace, "default-provider-namespace", getValueFromEnvironmentOrDefault("POD_NAMESPACE", "kube-system"),
		"Namespace to fall back if provider does not exists.")

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080",
		"The address the metric endpoint binds to.")

	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	flag.BoolVar(&printVersionAndExit, "version", false,
		"Print version and exit.")

	flag.DurationVar(&syncPeriod, "sync-period", 15*time.Minute,
		"Synchronization period after which caches are invalidated.")

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
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "9f7013bc.cloud.sap",
		Namespace:          namespace,
		SyncPeriod:         &syncPeriod,
	})
	handleError(err, "unable to start manager")

	err = (&certmanagerv1beta1controller.DigicertIssuerReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("DigicertIssuer"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr)
	handleError(err, "unable to initialize controller", "controller", "digicertIssuer")

	err = (&certmanagerv1beta1controller.CertificateRequestReconciler{
		Client:                             mgr.GetClient(),
		Log:                                ctrl.Log.WithName("controllers").WithName("CertificateRequest"),
		Scheme:                             mgr.GetScheme(),
		BackoffDurationProvisionerNotReady: backoffDurationProvisionerNotReady,
		BackoffDurationRequestPending:      backoffDurationRequestPending,
		DefaultProviderNamespace:           defaultProviderNamespace,
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
