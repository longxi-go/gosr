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

package annotation

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/longxi-go/gosr/pkg/codegen/diagnostic"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

// Parse ignores unrelated comments, rejects unknown/misplaced/duplicate markers
// and preserves parser causes and positions. Empty values are a schema policy.
func (r *Registry) Parse(fs *token.FileSet, group *ast.CommentGroup, target markers.TargetType) ([]Value, error) {
	var result []Value
	if group == nil {
		return result, nil
	}
	seen := make(map[string]bool)
	for _, comment := range group.List {
		if !strings.HasPrefix(comment.Text, "//") {
			continue
		}
		raw := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if !strings.HasPrefix(raw, r.prefix) {
			continue
		}
		var pos token.Position
		if fs != nil {
			pos = fs.Position(comment.Pos())
		}
		definition := r.registry.Lookup(raw, target)
		if definition == nil {
			return nil, diagnostic.Wrap("annotation", pos, "unknown or misplaced marker "+raw, nil)
		}
		value, err := definition.Parse(raw)
		if err != nil {
			return nil, diagnostic.Wrap("annotation", pos, raw, err)
		}
		key := fmt.Sprintf("%d:%s", target, definition.Name)
		if seen[key] && !r.repeat[key] {
			return nil, diagnostic.Wrap("annotation", pos, "duplicate marker "+definition.Name, nil)
		}
		seen[key] = true
		result = append(result, Value{
			Name:     definition.Name,
			Value:    value,
			Position: pos,
			Raw:      raw,
		})
	}
	return result, nil
}
