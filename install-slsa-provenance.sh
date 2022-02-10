#!/usr/bin/env bash

set -e

# default to relative path if INSTALL_PATH is not set
INSTALL_PATH=${INSTALL_PATH:-$(realpath ./.slsa-provenance)}

mkdir -p "${INSTALL_PATH}"

VERSION=v0.6.0
RELEASE="https://github.com/philips-labs/slsa-provenance-action/releases/download/${VERSION}"

OS=${RUNNER_OS:-linux}
ARCH=${RUNNER_ARCH:-amd64}

echo "Installing slsa-provenance at ${INSTALL_PATH}/bin"

if [ "${OS}" == "Windows" ] ; then
  OS=windows
elif [ "${OS}" == "Linux" ] ; then
  OS=linux
fi

if [ "${ARCH}" == "X64" ] ; then
  ARCH=amd64
fi

mkdir -p "$INSTALL_PATH"

trap "popd >/dev/null" EXIT
pushd "$INSTALL_PATH" > /dev/null || exit

echo "Downloading slsa-provenance_${VERSION/v}_${OS}_${ARCH}.tar.gz"
curl -sLo slsa-provenance.tar.gz "$RELEASE/slsa-provenance_${VERSION/v}_${OS}_${ARCH}.tar.gz"

if [ -x "$(command -v cosign)" ] ; then
  echo "Downloading slsa-provenance_${VERSION/v}_${OS}_${ARCH}.tar.gz.sig"
  curl -sLo slsa-provenance.tar.gz.sig "$RELEASE/slsa-provenance_${VERSION/v}_${OS}_${ARCH}.tar.gz.sig"
  echo "Downloading cosign.pub"
  curl -sLo cosign.pub "$RELEASE/cosign.pub"

  cosign verify-blob --key cosign.pub --signature slsa-provenance.tar.gz.sig slsa-provenance.tar.gz
  rm slsa-provenance.tar.gz.sig cosign.pub
else
  echo >&2
  echo "  cosign binary not installed in PATH. Unable to verify signature" >&2
  echo >&2
fi

tar -xzf slsa-provenance.tar.gz slsa-provenance
rm slsa-provenance.tar.gz

# for testing purposes fall back to "$INSTALL_PATH/GITHUB_PATH"
echo "$INSTALL_PATH" >> "${GITHUB_PATH:-"$INSTALL_PATH/GITHUB_PATH"}"
