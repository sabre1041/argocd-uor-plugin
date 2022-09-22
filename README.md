# argocd-uor-plugin

[Argo CD Plugin](https://argo-cd.readthedocs.io/en/stable/user-guide/config-management-plugins) to retrieve assets from [Universal Object Reference (UOR)](https://universalreference.io) Collections.

## Installing the plugin in your environment

The plugin is designed for Argo CD version 2.4+ and is designed to leverage the _sidecar_ model of a Config Management Plugin (CMP) and the [Argo CD Operator](https://argocd-operator.readthedocs.io).

First, add the CMP ConfigMap to the namespace containing Argo CD:

```shell
kubectl apply -f manifests/cmp-plugin.yaml
```

Next, patch the repo server containers to make the plugin available to Argo CD:

```shell
kubectl apply -f manifests/argocd.yaml
```

NOTE: If an existing `ArgoCD` resource has been defined, be sure to match the name of the existing resource.

The plugin has been applied successfully when the repo-server pod has completed its' initialization and includes the necessary sidecar containers

## Getting Started

The following is a simple demonstration using the plugin.

NOTE: The getting started is currently designed for use within an OpenShift environment.

The `simple` example that will be demonstrated as part of this getting started is located in the [examples/simple] directory.

Navigate to the [examples/simple] directory:

```shell
cd examples/simple
```

A _Collection_ has previously been built and is available within a publicly accessible registry. The assets related to the collection are located in the [collection](examples/simple/collection) directory along with a Dataset Configuration file ([dataset-config.yaml](examples/simple/dataset-config.yaml)).

Two Argo CD applications are available to demonstrate the retrieval of Kubernetes assets within a collection. The collection contains Kubernetes resources that have metadata associated to development (dev) and production (prod) environments through the use of _attributes_.

The `uor-sample-dev` application makes use of a plugin environment variable `ATTRIBUTE_QUERY` that defines the location of a `AttributeQuery` that retrieves resources from the collection that contain the attribute `dev: true`.

The Argo CD Application defines two plugin environment variable that it makes use of at runtime:

* Location of the UOR collection
* Location of the Attribute file within the source

Each Argo CD Application references a different Attribute Query resource located in the [source](examples/simple/source) directory.

Apply both Argo Applications to the cluster.

```
kubectl apply -f argocd
```

Note: Be sure to specify the namespace Argo CD is deployed within

IMPORTANT: Be sure that the Argo CD controller has the necessary permissions to create namespaces and resources within the newly created namespaces.

Once Argo CD synchronizes the repository, two namespaces will be created containing a single mircoservice application.

* simple-uor-dev
* simple-uor-prod

A _Route_ exposes the application so that they can be viewed in a browser. Open a browser to view the application. Note that the background color in each deployment is different (blue vs red) that differentiates the environments as the Argo CD plugin retrieved from the UOR collection as defined by the presence of the `ENVIRONMENT` environment variable in the production Deployment.  

Verify the application can be viewed in a browser to confirm the end to end capability of the Argo CD plugin.
