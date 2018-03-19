package bluegreendeployment_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "k8s.io/bgd/pkg/apis/controller/v1alpha1"
	. "k8s.io/bgd/pkg/client/clientset_generated/clientset/typed/controller/v1alpha1"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement controller logic tests

var _ = Describe("BlueGreenDeployment controller", func() {
	var instance BlueGreenDeployment
	var expectedKey types.ReconcileKey
	var client BlueGreenDeploymentInterface

	BeforeEach(func() {
		instance = BlueGreenDeployment{}
		instance.Name = "instance-1"
		expectedKey = types.ReconcileKey{
			Namespace: "default",
			Name:      "instance-1",
		}
	})

	AfterEach(func() {
		client.Delete(instance.Name, &metav1.DeleteOptions{})
	})

	Describe("when creating a new object", func() {
		It("invoke the reconcile method", func() {
			after := make(chan struct{})
			ctrl.AfterReconcile = func(key types.ReconcileKey, err error) {
				defer func() {
					// Recover in case the key is reconciled multiple times
					defer func() { recover() }()
					close(after)
				}()
				defer GinkgoRecover()
				Expect(key).To(Equal(expectedKey))
				Expect(err).ToNot(HaveOccurred())
			}

			// Create the instance
			client = cs.ControllerV1alpha1().BlueGreenDeployments("default")
			_, err := client.Create(&instance)
			Expect(err).ShouldNot(HaveOccurred())

			// Wait for reconcile to happen
			Eventually(after, "10s", "100ms").Should(BeClosed())

			// INSERT YOUR CODE HERE - test conditions post reconcile
		})
	})
})
