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
	"go/token"
	"strings"
	"testing"

	"github.com/longxi-go/gosr/pkg/codegen/diagnostic"
)

func TestPositionAndCause(t *testing.T) {
	cause := errors.New("invalid field")
	err := diagnostic.Wrap("analyze", token.Position{
		Filename: "config.go",
		Line:     7,
		Column:   3,
	}, "Port", cause)
	var d *diagnostic.Error
	if !errors.As(err, &d) ||
		!errors.Is(err, cause) ||
		d.Phase != "analyze" ||
		!strings.Contains(err.Error(), "config.go:7:3") {
		t.Fatal(err)
	}
}
