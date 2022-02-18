#!/usr/bin/env bash

shopt -s expand_aliases

if [ -z "$NO_COLOR" ]; then
  alias log_info="echo -e \"\033[1;32mINFO\033[0m:\""
  alias log_error="echo -e \"\033[1;31mERROR\033[0m:\""
  alias log_warning="echo -e \"\033[1;33mWARN\033[0m:\""
else
  alias log_info="echo \"INFO:\""
  alias log_error="echo \"ERROR:\""
  alias log_warning="echo \"WARN:\""
fi

set -e

GITHUB_API=${GITHUB_API:-'https://api.github.com'}

# default to relative path if INSTALL_PATH is not set
INSTALL_PATH=${INSTALL_PATH:-$(realpath ./.slsa-provenance)}

mkdir -p "${INSTALL_PATH}"

VERSION=v0.7.2
RELEASE="https://github.com/philips-labs/slsa-provenance-action/releases/download/${VERSION}"

if [[ "$VERSION" == *-draft ]] ; then
  curl_args=(-H "Authorization: token $GITHUB_TOKEN")
  assets=$(curl "${curl_args[@]}" -s "${GITHUB_API}/repos/philips-labs/slsa-provenance-action/releases?per_page=10" | jq "map(select(.name == \"${VERSION}\"))" | jq -r '.[0].assets')
fi

function download {
  url="${2}"
  if [[ "$VERSION" == *-draft ]] ; then
    url="$(echo "${assets}" | jq "map(select(.name == \"$1\"))" | jq -r '.[0].url')"
    curl_args+=(-H 'Accept: application/octet-stream')
  fi
  log_info "Downloading ${1} from ${url}…"
  curl -sLo "${1}" --show-error "${curl_args[@]}" "${url}"
  echo
}

OS=${RUNNER_OS:-Linux}
ARCH=${RUNNER_ARCH:-X64}

case "${ARCH}" in
  X64)
    ARCH=amd64
  ;;
  ARM64)
    ARCH=arm64
  ;;
  *)
    log_error "unsupported ARCH ${ARCH}"
    exit 1
  ;;
esac

BINARY=slsa-provenance
case "${OS}" in
  Linux)
    OS=linux
    ARCHIVE="slsa-provenance_${VERSION/v}_${OS}_${ARCH}.tar.gz"
  ;;
  macOS)
    ARCHIVE="slsa-provenance_${VERSION/v}_${OS}_${ARCH}.tar.gz"
  ;;
  Windows)
    OS=windows
    ARCHIVE="slsa-provenance_${VERSION/v}_${OS}_${ARCH}.zip"
    BINARY="${BINARY}.exe"
  ;;
  *)
    log_error "unsupported OS ${OS}"
    exit 1
  ;;
esac

DOWNLOAD="${RELEASE}/${ARCHIVE}"

log_info "Installing ${BINARY} (${OS}/${ARCH}) at ${INSTALL_PATH}"
mkdir -p "$INSTALL_PATH"

trap "popd >/dev/null" EXIT
pushd "$INSTALL_PATH" > /dev/null || exit

download "${ARCHIVE}" "${DOWNLOAD}"

if [ -x "$(command -v cosign)" ] ; then
  download ${ARCHIVE}.sig "${DOWNLOAD}.sig"
  download cosign.pub "$RELEASE/cosign.pub"

  log_info "Verifying signature…"
  cosign verify-blob --key cosign.pub --signature "${ARCHIVE}.sig" "${ARCHIVE}"
  rm "${ARCHIVE}.sig" cosign.pub
else
  log_warning >&2
  log_warning "  cosign binary not installed in PATH. Unable to verify signature!" >&2
  log_warning >&2
  log_warning "  Consider installing cosign first, to be able to verify the signature!" >&2
  log_warning >&2
fi

log_info "extracting ${BINARY} from ${ARCHIVE}"
tar -xzf "${ARCHIVE}" "${BINARY}"
rm "${ARCHIVE}"

# for testing purposes fall back to "$INSTALL_PATH/GITHUB_PATH"
echo "$INSTALL_PATH" >> "${GITHUB_PATH:-"$INSTALL_PATH/GITHUB_PATH"}"
