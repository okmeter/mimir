#!/bin/bash

set -uo pipefail

cat /mnt/envoy.yaml | envsubst > /etc/envoy/envoy.yaml

exit $?
