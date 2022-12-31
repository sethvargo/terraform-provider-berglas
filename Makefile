# Copyright 2019 Seth Vargo
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

dev:
	@go install ./...
.PHONY: dev

generate:
	@go generate ./...
.PHONY: generate

test:
	@go test -count=1 -shuffle=on -short ./...
.PHONY: test

test-acc:
	@TF_ACC=1 go test -count=1 -shuffle=on -race ./...
.PHONY: test-acc
