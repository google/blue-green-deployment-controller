package inject

import (
	"k8s.io/bgd/pkg/client/informers_generated/externalversions"
	"k8s.io/bgd/pkg/inject/args"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	controllerv1alpha1 "k8s.io/bgd/pkg/apis/controller/v1alpha1"
	"k8s.io/bgd/pkg/controller/bluegreendeployment"
)

func init() {
	// Inject Informers
	SetInformers = func(arguments args.InjectArgs, factory externalversions.SharedInformerFactory) {
		arguments.ControllerManager.AddInformerProvider(&controllerv1alpha1.BlueGreenDeployment{}, factory.Controller().V1alpha1().BlueGreenDeployments())
	}

	// Inject Controllers
	Controllers = append(Controllers, bluegreendeployment.ProvideController)
	// Inject CRDs
	CRDs = append(CRDs, &controllerv1alpha1.BlueGreenDeploymentCRD)
	// Inject PolicyRules
	PolicyRules = append(PolicyRules, rbacv1.PolicyRule{
		APIGroups: []string{"controller.google.com"},
		Resources: []string{"*"},
		Verbs:     []string{"*"},
	})
	// Inject GroupVersions
	GroupVersions = append(GroupVersions, schema.GroupVersion{
		Group:   "controller.google.com",
		Version: "v1alpha1",
	})
}
