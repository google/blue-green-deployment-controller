

package args

import (
	"github.com/kubernetes-sigs/kubebuilder/pkg/inject/args"
    "k8s.io/client-go/rest"

    "k8s.io/bgd/pkg/client/clientset_generated/clientset"
)

// InjectArgs are the arguments need to initialize controllers
type InjectArgs struct {
    args.InjectArgs

    Clientset *clientset.Clientset
}


// CreateInjectArgs returns new controller args
func CreateInjectArgs(config *rest.Config) InjectArgs {
    return InjectArgs{
        InjectArgs: args.CreateInjectArgs(config),
        Clientset: clientset.NewForConfigOrDie(config),
    }
}
