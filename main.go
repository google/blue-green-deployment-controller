/*
Copyright 2018 The Kubernetes Authors.

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
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"k8s.io/client-go/rest"

	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	demov1 "k8s.io/bgd-controller/pkg/apis/demo/v1"
	bgdclientset "k8s.io/bgd-controller/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	svc                          *corev1.Service
	activeRS, inactiveRS         *extensionsv1beta1.ReplicaSet
	bgdPodTemplateSpecAnnotation = "demo.google.com/bgd-pod-template-spec"
	pollInterval                 = 100 * time.Millisecond
	pollTimeout                  = 5 * time.Second
)

func createReplicaSetsAndService(bgd *demov1.BlueGreenDeployment, kubeClient *kubernetes.Clientset) error {
	rsClient := kubeClient.ExtensionsV1beta1().ReplicaSets(bgd.Namespace)
	encodedPodTemplateSpec, err := json.Marshal(bgd.Spec.Template.Spec)
	if err != nil {
		return fmt.Errorf("failed to encode bgd pod template spec: %v", err)
	}
	encodedPodTemplateSpecStr := fmt.Sprintf("%s", encodedPodTemplateSpec)

	// create blue RS if it doesn't exist
	blueRS, err := rsClient.Create(newReplicaSet("blue-rs", "blue", bgd, false))
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create blue RS: %v", err)
	} else if err == nil {
		fmt.Printf("created blue RS\n\n")
		blueRS.Annotations[bgdPodTemplateSpecAnnotation] = encodedPodTemplateSpecStr
		activeRS = blueRS
	}

	// create green RS if it doesn't exist
	greenRS, err := rsClient.Create(newReplicaSet("green-rs", "green", bgd, true))
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create green RS: %v", err)
	} else if err == nil {
		fmt.Printf("created green RS\n\n")
		greenRS.Annotations[bgdPodTemplateSpecAnnotation] = encodedPodTemplateSpecStr
		inactiveRS = greenRS
	}

	// create service points to current active RS if it doesn't exist
	// service will only be created once; initially point to blue RS
	svc, err = kubeClient.CoreV1().Services(bgd.Namespace).Create(newService(bgd.Namespace))
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create service: %v", err)
	} else if err == nil {
		fmt.Printf("created service %q\n\n", svc.Name)
	}

	return nil
}

// 1. delete the inactive RS to give way to the room for creating a new RS
// 2. create and scale up the new RS
// 3. point the service to the new RS
// 4. scale down previous active RS to make it inactive
func createNewActiveRSToReplaceInactiveRS(bgd *demov1.BlueGreenDeployment, kubeClient *kubernetes.Clientset) error {
	rsClient := kubeClient.ExtensionsV1beta1().ReplicaSets(bgd.Namespace)
	newActiveRSName := inactiveRS.Name
	newActiveRSColor := inactiveRS.Spec.Template.Labels["color"]
	err := rsClient.Delete(inactiveRS.Name, &metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete inactive RS %q: %v", newActiveRSName, err)
	}

	if err := waitRSDeletionToComplete(inactiveRS, rsClient); err != nil {
		fmt.Errorf("timeout waiting inactive RS %q deletion: %v", inactiveRS.Name, err)
	}

	newRS, err := rsClient.Create(newReplicaSet(newActiveRSName, newActiveRSColor, bgd, false))
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create a new RS %q: %v", newRS.Name, err)
	} else if err == nil {
		fmt.Printf("created a new RS %q\n\n", newRS.Name)
	}

	newActiveRS, err := scaleRS(newRS, bgd.Spec.Replicas, rsClient)
	if err != nil {
		return fmt.Errorf("failed to scale up new RS %q: %v", newRS.Name, err)
	}

	encodedPodTemplateSpec, err := json.Marshal(bgd.Spec.Template.Spec)
	if err != nil {
		return fmt.Errorf("failed to encode bgd pod template spec: %v", err)
	}
	newActiveRS.Annotations[bgdPodTemplateSpecAnnotation] = fmt.Sprintf("%s", encodedPodTemplateSpec)

	svc, err = updateService(svc.Name, bgd.Namespace, kubeClient, func(service *corev1.Service) {
		updatedLabels := newActiveRS.Spec.Template.Labels
		service.Labels = updatedLabels
		service.Spec.Selector = updatedLabels
	})
	if err != nil {
		return fmt.Errorf("failed to update service %q to point to new active RS %q: %v", svc.Name, newActiveRS.Name, err)
	}

	newInactiveRS, err := scaleRS(activeRS, 0, rsClient)
	if err != nil {
		return fmt.Errorf("failed to scale down active RS %q: %v", activeRS.Name, err)
	}
	newInactiveRS.Annotations[bgdPodTemplateSpecAnnotation] = activeRS.Annotations[bgdPodTemplateSpecAnnotation]

	activeRS = newActiveRS
	inactiveRS = newInactiveRS
	fmt.Printf("active RS: %q, inactive RS: %q, service label selectors: %q\n\n", activeRS.Name, inactiveRS.Name, svc.Spec.Selector)
	return nil
}

// 1. scale up the inactive RS
// 2. point the service to the newly active RS
// 3. scale down previous active RS to make it inactive
func switchToInactiveRS(bgd *demov1.BlueGreenDeployment, kubeClient *kubernetes.Clientset) error {
	rsClient := kubeClient.ExtensionsV1beta1().ReplicaSets(bgd.Namespace)
	newActiveRS, err := scaleRS(inactiveRS, bgd.Spec.Replicas, rsClient)
	if err != nil {
		return fmt.Errorf("failed to scale up inactive RS %q: %v", inactiveRS.Name, err)
	}
	newActiveRS.Annotations[bgdPodTemplateSpecAnnotation] = inactiveRS.Annotations[bgdPodTemplateSpecAnnotation]

	svc, err = updateService(svc.Name, bgd.Namespace, kubeClient, func(service *corev1.Service) {
		updatedLabels := newActiveRS.Spec.Template.Labels
		service.Labels = updatedLabels
		service.Spec.Selector = updatedLabels
	})
	if err != nil {
		return fmt.Errorf("failed to update service %q to point to new active RS %q: %v", svc.Name, newActiveRS.Name, err)
	}

	newInactiveRS, err := scaleRS(activeRS, 0, rsClient)
	if err != nil {
		return fmt.Errorf("failed to scale down active RS %q: %v", activeRS.Name, err)
	}
	newInactiveRS.Annotations[bgdPodTemplateSpecAnnotation] = activeRS.Annotations[bgdPodTemplateSpecAnnotation]

	activeRS = newActiveRS
	inactiveRS = newInactiveRS
	fmt.Printf("active RS: %q, inactive RS: %q, service label selectors: %q\n\n", activeRS.Name, inactiveRS.Name, svc.Spec.Selector)
	return nil
}

func syncReplicaSetsAndService(bgd *demov1.BlueGreenDeployment, kubeClient *kubernetes.Clientset) error {
	rss, err := kubeClient.ExtensionsV1beta1().ReplicaSets(bgd.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to get RSs: %v", err)
	}

	// create RSs and service in the beginning
	if len(rss.Items) == 0 {
		return createReplicaSetsAndService(bgd, kubeClient)
	}

	encodedPodTemplateSpec, err := json.Marshal(bgd.Spec.Template.Spec)
	if err != nil {
		return fmt.Errorf("failed to encode bgd pod template spec: %v", err)
	}
	encodedPodTemplateSpecStr := fmt.Sprintf("%s", encodedPodTemplateSpec)

	// check whether both current active RS and the BGD object have the same pod template spec
	if activeRS.Annotations[bgdPodTemplateSpecAnnotation] != encodedPodTemplateSpecStr {
		fmt.Printf("bgd pod template spec changed\n\n")
		if inactiveRS.Annotations[bgdPodTemplateSpecAnnotation] != encodedPodTemplateSpecStr {
			return createNewActiveRSToReplaceInactiveRS(bgd, kubeClient)
		} else {
			return switchToInactiveRS(bgd, kubeClient)
		}
	}
	return nil
}

func resync(bgdClient *bgdclientset.Clientset, kubeClient *kubernetes.Clientset) error {
	bgd, err := bgdClient.DemoV1().BlueGreenDeployments("default").Get("blue-green-deployment", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("bgd object is not created yet or has been deleted\n\n")
		} else {
			fmt.Errorf("error getting bgd object: %v", err)
		}
	} else {
		if err := syncReplicaSetsAndService(bgd, kubeClient); err != nil {
			fmt.Errorf("failed to update replicaset and service: %v", err)
		}
	}
	return nil
}

// GetClientConfig returns rest config, if path not specified assume in cluster config
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func main() {
	kubeconf := flag.String("kubeconf", "admin.conf", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	config, err := GetClientConfig(*kubeconf)
	if err != nil {
		panic(fmt.Errorf("failed to get client config: %v", err))
	}

	// Create a CRD client interface
	bgdClient, err := bgdclientset.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("error building CRD clientset: %v", err))
	}

	// Create a kubernetes client interface
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("Error building kubernetes clientset: %v", err))
	}

	for {
		if err := resync(bgdClient, kubeClient); err != nil {
			fmt.Errorf("sync error: %v", err)
		}
		time.Sleep(3 * time.Second)
	}
}
