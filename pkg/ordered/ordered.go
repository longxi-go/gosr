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

package ordered

import "math"

const (
	// HighestPrecedence precedes every larger order within the same group.
	// Its native-int width matches the return type of Ordered.Order.
	HighestPrecedence int = math.MinInt

	// LowestPrecedence is the fallback for objects without Ordered. Priority
	// membership is compared first, so this does not move a priority object
	// behind an ordinary object with HighestPrecedence.
	LowestPrecedence int = math.MaxInt
)

// Ordered supplies relative precedence within one extension set. Smaller values
// run first; consumers sample once per object per round before sorting copied
// keys. Implementations must keep reads quick, side-effect-free and safe for
// concurrent use. Consumers define equal-key identity or stability rules.
// Order does not replace creation dependencies or lifecycle phases.
type Ordered interface {
	Order() int
}

// PriorityOrdered belongs to the group before ordinary Ordered and unordered
// objects, regardless of numeric order. The extra marker is necessary because
// embedding Ordered alone cannot distinguish structural interfaces in Go.
// Consumers inspect membership without calling PriorityOrder.
type PriorityOrdered interface {
	Ordered
	PriorityOrder()
}
