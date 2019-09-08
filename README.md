# Pingdom Operator

[![tag-badge][]]() [![goreport-badge][]][goreport-target] [![pipeline-badge][]][pipeline-target]

> Kubernetes Operator that maintains resources in Pingdom.

An operator that monitors CRDs and creates, updates or deletes resources in
Pingdom (currently checks) to ensure they match the specification.

![pingdom-operator-demo][demo-gif]

Features:

* support for **check** resources ✓
  * supports `name`, `host`, `type`, `port`, `resolution`, `userids`, `url`
    and `encryption` parameters
  * supports pausing/un-pausing
* support for _HTTP_, _TCP_, _Ping_, _SMTP_, _POP3_ and _IMAP_ check types
  with common parameters
* per-resource credentials (allows maintaining multiple Pingdom accounts from
  a single Kubernetes installation)

Built with the help of [Kubebuilder][] framework.

## Installing

TL;DR

With `kubectl` configured and authenticated to run against your Kubernetes
cluster:

``` sh
make install deploy
```

Assuming you run this in a tree with a tagged commit or any commit on master
checked out, it will **install** the CRDs and then **deploy** the operator
controller application.

Container images are built automatically by the [CI Pipeline][pipeline-target]
in [GitLab][] for every commit, and published to a [registry][], with image
tags matching the format `PROJECT/BRANCH:COMMIT_ID_OR_TAG`.

The `deploy` target in Makefile will infer the correct image tag for current
commit, update the manifests and apply them using `kubectl apply`. This can be
used to easily deploy any commit that's already been built by CI by simply
checking it out locally and running `make deploy`.

Finally to deploy a local-only commit (e.g. in development) you need to build
the image, push it to a registry and then deploy, like so:

``` sh
make docker-build docker-push deploy

# you can also override the tag, e.g. to use a different registry
make IMG=docker.io/user/pingdom-operator:latest docker-build docker-push deploy
```

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

Then there are sample manifests in [config/samples/](config/samples/) directory
for different types of checks, which you will need to modify to point to your
secret and then you can apply them with:

``` sh
$ kubectl apply -f config/samples/observability_v1alpha1_check_http.yaml
check.observability.pingdom.mig4.gitlab.io/sample-1 created
$ kubectl apply -f config/samples/observability_v1alpha1_check_ping.yaml
check.observability.pingdom.mig4.gitlab.io/sample-2 created
```

If the controller manager is running you should then see it pick up the
resource and create a check on Pingdom. See a list of checks:

``` sh
$ kubectl get checks
NAME           ID        TYPE   STATUS   HOST
sample-1       5399365   http   up       wtfismyip.com
sample-2       5401834   ping   down     error-service.io
```

## Uninstalling

To remove all resources created by `make deploy` run `make destroy` which will
run `kubectl delete` on the [default manifests](config/default/).

## Contributing

See known issues on [GitLab][gl-issues] (or [GitHub][gh-issues]), if you found
one that's not on the list or have a suggestion for improvement, open a new
issue. If you can, fork and send an MR/PR, it will be appreciated 💖.

## Hacking

`make` will generate code and compile a binary in `./bin/manager`. Then
`make install` will generate the manifests and install them to a Kubernetes
cluster (needs `kubectl` configured correctly).

You can then run the manager locally, just `./bin/manager` or `make run`.

Run tests with `make test` or `make gtest` (using
[Ginkgo](http://onsi.github.io/ginkgo/) runner).

## Acknowledgements 👍

* [kubebuilder][]
* [go-pingdom][]

## License

[![license-badge][]](LICENSE)


[tag-badge]: https://img.shields.io/github/v/tag/mig4/pingdom-operator
[goreport-badge]: https://goreportcard.com/badge/gitlab.com/mig4/pingdom-operator
[goreport-target]: https://goreportcard.com/report/gitlab.com/mig4/pingdom-operator
[pipeline-badge]: https://gitlab.com/mig4/pingdom-operator/badges/master/pipeline.svg
[pipeline-target]: https://gitlab.com/mig4/pingdom-operator/pipelines
[demo-gif]: https://gitlab.com/mig4/pingdom-operator/uploads/d88c23b703d080e80a72356f0c27826e/pingdom-operator-demo.gif
[gitlab]: https://gitlab.com/
[registry]: https://gitlab.com/mig4/pingdom-operator/container_registry
[gl-issues]: https://gitlab.com/mig4/pingdom-operator/issues
[gh-issues]: https://github.com/mig4/pingdom-operator/issues
[kubebuilder]: https://github.com/kubernetes-sigs/kubebuilder
[go-pingdom]: https://github.com/russellcardullo/go-pingdom
[license-badge]: https://img.shields.io/github/license/mig4/pingdom-operator?style=for-the-badge
