package main

import (
	"flag"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	"github.com/petomalina/krane/operator/pkg/apis"
	"github.com/petomalina/krane/operator/pkg/controller"
	"log"
	"runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

const (
	api            = "krane.petomalina.com/v1alpha1"
	deploymentKind = "Canary"
	resyncPeriod   = 0

	version = "0.0.1"
)

func printVersion() {
	log.Printf("Go Version: %s", runtime.Version())
	log.Printf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	log.Printf("operator-sdk Version: %v", version)
}

func main() {
	printVersion()
	flag.Parse()

	// get all namespaces to watch
	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		log.Fatalf("failed to get watch namespace: %v", err)
	}

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{Namespace: namespace})
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Fatal(err)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Fatal(err)
	}

	log.Print("Starting the Cmd.")

	// Start the Cmd
	log.Fatal(mgr.Start(signals.SetupSignalHandler()))

	//handler := &internal.Handler{}
	//
	//sdk.Watch(api, deploymentKind, namespace, resyncPeriod)
	//sdk.Handle(handler)
	//sdk.Run(context.Background())
}
