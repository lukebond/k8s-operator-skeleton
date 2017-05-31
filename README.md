# Kubernetes Operator Skeleton

Based on the _Third Party Resources Example_ in the `k8s.io/client-go`
codebase.

## Building

Dependencies are vendored, using [vndr](https://github.com/LK4D4/vndr).

```
$ make deps
```

To build:

```
$ make build
```

## Running

To run the operator outside the cluster, pointing at it, there are two
options: providing a kube config file or providing an API URL.

For the former, run with `--kubeconfig`, passing a path to a kube config file,
remembering to mount in the volume. E.g.:

```
$ docker run --rm --net=host -v /home/luke/.kube:/home/luke/.kube -v /home/luke/.minikube:/home/luke/.minikube quay.io/lukebond/k8s-operator-skeleton:v1.0.0 --kubeconfig=/home/luke/.kube/config
```

For the latter, pass a URL, e.g. like the following to point to a local
kube proxy:

```
$ docker run --rm --net=host quay.io/lukebond/k8s-operator-skeleton:v1.0.0 -master=http://127.0.0.1:8080
```

Omit both these arguments to run in-cluster and use the service account token.

## Development Workflow

I suggest having a minikube cluster running, with kubectl configured to point
thereto, and work like so:

- In a separate terminal in the current directory, run `minikube mount .` and
  keep it running
- Edit code
- Run `make build`, then `make save` and `make clean-minikube`
  - The last command will remove previous deployments from minikube so that
    you can start afresh
- Run `minikube ssh` then type `docker load -i /mount-9p/k8s-operator-skeleton.tar`
- Run `kubectl create -f kube/deployment.yaml` to launch the operator
- Run `kubectl create -f kube/example-cluster.yaml` to create an instance
  of the TPR
