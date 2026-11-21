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

package packageload_test

import (
	"testing"

	"github.com/longxi-go/gosr/tools/pkg/packageload"
)

func TestNilConfigUsesDefaultLoadOptions(t *testing.T) {
	loaded, err := packageload.Load(nil, "fmt")
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 1 || loaded[0].PkgPath != "fmt" || len(loaded[0].Errors) != 0 {
		t.Fatalf("default configuration did not load fmt: %v", loaded)
	}
}
