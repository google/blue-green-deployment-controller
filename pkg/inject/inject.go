

package inject

import (
    "time"

    "github.com/kubernetes-sigs/kubebuilder/pkg/controller"
    "github.com/kubernetes-sigs/kubebuilder/pkg/inject/run"
    apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
    rbacv1 "k8s.io/api/rbac/v1"
    "k8s.io/apimachinery/pkg/runtime/schema"
    appsv1 "k8s.io/api/apps/v1"

    "k8s.io/bgd/pkg/inject/args"
    "k8s.io/bgd/pkg/client/informers_generated/externalversions"
    "k8s.io/client-go/informers"
)

var (
    CRDs = []*apiextensionsv1beta1.CustomResourceDefinition{}

    PolicyRules = []rbacv1.PolicyRule{}

    GroupVersions = []schema.GroupVersion{}


    // Controllers provides the controllers to run
    // Should be set by code generation in this package.
    Controllers  = []func(args args.InjectArgs) (*controller.GenericController, error){}

    RunningControllers = map[string]*controller.GenericController{}

    // SetInformers adds the informers for the apis defined in this project.
    // Should be set by code generation in this package.
    SetInformers func(args.InjectArgs, externalversions.SharedInformerFactory)
)

// RunAll starts all of the informers and Controllers
func RunAll(options run.RunArguments, arguments args.InjectArgs) error {
    if SetInformers != nil {
        factory := externalversions.NewSharedInformerFactory(arguments.Clientset, time.Minute * 5)
        SetInformers(arguments, factory)
    }

    kubernetesinformers := informers.NewSharedInformerFactory(arguments.KubernetesClientSet, 30 * time.Second)
    arguments.ControllerManager.AddInformerProvider(&appsv1.ReplicaSet{}, kubernetesinformers.Apps().V1().ReplicaSets())
    for _, fn := range Controllers {
        if c, err := fn(arguments); err != nil {
            return err
        } else {
            arguments.ControllerManager.AddController(c)
        }
    }
    arguments.ControllerManager.RunInformersAndControllers(options)
    <-options.Stop
    return nil
}

type Injector struct {}

func (Injector) GetCRDs() []*apiextensionsv1beta1.CustomResourceDefinition {return CRDs}
func (Injector) GetPolicyRules() []rbacv1.PolicyRule {return PolicyRules}
func (Injector) GetGroupVersions() []schema.GroupVersion {return GroupVersions}
