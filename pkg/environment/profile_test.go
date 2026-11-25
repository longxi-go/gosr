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

package environment_test

import (
	"fmt"
	"testing"

	"github.com/longxi-go/gosr/pkg/environment"
)

func ExampleParseProfilePredicate() {
	name, negated, err := environment.ParseProfilePredicate("!prod")
	if err != nil {
		panic(err)
	}
	fmt.Println(name, negated)
	// Output: prod true
}

func TestParseProfilePredicate(t *testing.T) {
	for _, tc := range []struct {
		input, name string
		negated     bool
	}{
		{
			"prod",
			"prod",
			false,
		},
		{
			"!prod",
			"prod",
			true,
		},
	} {
		name, negated, err := environment.ParseProfilePredicate(tc.input)
		if err != nil || name != tc.name || negated != tc.negated {
			t.Fatalf("predicate=%q name=%q negated=%v err=%v", tc.input, name, negated, err)
		}
	}
	for _, invalid := range []string{"", "!", "!!prod", "prod|test", " prod"} {
		if _, _, err := environment.ParseProfilePredicate(invalid); err == nil {
			t.Fatalf("invalid predicate %q accepted", invalid)
		}
	}
}
