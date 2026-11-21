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

package render_test

import (
	"fmt"
	"testing/fstest"

	"github.com/longxi-go/gosr/pkg/codegen/render"
)

func Example() {
	set, err := render.ParseFS(fstest.MapFS{
		"file.tmpl": {
			Data: []byte("package {{.}}"),
		},
	}, nil, "*.tmpl")
	if err != nil {
		panic(err)
	}
	data, err := set.Execute("file.tmpl", "app")
	fmt.Println(string(data), err) // Output: package app <nil>
}
