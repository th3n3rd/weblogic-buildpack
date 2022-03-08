#!/bin/bash

set -e

SCRIPTS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

export DOMAIN="sample-domain"
export AUX_CONTAINER_IMAGE="weblogic-sample-app"
export WLS_CONTAINER_IMAGE="container-registry.oracle.com/middleware/weblogic:12.2.1.4"
export INGRESS_HOST="weblogic-$DOMAIN.example.com"
export INGRESS_VARIANT="nginx"
NAMESPACE="$DOMAIN-ns"

cd "$SCRIPTS_DIR/.."

set -o allexport
source .env
set +o allexport

echo "Creating a new deployment namespace (if does not exist)"
kubectl create namespace "$NAMESPACE" || true
kubectl label ns "$NAMESPACE" weblogic-operator=enabled || true

echo "Deploying weblogic secrets and domain"
envsubst < deployment/manifest.yaml.tpl | kubectl apply -n "$NAMESPACE" -f -

echo "Deploying ingress"
envsubst < "deployment/$INGRESS_VARIANT-ingress.yaml.tpl" | kubectl apply -n "$NAMESPACE" -f -

timeout --foreground 300s ./scripts/healthcheck.sh "http://$INGRESS_HOST/app/actuator/health"

echo "Running smoke test"
RESULT=$(curl -s "http://$INGRESS_HOST/app/")
if [[ ! "$RESULT" =~ "Hello World!" ]]; then
    echo "Smoke test failed"
    exit 1
fi

echo "The application has been successfully deployed!"
