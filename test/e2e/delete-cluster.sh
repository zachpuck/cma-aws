#!/bin/bash

CLUSTER_API=${CLUSTER_API:-cluster-manager-api.cnct.io}
CLUSTER_API_PORT=${CLUSTER_API_PORT:-443}
CLUSTER_NAME=${CLUSTER_NAME}
CMA_CALLBACK_URL=${CMA_CALLBACK_URL:-https://webhook.site/#/15a7f31c-5b57-41fc-bd70-a8dec0f56442}
CMA_CALLBACK_REQUESTID=${CMA_CALLBACK_REQUESTID:-12345}

# aws specific
AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
AWS_REGION=${AWS_REGION:-us-east-2}

[[ -n $DEBUG ]] && set -o xtrace
set -o errexit
set -o nounset
set -o pipefail

rawurlencode() {
  local string="${1}"
  local strlen=${#string}
  local encoded=""
  local pos c o

  for (( pos=0 ; pos<strlen ; pos++ )); do
     c=${string:$pos:1}
     case "$c" in
        [-_.~a-zA-Z0-9] ) o="${c}" ;;
        * )               printf -v o '%%%02x' "'$c"
     esac
     encoded+="${o}"
  done
  REPLY="${encoded}"
}


main() {
  
  rawurlencode "$AWS_SECRET_ACCESS_KEY"
  secret_encode="${REPLY}"
  rawurlencode "${CMA_CALLBACK_URL}"
  url_encode="${REPLY}"

  curl -X DELETE \
    "https://${CLUSTER_API}:${CLUSTER_API_PORT}/api/v1/cluster?name=${CLUSTER_NAME}&aws.secret_key_id=${AWS_ACCESS_KEY_ID}&aws.secret_access_key=${secret_encode}&aws.region=${AWS_REGION}&provider=aws&callback.url=${url_encode}&callback.request_id=${CMA_CALLBACK_REQUESTID}" \
    -H 'Cache-Control: no-cache' \
    -H 'Content-Type: application/json' \
    -iks
}

main
