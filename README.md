# WebLogic Buildpack PoC [WIP]

This repository hosts the code for a rudimentary cloud native buildpack implementation for Oracle WebLogic Applications.

Oracle supports running WebLogic Server and Fusion Middleware on Kubernetes via the [WebLogic Kubernetes Operator](https://oracle.github.io/weblogic-kubernetes-operator/).

The operator supports different [ways](https://oracle.github.io/weblogic-kubernetes-operator/userguide/managing-domains/choosing-a-model/) for supplying the WebLogic domain:

* Domain in PV: Locates WebLogic domain homes in a Kubernetes PersistentVolume (PV). This PV can reside in an NFS file system or other Kubernetes volume types.
* Domain in Image: Includes a WebLogic domain home in a container image.
* Model in Image: Includes WebLogic Deploy Tooling models and archives in a container image.

This buildpack supports only the [Model in Image with auxiliary images](https://oracle.github.io/weblogic-kubernetes-operator/userguide/managing-domains/model-in-image/auxiliary-images/).