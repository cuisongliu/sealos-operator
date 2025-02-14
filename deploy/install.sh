#!/bin/bash
NAMESPACE=${NAMESPACE:-"sealos-operator"}
HELM_OPTS=${HELM_OPTS:-""}
helm upgrade --install sealos-operator charts/sealos-operator --namespace "${NAMESPACE}" --create-namespace ${HELM_OPTS}
