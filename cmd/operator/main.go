/*
Copyright 2017 The Kubernetes Authors.

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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"flag"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	tprv1 "github.com/lukebond/k8s-operator-skeleton/pkg/apis/tpr/v1"
	exampleclient "github.com/lukebond/k8s-operator-skeleton/pkg/client"
	examplecontroller "github.com/lukebond/k8s-operator-skeleton/pkg/controller"
)

var (
	kubeconfigPath string
	masterUrl      string
)

func main() {
	flag.StringVar(&kubeconfigPath, "kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	flag.StringVar(&masterUrl, "master", "", "API server address. Omit to run in-cluster using the service account token. Not recommended for production.")
	flag.Parse()

	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := buildConfig(kubeconfigPath, masterUrl)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// initialize third party resource if it does not exist
	err = exampleclient.CreateTPR(clientset)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		panic(err)
	} else if err != nil {
		fmt.Println("TPR already exists")
	} else {
		fmt.Println("TPR created")
	}

	// make a new config for our extension's API group, using the first config as a baseline
	exampleClient, exampleScheme, err := exampleclient.NewClient(config)
	if err != nil {
		panic(err)
	}

	// wait until TPR gets processed
	err = exampleclient.WaitForExampleResource(exampleClient)
	if err != nil {
		panic(err)
	}

	// start a controller on instances of our TPR
	controller := examplecontroller.ExampleController{
		ExampleClient: exampleClient,
		ExampleScheme: exampleScheme,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go controller.Run(ctx)

	for {
		// Poll until Example object is handled by controller and gets status updated to "Processed"
		err = exampleclient.WaitForExampleInstanceProcessed(exampleClient, "example1")
		if err == nil {
			break
		}
	}
	fmt.Print("PROCESSED\n")

	// Fetch a list of our TPRs
	exampleList := tprv1.ExampleList{}
	err = exampleClient.Get().Resource(tprv1.ExampleResourcePlural).Do().Into(&exampleList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", exampleList)

	// infinite loop as placeholder for long-running functionality
	for {
	}
}

func buildConfig(kubeconfig string, masterUrl string) (*rest.Config, error) {
	if kubeconfig != "" || masterUrl != "" {
		fmt.Println("Creating out-of-cluster config")
		return clientcmd.BuildConfigFromFlags(masterUrl, kubeconfig)
	}
	fmt.Println("Creating in-cluster config")
	return rest.InClusterConfig()
}
