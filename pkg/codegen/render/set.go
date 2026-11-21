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

package render

import (
	"bytes"
	"go/token"
	"io/fs"
	"text/template"

	"github.com/longxi-go/gosr/pkg/codegen/diagnostic"
)

// Set contains parsed immutable templates. Execute is safe for concurrent calls;
// each call owns its output. Views and helper side effects remain caller-owned.
type Set struct {
	templates *template.Template
}

// ParseFS enables text/template's missingkey=error policy for dot-key access.
// Built-in index and helper functions retain their own semantics. Callers own
// template assets, helpers and the semantic shape of view models.
func ParseFS(source fs.FS, functions template.FuncMap, patterns ...string) (*Set, error) {
	t, err := template.New("codegen").
		Option("missingkey=error").
		Funcs(functions).
		ParseFS(source, patterns...)
	if err != nil {
		return nil, diagnostic.Wrap("render", token.Position{}, "parse templates", err)
	}
	return &Set{
		templates: t,
	}, nil
}

// Execute reports template identity and preserves the template engine's cause.
func (s *Set) Execute(name string, view any) ([]byte, error) {
	var out bytes.Buffer
	if err := s.templates.ExecuteTemplate(&out, name, view); err != nil {
		return nil, diagnostic.Wrap("render", token.Position{}, "template "+name, err)
	}
	return out.Bytes(), nil
}
