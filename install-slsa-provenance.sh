#!/usr/bin/env bash

shopt -s expand_aliases
if [ -z "$NO_COLOR" ]; then
  alias log_info="echo -e \"\033[1;32mINFO\033[0m:\""
  alias log_error="echo -e \"\033[1;31mERROR\033[0m:\""
else
  alias log_info="echo \"INFO:\""
  alias log_error="echo \"ERROR:\""
fi

set -e

# default to relative path if INSTALL_PATH is not set
INSTALL_PATH=${INSTALL_PATH:-$(realpath ./.slsa-provenance)}

mkdir -p "${INSTALL_PATH}"

VERSION=v0.7.0-rc
RELEASE="https://github.com/philips-labs/slsa-provenance-action/releases/download/${VERSION}"

OS=${RUNNER_OS:-linux}
ARCH=${RUNNER_ARCH:-amd64}

log_info "Installing slsa-provenance at ${INSTALL_PATH}"

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

log_info "Downloading slsa-provenance_${VERSION/v}_${OS}_${ARCH}.tar.gz"
curl -sLo slsa-provenance.tar.gz "$RELEASE/slsa-provenance_${VERSION/v}_${OS}_${ARCH}.tar.gz"

if [ -x "$(command -v cosign)" ] ; then
  log_info "Downloading slsa-provenance_${VERSION/v}_${OS}_${ARCH}.tar.gz.sig"
  curl -sLo slsa-provenance.tar.gz.sig "$RELEASE/slsa-provenance_${VERSION/v}_${OS}_${ARCH}.tar.gz.sig"
  log_info "Downloading cosign.pub"
  curl -sLo cosign.pub "$RELEASE/cosign.pub"

  log_info "Verifying signature…"
  cosign verify-blob --key cosign.pub --signature slsa-provenance.tar.gz.sig slsa-provenance.tar.gz
  rm slsa-provenance.tar.gz.sig cosign.pub
else
  log_error >&2
  log_error "  cosign binary not installed in PATH. Unable to verify signature" >&2
  log_error >&2
fi

log_info "extracting slsa-provenance from slsa-provenance.tar.gz"
tar -xzf slsa-provenance.tar.gz slsa-provenance
rm slsa-provenance.tar.gz

# for testing purposes fall back to "$INSTALL_PATH/GITHUB_PATH"
echo "$INSTALL_PATH" >> "${GITHUB_PATH:-"$INSTALL_PATH/GITHUB_PATH"}"
