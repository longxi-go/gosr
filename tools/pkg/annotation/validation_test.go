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
	"testing"

	"github.com/longxi-go/gosr/tools/pkg/annotation"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

func TestNilSampleReturnsAnError(t *testing.T) {
	_, err := annotation.New("+test:", []annotation.Definition{
		{
			Name:   "test:flag",
			Target: markers.DescribesType,
		},
	})
	if err == nil {
		t.Fatal("nil marker sample accepted")
	}
}

func TestMarkerPrefixRejectsCommentDelimitersAndWhitespace(t *testing.T) {
	for _, prefix := range []string{"/*", "demo:", "+demo: ", "+demo:\n"} {
		t.Run(prefix, func(t *testing.T) {
			if _, err := annotation.New(prefix, nil); err == nil {
				t.Fatalf("invalid annotation prefix %q accepted", prefix)
			}
		})
	}
}

func TestBlockCommentsAreIgnored(t *testing.T) {
	registry, err := annotation.New("+demo:", []annotation.Definition{
		{
			Name:   "demo:flag",
			Target: markers.DescribesType,
			Sample: struct{}{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	comments := &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text: "/* +demo:flag */",
			},
		},
	}
	values, err := registry.Parse(nil, comments, markers.DescribesType)
	if err != nil || len(values) != 0 {
		t.Fatalf("block comment parsed: %v, %v", values, err)
	}
}

func TestPlusPrefixSelectsAllRegisteredMarkers(t *testing.T) {
	registry, err := annotation.New("+", []annotation.Definition{
		{
			Name:   "demo:flag",
			Target: markers.DescribesType,
			Sample: struct{}{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	comments := &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text: "// +demo:flag",
			},
		},
	}
	values, err := registry.Parse(nil, comments, markers.DescribesType)
	if err != nil || len(values) != 1 || values[0].Name != "demo:flag" {
		t.Fatalf("registered marker not selected: %v, %v", values, err)
	}
}
