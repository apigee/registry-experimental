#!/bin/sh
# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script downloads and installs the latest versions of the registry-experimental tools.
# Binaries are installed in the user's home directory in $HOME/.registry/bin.
set -e

# Determine the operating system.
OS="$(uname)"
if [ "${OS}" = "Darwin" ] ; then
  OSEXT="Darwin"
else
  OSEXT="linux"
fi

# Determine the latest registry-experimental version by version number ignoring alpha, beta, and rc versions.
if [ "${REGISTRY_EXPERIMENTAL_VERSION}" = "" ] ; then
  REGISTRY_EXPERIMENTAL_VERSION="$(curl -sL https://github.com/apigee/registry-experimental/releases/latest | \
                  grep -i release | grep -o 'v[0-9].[0-9]*.[0-9]*' | tail -1)"
  REGISTRY_EXPERIMENTAL_VERSION="${REGISTRY_EXPERIMENTAL_VERSION##*/}"
fi

LOCAL_ARCH=$(uname -m)
if [ "${TARGET_ARCH}" ]; then
    LOCAL_ARCH=${TARGET_ARCH}
fi

case "${LOCAL_ARCH}" in
  x86_64|amd64|arm64)
    REGISTRY_ARCH=amd64
    ;;
  armv8*|aarch64*)
    REGISTRY_ARCH=arm64
    ;;
  *)
    echo "This system's architecture, ${LOCAL_ARCH}, isn't supported"
    exit 1
    ;;
esac

if [ "${REGISTRY_EXPERIMENTAL_VERSION}" = "" ] ; then
  printf "Unable to get latest registry-experimental version. Set REGISTRY_EXPERIMENTAL_VERSION env var and re-run. For example: export REGISTRY_EXPERIMENTAL_VERSION=v1.104"
  exit 1;
fi

# Download the registry-experimental release archive.
tmp=$(mktemp -d /tmp/registry-experimental.XXXXXX)
NAME="registry-experimental_$REGISTRY_EXPERIMENTAL_VERSION"

cd "$tmp" || exit
FILENAME="registry-experimental_${REGISTRY_EXPERIMENTAL_VERSION##v}_${OSEXT}_${REGISTRY_ARCH}.tar.gz"
URL="https://github.com/apigee/registry-experimental/releases/download/${REGISTRY_EXPERIMENTAL_VERSION}/${FILENAME}"
echo $URL

download_archive() {
  printf "\nDownloading %s from %s ...\n" "$NAME" "$URL"
  if ! curl -o /dev/null -sIf "$URL"; then
    printf "\n%s is not found, please specify a valid REGISTRY_EXPERIMENTAL_VERSION and TARGET_ARCH\n" "$URL"
    exit 1
  fi
  curl -fsLO "$URL"
  tar xzf "${FILENAME}"
}

download_archive

printf ""
printf "\nregistry-experimental %s Download Complete!\n" "$REGISTRY_EXPERIMENTAL_VERSION"
printf "\n"

# Setup registry-experimental
cd "$HOME" || exit
mkdir -p "$HOME/.registry/bin"
mv "${tmp}/registry-experimental" "$HOME/.registry/bin"
mv "${tmp}/registry-connect" "$HOME/.registry/bin"
printf "Copied registry-experimental tools into the $HOME/.registry/bin folder.\n"
chmod +x "$HOME/.registry/bin/registry-experimental"
chmod +x "$HOME/.registry/bin/registry-connect"

# Print message
printf "\n"
printf "Add the registry tools to your path with:"
printf "\n"
printf "  export PATH=\$PATH:\$HOME/.registry/bin \n"
printf "\n"
