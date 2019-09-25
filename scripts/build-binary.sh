#!/bin/sh
# Copyright 2019 The Morning Consult, LLC or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#         https://www.apache.org/licenses/LICENSE-2.0
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

set -e

ROOT=$( cd "$( dirname "${0}" )/.." && pwd )
cd "${ROOT}"

BIN_DIR="${ROOT}/bin"
TAG="${1}"
COMMIT="${2}"
REPO="${3}"

if [ "${REPO}" = "" ]; then
  echo "Project name (fourth argument) is missing. This should be the project name (e.g. \"github.com/example/project\"). Exiting."
  exit 1
fi

cd "${ROOT}"

# Set ldflags
version_ldflags="-X \"${REPO}/version.Date=$( date +"%b %d, %Y" )\""

if [ "${TAG}" != "" ]; then
  version_ldflags="${version_ldflags} -X \"${REPO}/version.Version=${TAG}\""
fi

if [ "${COMMIT}" != "" ]; then
  version_ldflags="${version_ldflags} -X \"${REPO}/version.Commit=${COMMIT}\""
fi

mkdir -p "${BIN_DIR}"

GO111MODULE=on CGO_ENABLED=0 go build \
	-installsuffix cgo \
	-a \
	-ldflags "-s -w ${version_ldflags}" \
	-o "${BIN_DIR}/$( basename ${REPO} )" \
	.
