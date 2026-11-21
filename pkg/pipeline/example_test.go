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

package pipeline_test

import (
	"context"
	"fmt"

	"github.com/longxi-go/gosr/pkg/pipeline"
)

type exampleState struct {
	completed []string
}

func analyzeExample(_ context.Context, state *exampleState) error {
	state.completed = append(state.completed, "analyze")
	return nil
}

func renderExample(_ context.Context, state *exampleState) error {
	state.completed = append(state.completed, "render")
	return nil
}

func ExampleRun() {
	state := &exampleState{}
	err := pipeline.Run(context.Background(), state,
		pipeline.NewStep("analyze", analyzeExample),
		pipeline.NewStep("render", renderExample),
	)
	fmt.Println(state.completed, err)
	// Output: [analyze render] <nil>
}
