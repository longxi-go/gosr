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

import "go/ast"

// TypeComments gives a type's own doc group precedence over its declaration
// group. It supports grouped declarations without inventing package scanning or
// bypassing the caller's generated-file and target selection rules.
func TypeComments(declaration *ast.GenDecl, spec *ast.TypeSpec) *ast.CommentGroup {
	if spec.Doc != nil {
		return spec.Doc
	}
	return declaration.Doc
}
