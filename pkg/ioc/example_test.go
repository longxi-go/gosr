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

package ioc_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/longxi-go/gosr/pkg/event"
	"github.com/longxi-go/gosr/pkg/ioc"
)

func ExampleBeanDefinition() {
	definition := ioc.BeanDefinition{
		ID:        "app/readDB",
		Name:      "readDB",
		Aliases:   []string{"primaryDB"},
		Contracts: []ioc.Contract{"example.Database"},
		Scope:     ioc.Singleton,
	}
	fmt.Println(definition.Validate())
	fmt.Println(definition.Matches(ioc.Selector{
		Name:     "primaryDB",
		Contract: "example.Database",
	}))
	// Output:
	// <nil>
	// true
}

func ExampleContextRefreshed() {
	// Notification types are portable values. Sending them at the correct
	// lifecycle boundary is the responsibility of an integrating container.
	var notification event.Event = ioc.ContextRefreshed{}
	_, refreshed := notification.(ioc.ContextRefreshed)
	_, closing := notification.(ioc.ContextClosing)
	fmt.Println(refreshed, closing)
	// Output: true false
}

func ExampleProvider() {
	provider := ioc.Provider[string](func(ctx context.Context) (string, error) {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		return "caller lifetime", nil
	})
	value, err := provider(context.Background())
	fmt.Println(value, err)
	// Output: caller lifetime <nil>
}

func ExampleOptional() {
	absent := ioc.Optional[int]{}
	present := ioc.Optional[int]{
		Present: true,
		Value:   0,
	}
	fmt.Println(absent.Present, present.Present, present.Value)
	// Output: false true 0
}

func ExampleDefinitionRegistry() {
	var registry ioc.DefinitionRegistry = &sliceRegistry{}
	err := registry.Register(ioc.BeanDefinition{
		ID:        "app/service",
		Name:      "service",
		Contracts: []ioc.Contract{"example.Service"},
		Scope:     ioc.Singleton,
	})
	fmt.Println(err, len(registry.Definitions()))
	// Output: <nil> 1
}

// bufferLease shows typed resource ownership without a container dependency.
type bufferLease struct {
	buffer *bytes.Buffer
}

func (l *bufferLease) Value() *bytes.Buffer {
	return l.buffer
}

func (l *bufferLease) Close(context.Context) error {
	l.buffer.Reset()
	return nil
}

func ExampleLease() {
	var lease ioc.Lease[*bytes.Buffer] = &bufferLease{
		buffer: bytes.NewBufferString("owned allocation"),
	}
	fmt.Println(lease.Value().String())
	if err := lease.Close(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println(lease.Value().Len())
	// Output:
	// owned allocation
	// 0
}
