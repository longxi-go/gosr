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

package output_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/longxi-go/gosr/pkg/codegen/artifact"
	"github.com/longxi-go/gosr/pkg/codegen/output"
)

func Example() {
	dir, err := os.MkdirTemp("", "codegen-example-*")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			panic(err)
		}
	}()
	set, err := artifact.New(artifact.File{
		Path:    filepath.Join(dir, "zz.gen.go"),
		Content: []byte("generated"),
	})
	if err != nil {
		panic(err)
	}
	err = output.Commit(context.Background(), set, output.Options{
		Owns: owned,
	})
	fmt.Println(err) // Output: <nil>
}
