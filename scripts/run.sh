#!/bin/bash

set -e

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

cd "$SCRIPT_DIR/.."

echo "Running the cloud native buildpack by trying to build a sample app"
pack build weblogic-sample-app \
  --path test/sample-app \
  --builder paketobuildpacks/builder:base \
  --buildpack paketo-buildpacks/bellsoft-liberica \
  --buildpack paketo-buildpacks/syft \
  --buildpack paketo-buildpacks/maven \
  --buildpack . \
  --verbose

docker run -it --rm weblogic-sample-app /bin/sh

echo "It works!"
