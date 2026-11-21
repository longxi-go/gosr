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

package diagnostic_test

import (
	"errors"
	"fmt"
	"go/token"

	"github.com/longxi-go/gosr/pkg/codegen/diagnostic"
)

func Example() {
	cause := errors.New("unsupported type")
	err := diagnostic.Wrap("analyze", token.Position{
		Filename: "app.go",
		Line:     3,
		Column:   1,
	}, "field Port", cause)
	fmt.Println(err)
	fmt.Println(errors.Is(err, cause))
	// Output:
	// app.go:3:1: analyze: field Port: unsupported type
	// true
}
