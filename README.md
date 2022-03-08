# WebLogic Buildpack PoC [WIP]

This repository hosts the code for a rudimentary cloud native buildpack implementation for Oracle WebLogic Applications.

Oracle supports running WebLogic Server and Fusion Middleware on Kubernetes via the [WebLogic Kubernetes Operator](https://oracle.github.io/weblogic-kubernetes-operator/).

The operator supports different [ways](https://oracle.github.io/weblogic-kubernetes-operator/userguide/managing-domains/choosing-a-model/) for supplying the WebLogic domain:

* Domain in PV: Locates WebLogic domain homes in a Kubernetes PersistentVolume (PV). This PV can reside in an NFS file system or other Kubernetes volume types.
* Domain in Image: Includes a WebLogic domain home in a container image.
* Model in Image: Includes WebLogic Deploy Tooling models and archives in a container image.

This buildpack supports only the [Model in Image with auxiliary images](https://oracle.github.io/weblogic-kubernetes-operator/userguide/managing-domains/model-in-image/auxiliary-images/).

## Caveats

1. Need to fix the Oracle WebLogic Operator installation as the one of the scripts it includes has a [bug](https://github.com/oracle/weblogic-kubernetes-operator/issues/2819#issuecomment-1060816388)
This is because the operator will spawn an init container based on the auxiliary image provided, which is based on Ubuntu Bionic,
and some scripts are written with non-POSIX compliant instructions and with the assumption that the OS shell (`/bin/sh`) is pointing to bash or equivalent.
This is how the operator can be fixed:
* Download the operator source code (by the way this step is also part their standard installation)
* Replace `#!/bin/sh` with `#!/bin/bash` in the `operator/src/main/resources/scripts/auxImage.sh` script
* Run `mvn clean package` in the root
* Run `./buildDockerImage.sh -t <your-custom-image-name>`
* Follow the installation instruction remembering to use `<your-custom-image-name>` image as reference

3. If the application will be deployed on a dedicated k8s cluster then is recommended to re-tag and publish the Oracle WebLogic Server image
into a container registry already accessible by the cluster, this is because there are some particular condition to be met before even attempting
the download of the base image.
* Requires a Oracle user account
* Requires that account to have agreed on the T&C for that specific image
* Requires to use the account credentials as pullSecrets.
Is far easier to pull the image locally and publish it onto a different registry.