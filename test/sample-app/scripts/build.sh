#!/bin/bash

set -e

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
export AUX_CONTAINER_IMAGE="weblogic-sample-app"

cd "$SCRIPT_DIR/.."

set -o allexport
source .env
set +o allexport

echo "Deploying weblogic sample application"
pack build "$AUX_CONTAINER_IMAGE" \
  --path . \
  --builder paketobuildpacks/builder:base \
  --buildpack paketo-buildpacks/bellsoft-liberica \
  --buildpack paketo-buildpacks/syft \
  --buildpack paketo-buildpacks/maven \
  --buildpack ../.. \
  --verbose

echo "Done!"
