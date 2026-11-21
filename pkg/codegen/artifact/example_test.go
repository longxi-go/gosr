// Copyright © 2026-present the longxi-go/gosr authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package artifact_test

import (
	"fmt"

	"github.com/longxi-go/gosr/pkg/codegen/artifact"
)

func Example() {
	set, err := artifact.New(artifact.File{
		Path:    "zz.gen.go",
		Content: []byte("package app"),
		Kind:    artifact.Go,
	})
	fmt.Println(len(set.Overlay()), err) // Output: 1 <nil>
}
