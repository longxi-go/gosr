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
	"errors"
	"fmt"

	"github.com/longxi-go/gosr/pkg/environment"
)

func ExampleValue() {
	precision, err := environment.Number("9007199254740993")
	if err != nil {
		panic(err)
	}
	config := environment.Object(map[string]environment.Value{
		"id":      precision,
		"enabled": environment.Bool(true),
	})
	fields := config.Fields()
	delete(fields, "enabled")
	id, scalar := config.Fields()["id"].Scalar()
	fmt.Println(id, scalar)
	fmt.Println(len(config.Fields()))
	fmt.Println(environment.Value{}.Kind() == environment.NullKind)
	// Output:
	// 9007199254740993 true
	// 2
	// true
}

// greetingSource shows that a source can implement the portable contract without
// depending on any framework, container, decoder or placeholder engine.
type greetingSource struct{}

func (greetingSource) Name() string {
	return "greeting"
}

func (greetingSource) Keys() []string {
	return []string{"greeting"}
}

func (greetingSource) Lookup(key string) (environment.Property, bool, error) {
	if err := environment.ValidateKey(key); err != nil {
		return environment.Property{}, false, err
	}
	if key != "greeting" {
		return environment.Property{}, false, nil
	}
	return environment.Property{
		Value: environment.String("Hello, ${name:world}"),
		Origin: environment.Origin{
			Source: "greeting",
			Key:    key,
		},
	}, true, nil
}

func ExamplePropertySource() {
	var source environment.PropertySource = greetingSource{}
	property, found, err := source.Lookup("greeting")
	if err != nil {
		panic(err)
	}
	raw, _ := property.Value.Scalar()
	fmt.Println(found, raw)
	fmt.Println(property.Origin.Source)
	// Output:
	// true Hello, ${name:world}
	// greeting
}

func ExampleError() {
	cause := errors.New("private connection details")
	err := &environment.Error{
		Code:  environment.SourceFailed,
		Key:   "db.password",
		Cause: cause,
	}
	fmt.Println(errors.Is(err, &environment.Error{
		Code: environment.SourceFailed,
	}))
	fmt.Println(errors.Is(err, cause))
	fmt.Println(err)
	// Output:
	// true
	// true
	// environment: source_failed: key=db.password
}
