#!/bin/bash

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

[[ -n $DEBUG ]] && set -o xtrace
set -o errexit
set -o nounset
set -o pipefail

get_kubeconfig(){
  export CURL_OPTIONS=ks
  kdata=$("${__dir}/get-cluster.sh")
  unset CURL_OPTIONS
  IFS='%'
  echo "$kdata" | sed 's/.*\"kubeconfig\"\:\"//g' | sed 's/}}//g' | sed 's/\\n/\'$'\n''/g' | sed 's/\",.*//'
  unset IFS
}

main() {
  get_kubeconfig
}

main
