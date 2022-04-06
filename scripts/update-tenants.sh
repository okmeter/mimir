#!/bin/sh

set -euo pipefail

mc config host add s3 "http://$S3_HOST:$S3_PORT" "$AWS_ACCESS_KEY_ID" "$AWS_SECRET_ACCESS_KEY"

mc ls --json "s3/$BUCKET_NAME/" | jq -sr '[ .[] | select(.type == "folder" and .key != "tenant.index.json/")   | .key | rtrimstr("/") ]' | mc pipe --json "s3/$BUCKET_NAME/tenant.index.json"
mc cat "s3/$BUCKET_NAME/tenant.index.json" | jq -r '"Discovered \(.|length) tenants"'
