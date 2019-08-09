#!/bin/bash
#
# Helper script for cross compiling local projects since it doesn't work with the current xgo version
# So we add the option to fix the pathing and use xgo within the xgo container which is an ubuntu based container image
#
# Required mounts:
# /src    - contains the source of the project
# /build  - contains the built binaries
#
# Needed environment variables:
# MODULE  - Module name/path to indicate where to copy it within the go root path (f.e. github.com/owner/name)
#
# Optional environment variables:
# PACKAGE - option to define the path to the main directory within the project
# TARGETS - option to define the build targets of xgo
# SOURCE  - option to set the repository source (branch/tag/commit)
# OUT     - option to set built binary prefix

# stop execution on errors to prevent wrong builds
set -e

# copy the mounted source into the final source
mkdir -p "${GOPATH}/src/${MODULE}" && cp -a "/src/." "${GOPATH}/src/${MODULE}" &&
  cd "${GOPATH}/src/${MODULE}"

if [[ ${SOURCE} ]]; then
  # clean up repository if a checkout source is set
  git reset --hard HEAD
  git clean -f -d
  # checkout repository to the passed source
  git checkout "${SOURCE}"
fi

params=()
if [[ ${PACKAGE} ]]; then
  params+=("--pkg=${PACKAGE}")
fi

if [[ ${TARGETS} ]]; then
  params+=("--targets=${TARGETS}")
fi

# shellcheck disable=SC2034
CGO_ENABLED=1 && xgo "${params[@]}" --dest /build/ --out provider .

# move into build directory for after-build operations
cd /build/

if [[ ${OUT} ]]; then
  # rename prefix of files in build directory if requested
  rename "s/^provider/${OUT}/" -- *
fi

# remove API level of windows builds from file name for auto updater
# which require the format {cmd}_{goos}_{goarch}{.ext} most of the times
rename "s/windows-4.0/windows/" -- *
