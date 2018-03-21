# Blue-Green Deployment controller

**WARNING: Please refer to kubebuilder branch for newest codes. This branch will be deleted in future.**

**This is not an official Google product**

This repository implements a simple blue-green deployment controller using a CustomResourceDefinition (CRD). The controller maintains 2 ReplicaSets (blue and green) all the time, alternating between the colors for new rollouts.

## Running Locally

3 terminals are needed to run the controller locally (one for running local cluster, another for running the controller, and last one for interacting with the controller).

```sh
### first terminal ###

# install "bgd-controller" directory alongside with "kubernetes" directory

# navigate to "bgd-controller" directory
cd bgd-controller

# copy whole directory into main repo
cp -r . ../kubernetes/staging/src/k8s.io/bgd-controller

# navigate to "kubernetes" directory
cd ../kubernetes

# create a symlink in vendor package
ln -s ../../staging/src/k8s.io/bgd-controller vendor/k8s.io/bgd-controller

# start a local cluster
hack/local-up-cluster.sh

### second terminal ###

# navigate to "bgd-controller" directory
cd kubernetes/staging/src/k8s.io/bgd-controller

# run the controller; kubeconfig is not required if operating in-cluster
go run *.go -kubeconf=/var/run/kubernetes/admin.kubeconfig

### third terminal ###

# navigate to "bgd-controller" directory
cd kubernetes/staging/src/k8s.io/bgd-controller

# set up kubeconfig
export KUBECONFIG=/var/run/kubernetes/admin.kubeconfig

# create a CustomResourceDefinition
kubectl create -f crd.yaml

# create a BlueGreenDeployment custom resource object
kubectl create -f bgd.yaml

# check ReplicaSets, pods, and service created through the custom resource
kubectl get all
```

When the `BlueGreenDeployment (BGD)` object is created, the controller creates 2 ReplicaSets based on pod template spec of the BGD object. `Blue` RS is active (has specified number of available pods) while `green` RS is inactive (0 available pods). The controller also creates a service pointing to the active `blue` RS.

## Details

The controller runs an infinite loop to check pod template spec difference between the active RS and BGD object.
a. If there is a discrepancy, the controller checks whether inactive RS has matching pod template spec.
b. If it is a match, the controller changes service to point to the inactive RS and makes the active RS inactive.
c. If it is not a match, the controller deletes the inactive RS, creates a new RS that has snapshot of jsonified version of the BGD’s current pod template spec stored in its `demo.google.com/bgd-pod-template-spec` annotation, and waits asynchronously through periodic polling for the new RS to become available.
d. If the new RS becomes available, the controller changes service to point to the new RS and makes the active RS inactive.

You can edit the BGD object using `kubectl`:

```sh
### third terminal ###

kubectl edit bluegreendeployment blue-green-deployment
```

## Cleanup

You can clean up the CRD with:

    kubectl delete crd bluegreendeployments.demo.google.com

CRD deletion cleans up the CRD, `bluegreendeployment` custom resource, and ReplicaSets.

You can also clean up the `bluegreendeployment` custom resource with:

    kubectl delete bluegreendeployment blue-green-deployment

Custom resource deletion cleans up the custom resource and ReplicaSets.

For now, you have to manually delete the `bgd-svc` service.

    kubectl delete service bgd-svc

_Note: The `bgd-svc` service **MUST** be deleted before restarting the custom controller to prevent a runtime error_.

## Limitations

The controller does not support some manual actions by the user, but this should not affect its main functionalities.
* When a ReplicaSet is deleted manually, the controller will not respawn it and this will break the controller.
* When an controller is turned off manually, all created resources will stay intact. The user has to manually delete all the resources before restarting the controller again, else there will be conflicts.

## References

* [metacontroller](https://github.com/GoogleCloudPlatform/kube-metacontroller/)
* [sample-controller](https://github.com/kubernetes/kubernetes/tree/master/staging/src/k8s.io/sample-controller)
* [kube-crd](https://github.com/yaronha/kube-crd)
