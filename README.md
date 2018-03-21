# Blue-Green Deployment controller

**This is not an official Google product**

This repository implements a simple blue-green deployment controller based on [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) framework.

The controller maintains 2 ReplicaSets (blue and green) all the time, alternating between the colors for new rollouts.

## Running Locally

3 terminals are needed to run the controller locally (one for running local cluster, another for running the controller, and last one for interacting with the controller).
Following step-by-step commands assume downloaded release files have been extracted to a directory named `bgd`. 

```sh
### first terminal ###

# start a local cluster within "kubernetes/kubernetes" directory
hack/local-up-cluster.sh

### second terminal ###

# navigate to "bgd" directory
cd <path>/<to>/bgd

# compile the codes
GOBIN=$(pwd)/bin go install <path>/<to>/bgd/cmd/controller-manager

# run the controller
bin/controller-manager --kubeconfig /var/run/kubernetes/admin.kubeconfig

### third terminal ###

# navigate to "bgd" directory
cd <path>/<to>/bgd

# export kubeconfig
export KUBECONFIG=/var/run/kubernetes/admin.kubeconfig

# create BlueGreenDeployment object
kubectl create -f hack/sample/bluegreendeployment.yaml 

# check ReplicaSets, pods, and service created through the custom resource
kubectl get all
```

When the `BlueGreenDeployment (BGD)` object is created, the controller creates 2 ReplicaSets based on pod spec of the BGD object. `Blue` RS is active (has specified number of available pods) while `green` RS is inactive (0 available pods). The controller also creates a service pointing to the active `blue` RS.

## Details

The controller runs an infinite loop to check pod spec difference between the active RS and BGD object.

1. If there is a discrepancy, the controller checks whether inactive RS has matching pod spec.

2. If it is a match, the controller changes service to point to the inactive RS and makes the active RS inactive.

3. If it is not a match, the controller deletes the inactive RS, creates a new RS, and waits for all pods of the new RS to become available.

4. After the new RS becomes available, the controller points the service to the new RS and makes the active RS inactive.

You can edit the BGD object to roll out new deployment versions:

```sh
### third terminal ###

kubectl edit bluegreendeployment blue-green-deployment
```

## Cleanup

You can clean up the CRD with:

    kubectl delete crd bluegreendeployments.controller.google.com

CRD deletion cleans up the CRD, `bluegreendeployment` custom resource, and ReplicaSets.

You can also clean up the `bluegreendeployment` custom resource with:

    kubectl delete bluegreendeployment blue-green-deployment

Custom resource deletion cleans up the custom resource and ReplicaSets.

For now, you have to manually delete the `bgd-svc` service.

    kubectl delete service bgd-svc

## Limitations

The controller does not support some manual actions by the user, but this should not affect its main functionalities.
* When a ReplicaSet is deleted manually, the controller will not respawn it and this will break the controller.
* When an controller is turned off manually, all created resources will stay intact. The user has to manually delete all the resources before restarting the controller again.
