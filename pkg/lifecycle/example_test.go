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

package lifecycle_test

import (
	"context"
	"fmt"

	"github.com/longxi-go/gosr/pkg/lifecycle"
)

type allocation struct {
	ready bool
}

func (a *allocation) Init(context.Context) error {
	a.ready = true
	return nil
}

func (a *allocation) Destroy(context.Context) error {
	a.ready = false
	return nil
}

var (
	_ lifecycle.Initializer = (*allocation)(nil)
	_ lifecycle.Destroyer   = (*allocation)(nil)
)

func Example() {
	resource := &allocation{}
	ctx := context.Background()
	if err := resource.Init(ctx); err != nil {
		panic(err)
	}
	fmt.Println(resource.ready)
	if err := resource.Destroy(ctx); err != nil {
		panic(err)
	}
	fmt.Println(resource.ready)
	// Output:
	// true
	// false
}
