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
	"errors"
	"testing"
	"testing/fstest"

	"github.com/longxi-go/gosr/pkg/codegen/diagnostic"
	"github.com/longxi-go/gosr/pkg/codegen/render"
)

func TestStrictImmutableTemplates(t *testing.T) {
	s, err := render.ParseFS(fstest.MapFS{
		"a.tmpl": {
			Data: []byte(`{{define "item"}}{{.Name}}{{end}}`),
		},
	}, nil, "*.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.Execute("item", map[string]string{"Name": "ok"})
	if err != nil || string(got) != "ok" {
		t.Fatal(string(got), err)
	}
	_, err = s.Execute("item", map[string]string{})
	var d *diagnostic.Error
	if !errors.As(err, &d) || d.Phase != "render" {
		t.Fatal(err)
	}
	if _, err := s.Execute("missing", nil); err == nil {
		t.Fatal("missing template accepted")
	}
	if _, err := render.ParseFS(fstest.MapFS{
		"bad.tmpl": {
			Data: []byte(`{{`),
		},
	}, nil, "*.tmpl"); err == nil {
		t.Fatal("bad template accepted")
	}
}
