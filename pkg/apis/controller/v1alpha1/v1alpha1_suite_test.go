


package v1alpha1_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "github.com/kubernetes-sigs/kubebuilder/pkg/test"
    "k8s.io/client-go/rest"

    "k8s.io/bgd/pkg/inject"
    "k8s.io/bgd/pkg/client/clientset_generated/clientset"
)

var testenv *test.TestEnvironment
var config *rest.Config
var cs *clientset.Clientset

func TestV1alpha1(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecsWithDefaultAndCustomReporters(t, "v1 Suite", []Reporter{test.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
    testenv = &test.TestEnvironment{CRDs: inject.CRDs}

    var err error
    config, err = testenv.Start()
    Expect(err).NotTo(HaveOccurred())

    cs = clientset.NewForConfigOrDie(config)
})

var _ = AfterSuite(func() {
    testenv.Stop()
})
