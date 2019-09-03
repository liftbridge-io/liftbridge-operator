#!/usr/bin/env bash

# Copyright 2019 The Liftbridge Operator Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euxo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel)"
MOD_NAME="github.com/liftbridge-io/liftbridge-operator"

bash "${ROOT_DIR}/vendor/k8s.io/code-generator/generate-groups.sh" "deepcopy,client,informer,lister" \
  "${MOD_NAME}/pkg/generated" "${MOD_NAME}/pkg/apis" \
  liftbridge:v1alpha1 \
  --go-header-file "${ROOT_DIR}/hack/codegen/header.txt"
