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

package annotation_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/longxi-go/gosr/tools/pkg/annotation"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

func TestPositionUnknownDuplicateAndTarget(t *testing.T) {
	r, err := annotation.New("+test:", []annotation.Definition{
		{
			Name:   "test:name",
			Target: markers.DescribesType,
			Sample: "",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "input.go", "package x\n// +test:name=hello\ntype X struct{}", parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	group := annotation.TypeComments(f.Decls[0].(*ast.GenDecl), f.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec))
	values, err := r.Parse(fs, group, markers.DescribesType)
	if err != nil || len(values) != 1 || values[0].Value != "hello" || values[0].Position.Line != 2 {
		t.Fatal(values, err)
	}
	if _, err := r.Parse(fs, group, markers.DescribesField); err == nil {
		t.Fatal("wrong target accepted")
	}
	group.List = append(group.List, group.List[0])
	if _, err := r.Parse(fs, group, markers.DescribesType); err == nil {
		t.Fatal("duplicate accepted")
	}
	group.List[0].Text = "// +test:unknown"
	if _, err := r.Parse(fs, group, markers.DescribesType); err == nil {
		t.Fatal("unknown accepted")
	}
}
