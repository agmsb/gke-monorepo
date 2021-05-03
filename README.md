# gke-monorepo

Little tweaks but almost entirely based on [viglesias/sample-app/tree/base-pattern](https://github.com/viglesiasce/sample-app/tree/base-pattern).

`apps` holds apps, Dockerfile, skaffold.yaml for a `frontend` and `backend` app. The `frontend` skaffold requires the `backend` skaffold. Both require the `base-image` skaffold. The `frontend` app refers to k8s manifests upstream. 
`infrastructure` holds Dockerfile and `base-image` skaffold to build base images for `go` apps, as well as k8s manifests for the `backend` app should they be in the same repo. 

## skaffold profiles

`skaffold.yaml` for the frontend will use a `ClusterIP` service for `pre-production` with `portForward` enabled. For `production` it will create a service `Type: LoadBalancer`. 

## TODO

* add policy
* incorporate secrets manager
* incorporate env repo pattern
* incorporate remote GCB builds in `/ci`
