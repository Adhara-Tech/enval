#!/bin/bash
set -eo pipefail

TAG=$1
VERSION_NAME=$(echo $TAG | sed 's/^v//')

echo "${VERSION_NAME}" > release_message.txt
# whiteline to separate release title from body
echo >> release_message.txt

changelog-parser CHANGELOG.md | jq -r ".versions[] | select(.version == \"${VERSION_NAME}\") | .body" >> release_message.txt


hub release create -F release_message.txt \
  -a bin/enval_linux_amd64 \
  -a bin/enval_darwin_amd64 \
  -a bin/enval_windows_amd64.exe \
  $TAG
