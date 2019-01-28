#!/bin/bash

CLUSTER_API=${CLUSTER_API:-cluster-manager-api.cnct.io}
CLUSTER_API_PORT=${CLUSTER_API_PORT:-443}
CLUSTER_NAME=${CLUSTER_NAME:-aws-test-$(date +%s)}
K8S_VERSION=${K8S_VERSION:-1.11.5}
CMA_CALLBACK_URL=${CMA_CALLBACK_URL:-https://webhook.site/#/15a7f31c-5b57-41fc-bd70-a8dec0f56442}
CMA_CALLBACK_REQUESTID=${CMA_CALLBACK_REQUESTID:-12345}

# aws specific
AWS_REGION=${AWS_REGION:-us-east-2}
AWS_AVAILABILITY_ZONE=${AWS_AVAILABILITY_ZONE:-us-east-2a}
AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
AWS_NODE_TYPE=${AWS_NODE_TYPE:-m5.large}

[[ -n $DEBUG ]] && set -o xtrace
set -o errexit
set -o nounset
set -o pipefail

DATA=$(
  cat <<JSON
{
  "name": "${CLUSTER_NAME}",
  "provider": {
    "name": "aws",
    "k8s_version": "${K8S_VERSION}",
    "aws": {
      "data_center": {
        "region": "${AWS_REGION}",
        "availability_zones": [ "${AWS_AVAILABILITY_ZONE}"
        ]
      },
      "credentials": {
        "secret_key_id": "${AWS_ACCESS_KEY_ID}",
        "secret_access_key": "${AWS_SECRET_ACCESS_KEY}",
        "region": "${AWS_REGION}"
      },
      "resources": {
        "vpc_id": "",
        "security_group_id": "",
        "iam_role_arn": ""
      },
      "instance_groups": [
        {
          "type": "${AWS_NODE_TYPE}",
          "min_quantity": 1,
          "max_quantity": 10
        }
      ]
    },
    "high_availability": true,
    "network_fabric": "flannel"
  },
  "callback": {
    "url": "${CMA_CALLBACK_URL}",
    "request_id": "${CMA_CALLBACK_REQUESTID}"
  }
}
JSON
)

main() {
  curl -X POST \
    "https://${CLUSTER_API}:${CLUSTER_API_PORT}/api/v1/cluster" \
    -H 'Cache-Control: no-cache' \
    -H 'Content-Type: application/json' \
    -d "${DATA}" \
    -iks
}

main
