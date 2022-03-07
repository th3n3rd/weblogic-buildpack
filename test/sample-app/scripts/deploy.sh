#!/bin/bash

set -e

SCRIPTS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
export DOMAIN=${1:-sample-domain}
export CONTAINER_IMAGE=${2:-"weblogic-sample-app"}
NAMESPACE="$DOMAIN-ns"

cd "$SCRIPTS_DIR/.."

echo "Creating a new deployment namespace (if does not exist)"
kubectl create namespace "$NAMESPACE" || true
kubectl label ns "$NAMESPACE" weblogic-operator=enabled || true

echo "Deploying weblogic secrets and domain"
envsubst < deployment/manifest.yaml.tpl | kubectl apply -n "$NAMESPACE" -f -

echo "Deploying ingress"
envsubst < deployment/ingress.yaml.tpl | kubectl apply -n "$NAMESPACE" -f -

timeout --foreground 300s ./scripts/healthcheck.sh http://sample-app.weblogic.k8s/actuator/health

echo "Running smoke test"
RESULT=$(curl -s http://sample-app.weblogic.k8s)
if [[ ! "$RESULT" =~ "Hello World!" ]]; then
    echo "Smoke test failed"
    exit 1
fi

echo "The application has been successfully deployed!"
