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

package pipeline

import (
	"context"
	"fmt"
	"strings"
)

// Run executes steps in caller-supplied order against the same state value. Pass
// a pointer when steps need to mutate shared invocation state. Cancellation is
// checked between steps; each step remains responsible for observing it while
// executing. The first failure stops execution and retains its cause.
func Run[S any](ctx context.Context, state S, steps ...Step[S]) error {
	for index, step := range steps {
		if err := ctx.Err(); err != nil {
			return err
		}
		if step == nil {
			return fmt.Errorf("pipeline step %d is nil", index)
		}
		name := step.Name()
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("pipeline step %d has a blank name", index)
		}
		if err := step.Execute(ctx, state); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}
