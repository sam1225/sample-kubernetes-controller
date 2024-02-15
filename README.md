# sample-kubernetes-controller
This sample kubernetes controller implements a custom resource named MyAppResource.

The controller is built using Kubebuilder framework (https://github.com/kubernetes-sigs/kubebuilder).

## Description
A popular open source app named podinfo (https://github.com/stefanprodan/podinfo) will be deployed and managed as part of this implementation along with its data source (Redis).
The lifecycle management of this app will be taken care by a CustomResourceDefinition (CRD) named MyAppResource.

## Getting Started

### Prerequisites
- go version v1.21.0+
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.
- Need to have **cluster-admin** access on the K8s cluster connected to.

### To Deploy on the cluster

**Clone source code:**

```sh
git clone https://github.com/sam1225/sample-kubernetes-controller.git
cd sample-kubernetes-controller
```

**Install the CRD:**

```sh
kubectl apply -f config/crd/bases/k8s.myappresources.io_myappresources.yaml
```

**Run controller (on a seperate tab):**

```sh
make run
```

**Install the custom resource:**

```sh
kubectl apply -f config/samples/k8s_v1alpha1_myappresource.yaml 
```

>**NOTE**: You can retrieve the installed resource using any of these names: myappresource / mr

**To look at the installed resources:**

```sh
kubectl get deploy,po,svc,mr --show-labels -l app.kubernetes.io/name=podinfo
```

**To edit the resource:**

```sh
kubectl edit mr podinfo
```

**To view in UI:**

```sh
kubectl port-forward svc/podinfo --address=0.0.0.0 8181:9898
```
Open the link in browser: http://localhost:8181

**To interact with the cache:**

```sh
curl -X POST -d "cache value..." http://localhost:8181/cache/mykey
curl http://localhost:8181/cache/mykey
```

### To Uninstall

```sh
kubectl apply -f config/samples/k8s_v1alpha1_myappresource.yaml
kubectl apply -f config/crd/bases/k8s.myappresources.io_myappresources.yaml
```

