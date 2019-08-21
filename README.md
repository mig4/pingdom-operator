# Pingdom Operator

A Kubernetes Operator that maintains resources (currently checks) in Pingdom.

Built with the help of [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)
framework.

## Hacking

`make` will generate code and compile a binary in `./bin/manager`. Then
`make install` will generate the manifests and install them to a Kubernetes
cluster (needs `kubectl` configured correctly).

You can then run the manager locally, just `./bin/manager` or `make run`.

## Usage

Once the CRD is installed (with `make install`) you can deploy Pingdom checks.

First thing you'll need is a
[secret](https://kubernetes.io/docs/concepts/configuration/secret/) with
Pingdom API credentials, e.g.:

``` sh
kubectl create secret generic my-pd-secret \
  --from-literal=username=PINGDOM_USER \
  --from-literal=password=PINGDOM_PASS
```

Or create a YAML manifest and apply it.

Then there is a sample manifest in
[config/samples/observability_v1alpha1_check.yaml](config/samples/observability_v1alpha1_check.yaml)
you will need to modify it to point to your secret and then apply it with:

``` sh
$ kubectl apply -f config/samples/observability_v1alpha1_check.yaml
check.observability.pingdom.mig4.gitlab.io/check-sample created
```

If the controller manager is running you should then see it pick up the
resource and create a check on Pingdom. See a list of checks:

``` sh
$ kubectl get checks
NAME           ID        TYPE   STATUS   PAUSED
check-sample   5392656   http   up       false
```
