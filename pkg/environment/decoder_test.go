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
	"slices"
	"testing"

	"github.com/longxi-go/gosr/pkg/environment"
)

type applicationDecoder func(environment.Value) (string, error)

func decodeApplicationString(value environment.Value) (string, error) {
	text, _ := value.Scalar()
	return text, nil
}

func TestValueDecoderFuncCompatibility(t *testing.T) {
	var application applicationDecoder = decodeApplicationString
	var shared environment.ValueDecoderFunc[string] = application
	application = shared
	got, err := application(environment.String("typed"))
	if err != nil || got != "typed" {
		t.Fatal(got, err)
	}
}

func TestPlaceholderSyntaxContract(t *testing.T) {
	got := []string{
		environment.PlaceholderPrefix,
		environment.PlaceholderSuffix,
		environment.PlaceholderDefaultSeparator,
		environment.PlaceholderEscapePrefix,
		environment.EscapedPlaceholderPrefix,
		environment.DefaultProfile,
	}
	want := []string{"${", "}", ":", "\\", "\\${", "default"}
	if !slices.Equal(got, want) {
		t.Fatalf("syntax: got %q, want %q", got, want)
	}
}

func ExampleValueDecoderFunc() {
	decode := environment.ValueDecoderFunc[string](decodeApplicationString)
	text, _ := decode(environment.String("hello"))
	fmt.Println(text)
	// Output: hello
}
